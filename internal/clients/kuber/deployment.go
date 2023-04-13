package kuber

import (
	"fmt"
	"math"
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (client *Client) FindDeployment(namespaceName, name string) (*appsv1.Deployment, error) {
	deployments := client.clientset.AppsV1().Deployments(namespaceName)

	return deployments.Get(client.context, name, metav1.GetOptions{})
}

func (client *Client) findOrInitializeDeployment(
	namespace *corev1.Namespace,
	deploymentName string,
	containerList *domainTypes.ContainerList,
) (*appsv1.Deployment, error) {
	deployments := client.clientset.AppsV1().Deployments(namespace.Name)

	deployment, err := deployments.Get(client.context, deploymentName, metav1.GetOptions{})

	if err != nil && !errors.IsNotFound(err) {
		return deployment, err
	}

	if len(deployment.UID) == 0 {
		deployment = initializeDeployment(namespace, deploymentName, containerList)
	}

	return deployment, nil
}

// Given a memory resource, provide the proportional CPU resource in millicores.
func cpuProportion(memoryQuantity resource.Quantity) *resource.Quantity {
	var (
		numMilliCores  = global.Settings.PoolMachineTotalCpuMillicores
		numMemoryBytes = global.Settings.PoolMachineTotalMemoryBytes
	)

	return resource.NewMilliQuantity(memoryQuantity.Value()*numMilliCores/numMemoryBytes, resource.DecimalSI)
}

func podCpuProportion(totalMemory *resource.Quantity) *resource.Quantity {
	const totalPercentage = 100.0

	return resource.NewMilliQuantity(
		int64(math.Ceil(float64(cpuProportion(*resource.NewQuantity(totalMemory.Value(), resource.DecimalSI)).MilliValue())*
			float64(global.Settings.DefaultAutoscalingCpuThreshold+global.Settings.DefaultAutoscalingCpuThresholdEpsilon)/totalPercentage)), resource.DecimalSI)
}

func ephemeralStorageProportion(memoryQuantity resource.Quantity) *resource.Quantity {
	ephemeralStorageCoefficient := global.Settings.EphemeralStorageCoefficient
	return resource.NewQuantity(int64(math.Round(float64(memoryQuantity.Value())*ephemeralStorageCoefficient)), resource.DecimalSI)
}

func (client *Client) updateDeploymentAttributes(
	namespace *corev1.Namespace,
	deployment *appsv1.Deployment,
	containerList domainTypes.ContainerList,
	composeFile domainTypes.ComposeFile,
	hostVolumeFileList *domainTypes.HostVolumeFileList,
) (*appsv1.Deployment, error) {
	replicaCount := global.Settings.CustomerDefaultReplicationFactor

	deployment.Spec.Replicas = &replicaCount

	if len(deployment.Spec.Template.ObjectMeta.Annotations) == 0 {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}

	deployment.Spec.Template.Spec.Volumes = prepareDeploymentVolumes(containerList, hostVolumeFileList)

	if containerList.IsAnyVolumeExists() {
		deployment.Spec.Strategy = buildRecreateDeploymentStrategy()
	} else {
		deployment.Spec.Strategy = buildDefaultDeploymentStrategy()
	}

	// DO NOT DELETE THIS LINE. This is necessary to get an up-to-date container image each time you deploy.
	deployment.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	containers, err := client.updateDeploymentContainers(namespace, containerList)

	if err != nil {
		return deployment, err
	}

	deployment.Spec.Template.Spec.Containers = containers
	deployment.Spec.Template.Spec.HostAliases = []corev1.HostAlias{
		{
			IP:        global.Settings.DefaultIp,
			Hostnames: buildAllowedHostnames(&containerList),
		},
	}

	initContainers := []corev1.Container{}
	initContainerForHostVolumes, err := buildInitContainerForHostVolumes(containerList, composeFile, hostVolumeFileList)

	if err != nil {
		return nil, err
	}

	if initContainerForHostVolumes.Name != "" {
		initContainers = append(initContainers, initContainerForHostVolumes)
	}

	deployment.Spec.Template.Spec.InitContainers = initContainers

	return deployment, nil
}

func (client *Client) updateDeploymentContainers(
	namespace *corev1.Namespace,
	containerList domainTypes.ContainerList,
) ([]corev1.Container, error) {
	var containers []corev1.Container

	for _, draftContainer := range containerList.Items {
		containerImage, err := draftContainer.NameWithTag()
		if err != nil {
			return containers, err
		}

		var containerPorts []corev1.ContainerPort

		if draftContainer.IsPublic() {
			defaultContainerPortName := "default-port"
			containerPort := corev1.ContainerPort{Name: defaultContainerPortName, ContainerPort: *draftContainer.Port}
			containerPorts = append(containerPorts, containerPort)
		}

		memoryRequest, err := resource.ParseQuantity(fmt.Sprintf("%vMi", draftContainer.MemoryRequest))
		if err != nil {
			return nil, err
		}

		// podCpuRequest := podCpuProportion(&memoryRequest)
		podCpuRequest := resource.NewMilliQuantity(int64(0), resource.DecimalSI) //nolint: gomnd

		requests := corev1.ResourceList{
			corev1.ResourceMemory: memoryRequest,
			corev1.ResourceCPU:    *podCpuRequest,
		}

		memoryLimit, err := resource.ParseQuantity(fmt.Sprintf("%vMi", draftContainer.MemoryLimit))
		if err != nil {
			return nil, err
		}

		// podCpuLimit := podCpuProportion(&memoryLimit)
		podCpuLimit := resource.NewMilliQuantity(int64(1000), resource.DecimalSI) //nolint: gomnd

		limits := corev1.ResourceList{
			corev1.ResourceMemory: memoryLimit,
			corev1.ResourceCPU:    *podCpuLimit,
		}

		name := global.Settings.ResourceName.ContainerSecret(draftContainer.ID)
		secret, err := client.GetSecret(namespace.Name, name)

		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}

		containerVariables := prepareContainerEnvironmentVariables(draftContainer)
		containerSecrets := prepareContainerSecrets(draftContainer, secret)

		container := corev1.Container{
			Name:            draftContainer.ControllerName,
			Image:           containerImage,
			ImagePullPolicy: "Always",
			Ports:           containerPorts,
			Resources: corev1.ResourceRequirements{
				Requests: requests,
				Limits:   limits,
			},
			Env:           append(containerVariables, containerSecrets...),
			VolumeMounts:  prepareContainerVolumeMounts(draftContainer),
			LivenessProbe: prepareContainerHealthcheck(draftContainer),
		}

		container.Command = draftContainer.Entrypoint
		container.Args = draftContainer.Command

		if draftContainer.Port != nil {
			port := intstr.FromInt(int(*draftContainer.Port))

			if draftContainer.TargetPort != nil {
				port = intstr.FromInt(int(*draftContainer.TargetPort))
			}

			container.StartupProbe = &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: port,
					},
				},
				InitialDelaySeconds: global.Settings.StartupProbeDelaySeconds,
				FailureThreshold:    global.Settings.StartupProbeFailureThreshold,
				PeriodSeconds:       global.Settings.StartupProbePeriodSettings,
			}
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func (client *Client) CreateOrUpdateDeployments(
	namespace *corev1.Namespace,
	deploymentName string,
	containerList domainTypes.ContainerList,
	credentials []domainTypes.Credential,
	composeFile domainTypes.ComposeFile,
	hostVolumeFileList *domainTypes.HostVolumeFileList,
) (*appsv1.Deployment, error) {
	deployments := client.clientset.AppsV1().Deployments(namespace.Name)

	deployment, err := client.findOrInitializeDeployment(namespace, deploymentName, &containerList)
	if err != nil {
		return nil, err
	}

	deployment, err = client.updateDeploymentAttributes(namespace, deployment, containerList, composeFile, hostVolumeFileList)
	if err != nil {
		return nil, err
	}

	if len(credentials) > 0 {
		deployment.Spec.Template.Spec.ImagePullSecrets = prepareCredentialsDeployment(credentials)
	}

	if len(deployment.UID) > 0 {
		deployment, err = deployments.Update(client.context, deployment, metav1.UpdateOptions{})
	} else {
		deployment, err = deployments.Create(client.context, deployment, metav1.CreateOptions{})
	}

	return deployment, err
}

func (client *Client) RemoveDeployments(namespaceName, name string) error {
	deployments := client.clientset.AppsV1().Deployments(namespaceName)

	err := deployments.Delete(client.context, name, metav1.DeleteOptions{})

	if err != nil {
		switch {
		case errors.IsNotFound(err):
			return nil
		default:
			return err
		}
	}

	return nil
}

package kuber

import (
	"fmt"
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	domainTypes "gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/types/domain"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

func initializeDeployment(
	namespace *corev1.Namespace,
	deploymentName string,
	containerList *domainTypes.ContainerList,
) *appsv1.Deployment {
	var deploymentStrategy appsv1.DeploymentStrategy

	if containerList.IsAnyVolumeExists() {
		deploymentStrategy = buildRecreateDeploymentStrategy()
	} else {
		deploymentStrategy = buildDefaultDeploymentStrategy()
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			Labels: map[string]string{
				"app":                          "controller",
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
				"app.kubernetes.io/instance":   namespace.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": deploymentName},
			},
			Strategy: deploymentStrategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentName,
					Labels: map[string]string{
						"app":                          deploymentName,
						"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
					},
				},
				Spec: corev1.PodSpec{
					Containers:  []corev1.Container{},
					Affinity:    getPodSpecAffinity(),
					Tolerations: getPodSpecTolerations(),
					HostAliases: []corev1.HostAlias{
						{
							IP:        global.Settings.DefaultIp,
							Hostnames: buildAllowedHostnames(containerList),
						},
					},
					AutomountServiceAccountToken: ptr.To(false), // False. Security, DO NOT REMOVE
				},
			},
		},
	}
}

func generateServiceName() string {
	now := time.Now()
	unix := now.Unix()

	return fmt.Sprintf("service-%v", unix)
}

func generateIngressName() string {
	now := time.Now()
	unix := now.Unix()

	return fmt.Sprintf("ingress-%v", unix)
}

func initializeService(namespace, serviceName string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
			Labels: map[string]string{
				"app":                          "controller",
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{},
			Type:  corev1.ServiceTypeClusterIP,
		},
	}
}

func initializeIngress(namespace, name string) *networkingV1.Ingress {
	return &networkingV1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app":                          "controller",
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
			Namespace:   namespace,
			Annotations: map[string]string{},
		},
		Spec: networkingV1.IngressSpec{},
	}
}

func initializeHorizontalPodAutoscaler(
	namespaceName, deploymentName, name string,
	minReplicas, maxReplicas int32,
) *autoscalingV1.HorizontalPodAutoscaler {
	var cpu int32 = int32(global.Settings.DefaultAutoscalingCpuThreshold) // issue #75

	return &autoscalingV1.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespaceName,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
			},
		},
		Spec: autoscalingV1.HorizontalPodAutoscalerSpec{
			MinReplicas:                    &minReplicas,
			MaxReplicas:                    maxReplicas,
			TargetCPUUtilizationPercentage: &cpu,
			ScaleTargetRef: autoscalingV1.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       deploymentName,
				APIVersion: "apps/v1",
			},
		},
	}
}

func getPodSpecAffinity() *corev1.Affinity {
	if global.Settings.SandboxEnabled {
		return &corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      "sandbox.gke.io/runtime",
									Operator: corev1.NodeSelectorOpIn,
									Values:   []string{"gvisor"},
								},
							},
						},
					},
				},
				PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
					{
						Preference: corev1.NodeSelectorTerm{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      "cloud.google.com/gke-spot",
									Operator: corev1.NodeSelectorOpIn,
									Values:   []string{"true"},
								},
							},
						},
						Weight: 100, //nolint: gomnd
					},
				},
			},
		}
	}

	return nil
}

func getPodSpecTolerations() []corev1.Toleration {
	if global.Settings.SandboxEnabled {
		return []corev1.Toleration{
			{
				Key:      "sandbox.gke.io/runtime",
				Operator: "Exists",
			},
		}
	}

	return nil
}

func buildDefaultDeploymentStrategy() appsv1.DeploymentStrategy {
	return appsv1.DeploymentStrategy{
		Type: appsv1.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: 0},
			MaxSurge:       &intstr.IntOrString{Type: intstr.Int, IntVal: 1},
		},
	}
}

func buildRecreateDeploymentStrategy() appsv1.DeploymentStrategy {
	return appsv1.DeploymentStrategy{
		Type: appsv1.RecreateDeploymentStrategyType,
	}
}

func buildAllowedHostnames(containerList *domainTypes.ContainerList) []string {
	hostnames := []string{}

	for _, container := range containerList.Items {
		hostname := container.ServiceName
		hostnames = append(hostnames, hostname)
	}

	return hostnames
}

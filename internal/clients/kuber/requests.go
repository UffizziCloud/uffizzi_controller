package kuber

import (
	"fmt"
	"time"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingV1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
)

func initializeDeployment(
	namespace *corev1.Namespace,
	deploymentName, deploymentSelectorName string) *appsv1.Deployment {
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
				MatchLabels: map[string]string{"app": deploymentSelectorName},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
					MaxSurge:       &intstr.IntOrString{Type: intstr.String, StrVal: "25%"},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentName,
					Labels: map[string]string{
						"app":                          deploymentSelectorName,
						"app.kubernetes.io/managed-by": global.Settings.ManagedApplication,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{},
					Tolerations: []corev1.Toleration{
						{
							Key:      "sandbox.gke.io/runtime",
							Operator: "Exists",
						},
					},

					AutomountServiceAccountToken: pointer.BoolPtr(false), // False. Security, DO NOT REMOVE
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

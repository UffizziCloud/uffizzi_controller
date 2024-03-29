package kuber

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"gitlab.com/dualbootpartners/idyl/uffizzi_controller/internal/global"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (client *Client) GetPods(namespace string) (*v1.PodList, error) {
	pods, err := client.clientset.CoreV1().Pods(namespace).List(client.context, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/managed-by=%v", global.Settings.ManagedApplication),
	})

	return pods, err
}

func (client *Client) GetPodsMetrics(namespace string) (*v1beta1.PodMetricsList, error) {
	metrics, err := client.metricClient.MetricsV1beta1().PodMetricses(namespace).List(client.context, metav1.ListOptions{})

	return metrics, err
}

func (client *Client) GetPodLogs(
	namespace string,
	podName string,
	containerName string,
	limit int64,
	previous bool,
) ([]string, error) {
	logOptions := &v1.PodLogOptions{
		Previous:   previous,
		Container:  containerName,
		TailLines:  &limit,
		Timestamps: true,
	}

	request := client.clientset.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	podLogs, err := request.Stream(client.context)

	if err != nil {
		return nil, err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)

	if err != nil {
		return nil, err
	}

	str := buf.String()

	return strings.Split(str, "\n"), nil
}

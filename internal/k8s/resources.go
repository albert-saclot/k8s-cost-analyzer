package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
)

type PodResources struct {
	Name         string
	Namespace    string
	CPURequest   string
	MemoryRequest string
	CPULimit     string
	MemoryLimit  string
}

// FetchPods retrieves all pods from the specified namespace
func FetchPods(ctx context.Context, client *kubernetes.Clientset, namespace string) ([]corev1.Pod, error) {
	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return pods.Items, nil
}

// ExtractResources parses resource requests and limits from a pod
func ExtractResources(pod corev1.Pod) PodResources {
	var cpuRequest, memoryRequest, cpuLimit, memoryLimit resource.Quantity

	for _, container := range pod.Spec.Containers {
		if container.Resources.Requests != nil {
			if cpu, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				cpuRequest.Add(cpu)
			}
			if mem, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				memoryRequest.Add(mem)
			}
		}
		if container.Resources.Limits != nil {
			if cpu, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				cpuLimit.Add(cpu)
			}
			if mem, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				memoryLimit.Add(mem)
			}
		}
	}

	return PodResources{
		Name:          pod.Name,
		Namespace:     pod.Namespace,
		CPURequest:    formatQuantity(cpuRequest),
		MemoryRequest: formatQuantity(memoryRequest),
		CPULimit:      formatQuantity(cpuLimit),
		MemoryLimit:   formatQuantity(memoryLimit),
	}
}

func formatQuantity(q resource.Quantity) string {
	if q.IsZero() {
		return "-"
	}
	return q.String()
}

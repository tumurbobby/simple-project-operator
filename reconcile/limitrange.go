package reconcile

import (
	"context"
	"log/slog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func CreateLimitRange(clientset *kubernetes.Clientset, namespace string, logger *slog.Logger) {
	lr := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-limits",
			Namespace: namespace,
		},
		Spec: corev1.LimitRangeSpec{
			Limits: []corev1.LimitRangeItem{
				{
					Type: corev1.LimitTypeContainer,
					Default: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("250m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
					DefaultRequest: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("125m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
				},
			},
		},
	}
	_, err := clientset.CoreV1().
		LimitRanges(lr.Namespace).
		Create(
			context.Background(),
			lr,
			metav1.CreateOptions{},
		)
	if err != nil {
		logger.Error(
			"Failed to create LimitRange",
			"namespace", lr.Namespace,
			"error", err,
		)
		return
	}
	logger.Info(
		"LimitRange created",
		"namespace", namespace,
	)
}

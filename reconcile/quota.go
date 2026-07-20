package reconcile

import (
	"context"
	"log/slog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateResourceQuota(clientset *kubernetes.Clientset, namespace string, logger *slog.Logger) {
	rq := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-quota",
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourcePods:   resource.MustParse("5"),
				corev1.ResourceCPU:    resource.MustParse("1000m"),
				corev1.ResourceMemory: resource.MustParse("1000Mi"),
			},
		},
	}
	_, err := clientset.CoreV1().
		ResourceQuotas(rq.Namespace).
		Create(
			context.Background(),
			rq,
			metav1.CreateOptions{},
		)
	if err != nil {
		logger.Error(
			"Failed to create ResourceQuota",
			"namespace", rq.Namespace,
			"error", err,
		)
		return
	}
	logger.Info(
		"ResourceQuota created",
		"namespace", namespace,
	)
}

package reconcile

import (
	"context"
	"log/slog"

	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateNetworkPolicy(clientset *kubernetes.Clientset, namespace string, logger *slog.Logger) {
	np := &networkv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny-ingress",
			Namespace: namespace,
		},
		Spec: networkv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkv1.PolicyType{
				networkv1.PolicyTypeIngress,
			},
		},
	}
	_, err := clientset.NetworkingV1().
		NetworkPolicies(np.Namespace).
		Create(
			context.Background(),
			np,
			metav1.CreateOptions{},
		)
	if err != nil {
		logger.Error(
			"Failed to create NetworkPolicy",
			"namespace", np.Namespace,
			"error", err,
		)
		return
	}
	logger.Info(
		"NetworkPolicy created",
		"namespace", namespace,
	)
}

package main

import (
	"context"
	"log/slog"
	"strings"

	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func shouldSkipNamespace(name string) bool {
	return strings.HasPrefix(name, "kube-")
}

func createResourceQuota(clientset *kubernetes.Clientset, namespace string, logger *slog.Logger) {
	rq := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-quota",
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourcePods:   resource.MustParse("2"),
				corev1.ResourceCPU:    resource.MustParse("500m"),
				corev1.ResourceMemory: resource.MustParse("500Mi"),
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

func createLimitRange(clientset *kubernetes.Clientset, namespace string, logger *slog.Logger) {
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
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("512Mi"),
					},
					DefaultRequest: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
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

func createNetworkPolicy(clientset *kubernetes.Clientset, namespace string, logger *slog.Logger) {
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

func main() {
	logger := slog.Default()
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	for {
		watcher, err := clientset.CoreV1().
			Namespaces().
			Watch(
				context.Background(),
				metav1.ListOptions{},
			)
		if err != nil {
			panic(err)
		}
		for event := range watcher.ResultChan() {

			ns, ok := event.Object.(*corev1.Namespace)
			if !ok {
				continue
			}
			if shouldSkipNamespace(ns.Name) {
				continue
			}
			if event.Type == watch.Added {
				logger.Info(
					"Namespace added",
					"namespace", ns.Name,
				)
				createResourceQuota(clientset, ns.Name, logger)
				createLimitRange(clientset, ns.Name, logger)
				createNetworkPolicy(clientset, ns.Name, logger)
			}
		}
		logger.Warn("Watch ended, reconnecting")
	}

}

package main

import (
	"context"
	"log/slog"
	"strings"

	"github.com/tumurbobby/namespace-bootstrap-controller/reconcile"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func shouldSkipNamespace(name string) bool {
	return strings.HasPrefix(name, "kube-")
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
				reconcile.CreateResourceQuota(clientset, ns.Name, logger)
				reconcile.CreateLimitRange(clientset, ns.Name, logger)
				reconcile.CreateNetworkPolicy(clientset, ns.Name, logger)
			}
		}
		logger.Warn("Watch ended, reconnecting")
	}

}

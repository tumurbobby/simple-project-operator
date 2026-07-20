package main

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Load sa
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	// Create Kubernetes client
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
			if !strings.HasPrefix(ns.Name, "bobby-") {
				continue
			}
			if event.Type == watch.Added {
				rq := &corev1.ResourceQuota{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "default-quota",
						Namespace: ns.Name,
					},
					Spec: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourcePods:   resource.MustParse("2"),
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("500Mi"),
						},
					},
				}
				_, err = clientset.CoreV1().
					ResourceQuotas(rq.Namespace).
					Create(
						context.Background(),
						rq,
						metav1.CreateOptions{},
					)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			fmt.Println("Watch ended, reconnecting...")
		}
	}

}

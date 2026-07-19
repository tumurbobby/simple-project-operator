package main

import (
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Load kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join("/home/vagrant/.kube", "config"))
	if err != nil {
		panic(err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "bobby-demo",
		},
	}
	_, err = clientset.CoreV1().
		Namespaces().
		Create(
			context.Background(),
			ns,
			metav1.CreateOptions{},
		)
	if err != nil {
		panic(err)
	}

	fmt.Println("Namespace created!")
}

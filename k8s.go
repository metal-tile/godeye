package main

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func updatePod() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "rest.InClusterConfig")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "kubernetes.NewForConfig")
	}

	{
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "Pod.List")
		}
		fmt.Printf("There are %d Pod in the cluster\n", len(pods.Items))
		for _, item := range pods.Items {
			fmt.Printf("Pod %s exists.\n", item.Name)
		}
	}

	{
		rcl, err := clientset.CoreV1().ReplicationControllers("").List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "ReplicationControllers.List")
		}
		fmt.Printf("There are %d ReplicationController in the cluster\n", len(rcl.Items))
		for _, item := range rcl.Items {
			fmt.Printf("ReplicationController %s is size = %d\n", item.Name, item.Size())
		}
	}

	return nil
}

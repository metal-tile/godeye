package main

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func updateReplicas(namespace string, name string, replicas int32) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return errors.Wrap(err, "failed rest.InClusterConfig")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed kubernetes.NewForConfig")
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("failed get Deployment %+v\n", err)
		return errors.Wrap(err, "failed get deployment")
	}
	deployment.Spec.Replicas = &replicas
	fmt.Printf("Deployment %v\n", deployment)
	ug, err := clientset.AppsV1().Deployments(deployment.Namespace).Update(deployment)
	if err != nil {
		fmt.Printf("failed update Deployment %+v", err)
		return errors.Wrap(err, "failed update Deployment")
	}
	fmt.Printf("done update deployment %v\n", ug)

	return nil
}

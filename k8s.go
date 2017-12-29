package main

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		rsl, err := clientset.AppsV1beta2().ReplicaSets("").List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "ReplicaSet.List")
		}
		fmt.Printf("There are %d ReplicaSet in the cluster\n", len(rsl.Items))
		for _, item := range rsl.Items {
			fmt.Printf("ReplicaSet %s exists. \n", item.Name)
		}
	}

	{
		dl, err := clientset.AppsV1beta2().Deployments("").List(metav1.ListOptions{})
		if err != nil {

		}
		fmt.Printf("There are %d Deployments in the cluster\n", len(dl.Items))
		for _, item := range dl.Items {
			fmt.Printf("Deployments %s exists. \n", item.Name)
		}
		d, err := clientset.AppsV1beta2().Deployments("").Update(&v1beta2.Deployment{
			ObjectMeta: apiv1.ObjectMeta{
				Name: "godeye-node",
			},
			Status: v1beta2.DeploymentStatus{
				Replicas: 0,
			},
		})
		if err != nil {
			return errors.Wrap(err, "Update Deployment")
		}
		fmt.Printf("Updated Deployment %v", d)
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

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
		return errors.Wrap(err, "failed rest.InClusterConfig")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed kubernetes.NewForConfig")
	}

	{
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "failed Pod.List")
		}
		fmt.Printf("There are %d Pod in the cluster\n", len(pods.Items))
		for _, item := range pods.Items {
			fmt.Printf("Pod %s exists.\n", item.Name)
		}
	}

	{
		rsl, err := clientset.AppsV1beta2().ReplicaSets("").List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "failed ReplicaSet.List")
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
		godeye, err := clientset.AppsV1beta2().Deployments("").Get("godeye-node", metav1.GetOptions{})
		if err != nil {
			return errors.Wrap(err, "failed get godeye deployment")
		}
		godeye.Status.Replicas = 0
		fmt.Printf("godeye Deployment %v", godeye)
		ug, err := clientset.AppsV1beta2().Deployments("").Update(godeye)
		if err != nil {
			return errors.Wrap(err, "Update Deployment")
		}
		fmt.Printf("Updated Deployment %v", ug)
	}

	{
		rcl, err := clientset.CoreV1().ReplicationControllers("").List(metav1.ListOptions{})
		if err != nil {
			return errors.Wrap(err, "failed ReplicationControllers.List")
		}
		fmt.Printf("There are %d ReplicationController in the cluster\n", len(rcl.Items))
		for _, item := range rcl.Items {
			fmt.Printf("ReplicationController %s is size = %d\n", item.Name, item.Size())
		}
	}

	return nil
}

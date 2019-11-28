package resource

import (
	"github.com/jbaojunior/zeropods/auth"
	"k8s.io/client-go/kubernetes"
)

// clientset
var clientSet *kubernetes.Clientset

// ActionDeployments is a function to create connection and Up/Down the deployments
func ActionDeployments(action, namespace, connection, kubeconfig string) {

	clientSet = auth.ConnectCluster(connection, kubeconfig)

	if action == "up" {
		upScaleDeployments(namespace)
		upScaleStatefulSets(namespace)
	} else if action == "down" {
		downScaleDeployments(namespace)
		downScaleStatefulSets(namespace)
	}
}

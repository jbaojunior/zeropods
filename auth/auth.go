package auth

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// To connect inside cluster using service account and token
func clusterkubeConnect() *kubernetes.Clientset {
	// Create internal cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("[ERROR] Problems to create cluster config\n")
		panic(err)
	}

	// Creates clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("[ERROR] Problems to create config\n")
		panic(err)
	}
	return clientSet
}

// To connect using kubeconfig. To external tools or tests
// Extract from https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go
func configKubeConnect(kubeconfig string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

// ConnectCluster is a function to auth in cluster
func ConnectCluster(kind, kubeconfig string) *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	if kind == "cluster" {
		clientset = clusterkubeConnect()
	} else if kind == "config" {
		if kubeconfig != "" {
			clientset = configKubeConnect(kubeconfig)
		} else {
			log.Fatalln("Connection method \"config\" but the kubeconfig is not defined. Please verify.")
		}

	} else {
		log.Fatalln("Wrong type of connection method (argument --conn). The possible values are \"cluster\" or \"config\"")
	}

	return clientset
}

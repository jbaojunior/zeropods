package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	resource "github.com/jbaojunior/zeropods/resources"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
)

// ClientSet is a old method
var ClientSet *kubernetes.Clientset

func main() {

	fmt.Println("OS Args:", os.Args)

	var action = flag.String("action", "", "Action to do. Possible values are \"up\" or \"down\"")
	var namespace = flag.String("n", "", "Namespace to do the action")
	var connection = flag.String("conn", "cluster", "Connect method to cluster. Possible values are \"cluster\" and \"config\".\n\t - \"cluster\" is to deploy inside a cluster, using a Service Account.\n\t - \"config\" is to using outsite of cluster, with a kubeconfig")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file. Usage when the parameter \"connection\" is \"config\"")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file. Usage when the parameter \"connection\" is \"config\"")
	}

	flag.Parse()

	if *namespace == "" {
		log.Fatal("Namespace not defined (option -n ). Please verify.")
	}

	if *action == "up" || *action == "down" {
		resource.ActionDeployments(*action, *namespace, *connection, *kubeconfig)
	} else {
		log.Fatal("Wrong scale action (argument -c). The possible values are \"up\" or \"down\"")
	}
}

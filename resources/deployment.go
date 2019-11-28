package resource

import (
	"fmt"
	"log"
	"strconv"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func downScaleDeployments(namespace string) {
	resource := clientSet.AppsV1().Deployments(namespace)
	resourceMetadata, err := resource.List(metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, items := range resourceMetadata.Items {
		var numReplicas int32 = 0
		resourceName := items.ObjectMeta.Name
		resourceNamespace := items.ObjectMeta.Namespace
		actualReplicas := *items.Spec.Replicas
		if actualReplicas != 0 {
			deploymentSetScale(resourceNamespace, resourceName, numReplicas)
			deploymentSetAnnotations(resourceNamespace, resourceName, actualReplicas)
		} else {
			log.Println("Deployment - Nothing to scale. Scale already is zero.")
		}
	}
}

func upScaleDeployments(namespace string) {
	resource := clientSet.AppsV1().Deployments(namespace)
	resourceMetadata, err := resource.List(metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, items := range resourceMetadata.Items {
		lastNumReplicas, err := strconv.Atoi(items.ObjectMeta.Annotations["zeropods/last-scale-number"])
		if err != nil {
			log.Println("Deployment - Problems to get annotation \"zeropods/last-scale-number\". Please execute a scale down first or create this annotation.")
		} else {
			numReplicas := int32(lastNumReplicas)
			resourceName := items.ObjectMeta.Name
			resourceNamespace := items.ObjectMeta.Namespace
			if numReplicas != 0 {
				deploymentSetScale(resourceNamespace, resourceName, numReplicas)
			} else {
				log.Println("Deployment - Nothing to scale. Scale already is zero.")
			}
		}
	}
}

// replicaScale scale the deployment
func deploymentSetScale(namespace, name string, numReplicas int32) {
	resource := clientSet.AppsV1().Deployments(namespace)
	scale, _ := resource.GetScale(name, metaV1.GetOptions{})
	scale.Spec.Replicas = numReplicas

	_, err := resource.UpdateScale(name, scale)
	if err != nil {
		log.Printf("[ERROR] Deployment %s scale to %v replicas failed: %v\n", name, numReplicas, err)
		log.Fatal()
	} else {
		log.Printf("Deployment: %s, Replicas: %v - scale with success.\n", name, numReplicas)
	}
}

// setAnnotation is a annotation to known the number of replicas before scale down. With this is possible return to the correct number of replicas.
func deploymentSetAnnotations(namespace, name string, numReplicas int32) {
	var dataPatch []byte
	dataPatch = append(dataPatch, fmt.Sprintf("{\"metadata\":{\"annotations\":{\"zeropods/last-scale-number\":\"%v\"}}}", numReplicas)...)
	_, err := clientSet.AppsV1().Deployments(namespace).Patch(name, "application/strategic-merge-patch+json", dataPatch)

	if err != nil {
		log.Printf("[ERROR] Deployment annotation error: %s\n", err)
		log.Fatal()
	}
}

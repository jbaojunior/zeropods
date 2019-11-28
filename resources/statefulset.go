package resource

import (
	"fmt"
	"log"
	"strconv"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func downScaleStatefulSets(namespace string) {
	resource := clientSet.AppsV1().StatefulSets(namespace)
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
			statefulSetScale(resourceNamespace, resourceName, numReplicas)
			statefulSetAnnotations(resourceNamespace, resourceName, actualReplicas)
		} else {
			log.Println("StatefulSet - Nothing to scale. Scale already is zero.")
		}
	}
}

func upScaleStatefulSets(namespace string) {
	resource := clientSet.AppsV1().StatefulSets(namespace)
	resourceMetadata, err := resource.List(metaV1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, items := range resourceMetadata.Items {
		lastNumReplicas, err := strconv.Atoi(items.ObjectMeta.Annotations["zeropods/last-scale-number"])
		if err != nil {
			log.Println("StatefulSet - Problems to get annotation \"zeropods/last-scale-number\". Please execute a scale down first or create this annotation.")
		} else {
			numReplicas := int32(lastNumReplicas)
			resourceName := items.ObjectMeta.Name
			resourceNamespace := items.ObjectMeta.Namespace
			if numReplicas != 0 {
				statefulSetScale(resourceNamespace, resourceName, numReplicas)
			} else {
				log.Println("StatefulSet - Nothing to scale. Scale already is zero.")
			}
		}
	}
}

// replicaScale scale the deployment
func statefulSetScale(namespace, name string, numReplicas int32) {
	resource := clientSet.AppsV1().StatefulSets(namespace)
	scale, _ := resource.GetScale(name, metaV1.GetOptions{})
	scale.Spec.Replicas = numReplicas

	_, err := resource.UpdateScale(name, scale)
	if err != nil {
		log.Printf("[ERROR] StatefulSet %s scale to %v replicas failed: %v\n", name, numReplicas, err)
		log.Fatal()
	} else {
		log.Printf("StatefulSet: %s, Replicas: %v - scale with success.\n", name, numReplicas)
	}
}

// setAnnotation is a annotation to known the number of replicas before scale down. With this is possible return to the correct number of replicas.
func statefulSetAnnotations(namespace, name string, numReplicas int32) {
	var dataPatch []byte
	dataPatch = append(dataPatch, fmt.Sprintf("{\"metadata\":{\"annotations\":{\"zeropods/last-scale-number\":\"%v\"}}}", numReplicas)...)
	_, err := clientSet.AppsV1().StatefulSets(namespace).Patch(name, "application/strategic-merge-patch+json", dataPatch)

	if err != nil {
		log.Printf("[ERROR] StatefulSet annotation error: %s\n", err)
		log.Fatal()
	}
}

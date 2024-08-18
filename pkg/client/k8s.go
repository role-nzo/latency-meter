package main

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NodeInfo contiene le informazioni sugli indirizzi IP di un nodo
type NodeInfo struct {
	Name       string
	PodName    string
	Namespace  string
	InternalIP string
	ExternalIP string
}

// GetNodeIPs recupera e restituisce i nomi dei nodi e i loro indirizzi IP
func GetNodeIPs(kubeconfig *string, labelSelector *string, namespace *string) ([]NodeInfo, error) {

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// Crea un clientset Kubernetes
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err.Error())
	}

	// Ottieni i pod con l'etichetta specificata
	pods, err := clientset.CoreV1().Pods(*namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: *labelSelector,
	})
	if err != nil {
		log.Fatalf("Error listing pods: %s", err.Error())
	}

	fmt.Printf("Found %d pods\n", len(pods.Items))

	// Mappa dei nodi per evitare duplicati
	var nodeInfos []NodeInfo
	for _, pod := range pods.Items {
		fmt.Printf("Pod: %s\n", pod.Name)

		nodeName := pod.Spec.NodeName
		// Ottieni i dettagli del nodo
		node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
		if err != nil {
			log.Printf("error getting node %s: %s", nodeName, err)
			continue
		}

		nodeInfo := NodeInfo{Name: nodeName, PodName: pod.Name, Namespace: pod.Namespace}
		for _, address := range node.Status.Addresses {
			if address.Type == corev1.NodeInternalIP {
				nodeInfo.InternalIP = address.Address
			} else if address.Type == corev1.NodeExternalIP {
				nodeInfo.ExternalIP = address.Address
			}
		}
		nodeInfos = append(nodeInfos, nodeInfo)
	}

	return nodeInfos, nil
}

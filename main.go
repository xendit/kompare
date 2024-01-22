package main

import (
	"fmt"
	"kompare/kubernetes"
	"kompare/kubernetes/query"
	"os"
	"path"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeconfigPath := path.Join(homeDir, ".kube", "config")
	client, err := kubernetes.NewKubernetesClient(kubeconfigPath)
	if err != nil {
		panic(err)
	}

	namespaces, err := query.ListNamespaces(client)
	if err != nil {
		panic(err)
	}

	for _, ns := range namespaces.Items {
		fmt.Printf("namespace: %s\n---\n", ns.Name)
		deployments, err := query.ListDeployments(client, ns.Name)
		if err != nil {
			panic(err)
		}
		for _, d := range deployments.Items {
			fmt.Printf("- %s\n", d.Name)
		}
	}

}

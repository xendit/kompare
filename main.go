package main

import (
	"fmt"
	"kompare/kubernetes"
	"kompare/kubernetes/diff"
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

	client.SwitchContext("trident-staging-0")
	namespacesSource, err := query.ListNamespaces(client)
	if err != nil {
		panic(err)
	}

	client.SwitchContext("trident-playground-0")
	namespacesTarget, err := query.ListNamespaces(client)
	if err != nil {
		panic(err)
	}

	namespaceDiff := diff.GetNamespaceDiffByName(namespacesSource, namespacesTarget)
	fmt.Println(namespaceDiff)

}

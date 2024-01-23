package main

import (
	"fmt"
	"kompare/compare"
	"kompare/connect"
	"kompare/query"
)

func main() {
	configFile := "/Users/abel.guzman/.kube/config"
	clientsetToSource, err := connect.ConnectNow(&configFile)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}
	nameSpacesList, err := query.ListNameSpaces(clientsetToSource)
	if err != nil {
		fmt.Printf("Error getting namespace list: %v\n", err)
		return
	}
	// fmt.Println(x)
	for _, ns := range nameSpacesList.Items {
		deploymentListOnSource, err := query.ListK8sDeployments(clientsetToSource, ns.Name)
		if err != nil {
			fmt.Printf("Error getting source cluster's deployment list: %v\n", err)
			return
		}
		// If you need to switch context
		clientsetToTarget, err := connect.ContextSwitch("arn:aws:eks:ap-southeast-1:705506614808:cluster/trident-playground-0", &configFile)
		if err != nil {
			fmt.Printf("Error switching context: %v\n", err)
			return
		}
		deploymentListOnTarget, err := query.ListK8sDeployments(clientsetToTarget, ns.Name)
		if err != nil {
			fmt.Printf("Error getting target cluster's deployment list: %v\n", err)
			return
		}
		// fmt.Println(query.ListK8sDeployments(x, "default"))
		// fmt.Println(query.ListNameSpaces(x))
		// compare.IterateSimpleDiff(deploymentListOnSource, deploymentListOnTarget)
		compare.DeepDeploySourceTargetCompare(deploymentListOnSource, deploymentListOnTarget)
		// fmt.Println(x)
	}

}

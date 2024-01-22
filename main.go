package main

import (
	"kompare/compare"
	"kompare/connect"
	"kompare/query"
)

func main() {
	configFile := "/Users/abel.guzman/.kube/config"
	clientsetToSource := connect.ConnectNow(&configFile)
	nameSpacesList := query.ListNameSpaces(clientsetToSource)
	// fmt.Println(x)
	for _, n := range nameSpacesList.Items {
		deploymentListOnSource := query.ListK8sDeployments(clientsetToSource, n.Name)
		// If you need to switch context
		clientsetToTarget := connect.ContextSwitxh("arn:aws:eks:ap-southeast-1:705506614808:cluster/trident-playground-0", &configFile)
		deploymentListOnTarget := query.ListK8sDeployments(clientsetToTarget, n.Name)
		// fmt.Println(query.ListK8sDeployments(x, "default"))
		// fmt.Println(query.ListNameSpaces(x))
		// compare.IterateSimpleDiff(deploymentListOnSource, deploymentListOnTarget)
		compare.DeepDeploySourceTargetCompare(deploymentListOnSource, deploymentListOnTarget)
		// fmt.Println(x)
	}

}

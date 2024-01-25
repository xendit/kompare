package main

import (
	"fmt"
	"kompare/compare"
	"kompare/connect"
	"kompare/query"
)

func main() {
	// fmt.Println(compare.ConvertTypeStringToHumanReadable("*Corev1.ServiceAccountList"))
	configFile := "/Users/abel.guzman/.kube/config"
	clientsetToSource, err := connect.ConnectNow(&configFile)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}
	nameSpacesList, err := query.ListNameSpaces(clientsetToSource)

	// sourceDeploymentList, err := query.ListK8sDeployments(clientsetToSource, "default")
	if err != nil {
		fmt.Printf("Error getting namespace list: %v\n", err)
		return
	}
	// cluster1Count := compare.GenericCountListElements(nameSpacesList1)
	// clientsetTotarget, err := connect.ContextSwitch("arn:aws:eks:ap-southeast-1:705506614808:cluster/trident-playground-0", &configFile)
	// nameSpacesList2, err := query.ListK8sDeployments(clientsetTotarget, "default")
	// if err != nil {
	// 	fmt.Printf("Error getting namespace list: %v\n", err)
	// 	return
	// }
	// cluster2Count := compare.GenericCountListElements(nameSpacesList2)
	// fmt.Println(cluster1Count, cluster2Count)
	// compare.CompareNumbersGenericOutput(cluster1Count, cluster2Count, nameSpacesList2)
	// compare.IterateGenericSimpleDiff(nameSpacesList1, nameSpacesList2)
	// messageTemplate1 := "- First cluster has %s %s, but it's not in the second cluster\n"
	// compare.CompareByName(nameSpacesList1, nameSpacesList2, messageTemplate1)
	// messageTemplate2 := "- Second cluster has %s %s, but it's not in the second cluster\n"
	// compare.CompareByName(nameSpacesList2, nameSpacesList1, messageTemplate2)
	// fieldsToCompre := []string{"Spec", "Status"}
	// fieldsToCompre := []string{"Spec.Template.Spec"}
	// fmt.Println(compare.DeepCompare(nameSpacesList1, nameSpacesList2, fieldsToCompre))
	// compare.CompareByName(nameSpacesList2, nameSpacesList1, "- Second cluster has deployment %s, but it's not in the first cluster\n")
	// compare.GenericCount()
	// if Object1 == Object2 {
	// 	compare.CompareNumkersGenericOutput(Object1, cluster1, cluster2)
	// }
	for _, ns := range nameSpacesList.Items {
		fmt.Println("Looping on NS: ", ns.Name)
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
		fieldsToCompre := []string{"Spec.Template.Spec"}
		results := compare.DeepCompare(deploymentListOnSource, deploymentListOnTarget, fieldsToCompre)
		if len(results) > 0 {
			fmt.Println("Diff:", results)
		}
		// 	// compare.IterateDeploymentsSimpleDiff(deploymentListOnSource, deploymentListOnTarget)
		// 	compare.DeepDeploySourceTargetCompare(deploymentListOnSource, deploymentListOnTarget)
		// }
		// fmt.Println(query.ListCRDs(configFile))
	}
}

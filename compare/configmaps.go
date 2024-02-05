package compare

import (
	"fmt"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareConfigMaps(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]DiffWithName, error) {

	var TheDiff []DiffWithName
	sourceConfigMaps, err := query.ListConfigMaps(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	targetConfigMaps, err := query.ListConfigMaps(clientsetToTarget, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Data", "Name", "Annotations"}
	return CompareVerboseVSNonVerbose(sourceConfigMaps, targetConfigMaps, diffCriteria, boolverboseDiffs)
}

func GenericCompareConfigMaps(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]DiffWithName, error) {
	diffCriteria := []string{"Data", "Name", "Annotations"}
	return GenericCompareResources(clientsetToSource, clientsetToTarget, namespaceName, query.ListConfigMapsGeneric, diffCriteria, boolverboseDiffs)
}

// func CompareConfigMaps(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]DiffWithName, error) {
// 	diffCriteria := []string{"Data", "Name", "Annotations"}
// 	return CompareResources(clientsetToSource, clientsetToTarget, namespaceName, query.ListConfigMaps, diffCriteria, boolverboseDiffs)
// }

package compare

import (
	"fmt"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareClusterRoles(clientsetToSource, clientsetToTarget *kubernetes.Clientset, boolverboseDiffs *bool) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	sourceClusterRoles, err := query.ListClusterRoles(clientsetToSource)
	if err != nil {
		fmt.Printf("Error getting cluster role list: %v\n", err)
		return TheDiff, err
	}
	targetClusterRoles, err := query.ListClusterRoles(clientsetToTarget)
	if err != nil {
		fmt.Printf("Error getting cluster role list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Rules", "Name", "Annotations"}
	return CompareVerboseVSNonVerbose(sourceClusterRoles, targetClusterRoles, diffCriteria, boolverboseDiffs)
}

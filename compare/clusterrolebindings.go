package compare

import (
	"fmt"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareClusterRoleBindings(clientsetToSource, clientsetToTarget *kubernetes.Clientset, boolverboseDiffs *bool) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	sourceClusterRoleBindings, err := query.ListClusterRoleBindings(clientsetToSource)
	if err != nil {
		fmt.Printf("Error getting cluster role list: %v\n", err)
		return TheDiff, err
	}
	targetClusterRoleBindings, err := query.ListClusterRoleBindings(clientsetToTarget)
	if err != nil {
		fmt.Printf("Error getting cluster role list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"RoleRef", "Name", "Annotations"}
	return CompareVerboseVSNonVerbose(sourceClusterRoleBindings, targetClusterRoleBindings, diffCriteria, boolverboseDiffs)
}

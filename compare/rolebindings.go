package compare

import (
	"fmt"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareRoleBindings(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]DiffWithName, error) {

	var TheDiff []DiffWithName
	sourceRoleBindings, err := query.ListRoleBindings(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting role bindings list: %v\n", err)
		return TheDiff, err
	}
	targetRoleBindings, err := query.ListRoleBindings(clientsetToTarget, namespaceName)
	if err != nil {
		fmt.Printf("Error getting role bindings list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"RoleRef", "Subjects", "Annotations"}
	return CompareVerboseVSNonVerbose(sourceRoleBindings, targetRoleBindings, diffCriteria, boolverboseDiffs)
}

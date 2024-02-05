package compare

import (
	"fmt"
	"kompare/cli"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareClusterRoleBindings(clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
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
	return CompareVerboseVSNonVerbose(sourceClusterRoleBindings, targetClusterRoleBindings, diffCriteria, &TheArgs.VerboseDiffs)
}

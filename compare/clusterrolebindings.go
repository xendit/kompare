package compare

import (
	"fmt"
	"kompare/DAO"
	"kompare/cli"
	"kompare/query"
	"kompare/tools"

	"k8s.io/client-go/kubernetes"
)

func CompareClusterRoleBindings(clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) ([]DAO.DiffWithName, error) {
	var TheDiff []DAO.DiffWithName
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
	var diffCriteria []string
	if TheArgs.FiltersForObject == "" {
		diffCriteria = []string{"RoleRef", "Name", "Annotations"}
	} else {
		diffCriteria = tools.ParseCommaSeparateList(TheArgs.FiltersForObject)
	}
	return CompareVerboseVSNonVerbose(sourceClusterRoleBindings, targetClusterRoleBindings, diffCriteria, TheArgs)
}

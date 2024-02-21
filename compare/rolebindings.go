package compare

import (
	"fmt"
	"kompare/DAO"
	"kompare/cli"
	"kompare/query"
	"kompare/tools"

	"k8s.io/client-go/kubernetes"
)

func CompareRoleBindings(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, TheArgs cli.ArgumentsReceivedValidated) ([]DAO.DiffWithName, error) {

	var TheDiff []DAO.DiffWithName
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
	var diffCriteria []string
	if TheArgs.FiltersForObject == "" {
		diffCriteria = []string{"RoleRef", "Subjects"}
	} else {
		diffCriteria = tools.ParseCommaSeparateList(TheArgs.FiltersForObject)
	}
	return CompareVerboseVSNonVerbose(sourceRoleBindings, targetRoleBindings, diffCriteria, TheArgs)
}

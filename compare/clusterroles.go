package compare

import (
	"fmt"
	"kompare/cli"
	"kompare/query"
	"kompare/tools"

	"k8s.io/client-go/kubernetes"
)

func CompareClusterRoles(clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
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
	var diffCriteria []string
	if TheArgs.FiltersForObject == "" {
		diffCriteria = []string{"Rules", "Name", "Annotations"}
	} else {
		diffCriteria = tools.ParseCommaSeparateList(TheArgs.FiltersForObject)
	}
	return CompareVerboseVSNonVerbose(sourceClusterRoles, targetClusterRoles, diffCriteria, TheArgs)
}

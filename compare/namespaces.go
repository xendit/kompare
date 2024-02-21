package compare

import (
	"fmt"
	"kompare/DAO"
	"kompare/cli"
	"kompare/query"
	"kompare/tools"

	"k8s.io/client-go/kubernetes"
)

// Compare actual namespaces comparison using generic functions from module "compare"
func CompareNameSpaces(clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) ([]DAO.DiffWithName, error) {
	var TheDiff []DAO.DiffWithName
	sourceNameSpacesList, err := query.ListNameSpaces(clientsetToSource)
	if err != nil {
		fmt.Printf("Error getting namespace list: %v\n", err)
		return TheDiff, err
	}
	targetNameSpacesList, err := query.ListNameSpaces(clientsetToTarget)
	if err != nil {
		fmt.Printf("Error getting namespace list: %v\n", err)
		return TheDiff, err
	}
	var diffCriteria []string
	if TheArgs.FiltersForObject == "" {
		diffCriteria = []string{"Spec", "Name", "Status.Phase"}
	} else {
		diffCriteria = tools.ParseCommaSeparateList(TheArgs.FiltersForObject)
	}
	return CompareVerboseVSNonVerbose(sourceNameSpacesList, targetNameSpacesList, diffCriteria, TheArgs)
}

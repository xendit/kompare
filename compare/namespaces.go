package compare

import (
	"fmt"
	"kompare/cli"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

// Compare actual namespaces comparison using generic functions from module "compare"
func CompareNameSpaces(clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
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
	diffCriteria := []string{"Spec", "Name", "Status.Phase"}
	return CompareVerboseVSNonVerbose(sourceNameSpacesList, targetNameSpacesList, diffCriteria, &TheArgs.VerboseDiffs)
}

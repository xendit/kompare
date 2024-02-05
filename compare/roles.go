package compare

import (
	"fmt"
	"kompare/cli"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareRoles(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, TheArgs cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {

	var TheDiff []DiffWithName
	sourceRoles, err := query.ListRoles(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting roles list: %v\n", err)
		return TheDiff, err
	}
	targetRoles, err := query.ListRoles(clientsetToTarget, namespaceName)
	if err != nil {
		fmt.Printf("Error getting roles list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Rules", "Name"}
	return CompareVerboseVSNonVerbose(sourceRoles, targetRoles, diffCriteria, &TheArgs.VerboseDiffs)
}

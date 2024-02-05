package compare

import (
	"fmt"
	"kompare/cli"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareServices(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, TheArgs cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	sourceServices, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	targetServices, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Spec", "Name"}
	return CompareVerboseVSNonVerbose(sourceServices, targetServices, diffCriteria, &TheArgs.VerboseDiffs)
}

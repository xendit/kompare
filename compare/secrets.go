package compare

import (
	"fmt"
	"kompare/cli"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareSecrets(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, TheArgs cli.ArgumentsReceivedValidated) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	sourceSecrets, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting secrets list: %v\n", err)
		return TheDiff, err
	}
	targetSecrets, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting secrets list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Annotations", "Name"}
	return CompareVerboseVSNonVerbose(sourceSecrets, targetSecrets, diffCriteria, &TheArgs.VerboseDiffs)
}

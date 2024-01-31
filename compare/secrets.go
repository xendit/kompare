package compare

import (
	"fmt"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareSecrets(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]DiffWithName, error) {
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
	return CompareVerboseVSNonVerbose(sourceSecrets, targetSecrets, diffCriteria, boolverboseDiffs)
}

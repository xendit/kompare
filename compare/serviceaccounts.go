package compare

import (
	"fmt"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func CompareServiceAccounts(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	sourceServiceAccounts, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting service accounts list: %v\n", err)
		return TheDiff, err
	}
	targetServiceAccounts, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting service accounts list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Annotations", "Name"}
	return CompareVerboseVSNonVerbose(sourceServiceAccounts, targetServiceAccounts, diffCriteria, boolverboseDiffs)
}

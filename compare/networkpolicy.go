package compare

import (
	"fmt"
	"kompare/DAO"
	"kompare/cli"
	"kompare/query"
	"kompare/tools"

	"k8s.io/client-go/kubernetes"
)

func CompareNetworkPolicies(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, TheArgs cli.ArgumentsReceivedValidated) ([]DAO.DiffWithName, error) {
	var TheDiff []DAO.DiffWithName
	sourceNetworkPolicies, err := query.ListNetworkPolicies(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	targetNetworkPolicies, err := query.ListNetworkPolicies(clientsetToTarget, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	var diffCriteria []string
	if TheArgs.FiltersForObject == "" {
		diffCriteria = []string{"Spec", "Name", "Annotations"}
	} else {
		diffCriteria = tools.ParseCommaSeparateList(TheArgs.FiltersForObject)
	}
	return CompareVerboseVSNonVerbose(sourceNetworkPolicies, targetNetworkPolicies, diffCriteria, TheArgs)
}

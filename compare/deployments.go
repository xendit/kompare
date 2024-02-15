package compare

import (
	"fmt"
	"kompare/DAO"
	"kompare/cli"
	"kompare/query"
	"kompare/tools"

	"k8s.io/client-go/kubernetes"
)

// compare deployments for a namespace
func CompareDeployments(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, TheArgs cli.ArgumentsReceivedValidated) ([]DAO.DiffWithName, error) {
	var TheDiff []DAO.DiffWithName
	sourceDeployments, err := query.ListDeployments(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	targetDeplotments, err := query.ListDeployments(clientsetToTarget, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	var diffCriteria []string
	if TheArgs.FiltersForObject == "" {
		diffCriteria = []string{"Spec.Template.Spec", "Name"}
	} else {
		diffCriteria = tools.ParseCommaSeparateList(TheArgs.FiltersForObject)
	}
	return CompareVerboseVSNonVerbose(sourceDeployments, targetDeplotments, diffCriteria, TheArgs)
}

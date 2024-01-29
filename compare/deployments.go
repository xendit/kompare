package compare

import (
	"fmt"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

// compare deployments for a namespace
func CompareDeployments(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]DiffWithName, error) {
	var TheDiff []DiffWithName
	sourceDeployments, err := query.ListK8sDeployments(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	targetDeplotments, err := query.ListK8sDeployments(clientsetToTarget, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Spec.Template.Spec", "Name"}
	return CompareVerboseVSNonVerbose(sourceDeployments, targetDeplotments, diffCriteria, boolverboseDiffs)
}

package compare

import (
	"fmt"
	"kompare/DAO"
	"kompare/cli"
	"kompare/query"
	"kompare/tools"

	"k8s.io/client-go/kubernetes"
)

func CompareCronJobs(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, TheArgs cli.ArgumentsReceivedValidated) ([]DAO.DiffWithName, error) {
	var TheDiff []DAO.DiffWithName
	sourceCronJobs, err := query.ListCronJobs(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	targetCronJobs, err := query.ListCronJobs(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	var diffCriteria []string
	if TheArgs.FiltersForObject == "" {
		diffCriteria = []string{"Spec", "Name"}
	} else {
		diffCriteria = tools.ParseCommaSeparateList(TheArgs.FiltersForObject)
	}
	return CompareVerboseVSNonVerbose(sourceCronJobs, targetCronJobs, diffCriteria, TheArgs)
}

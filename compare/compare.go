package compare

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/apps/v1"
)

// source vs target or old vs new
func CompareNumberOfDeployments(sourceDeployments, targetDeplotments *v1.DeploymentList) (int, int) {
	//Print quantity of deployments
	len_sourceDeployments := len(sourceDeployments.Items)
	len_targetDeplotments := len(targetDeplotments.Items)
	fmt.Printf("There are %d Deployments(apps) in the source cluster and %d in the target Cluster\n",
		len_sourceDeployments, len_targetDeplotments)
	// IF deployment quantities in both clusters are different find those diferent apps in a later function
	return len_sourceDeployments, len_targetDeplotments
}

func iterateSimpleDiff(sourceDeployments, targetDeplotments *v1.DeploymentList) {
	len_sourceDeployments, len_targetDeplotments := CompareNumberOfDeployments(sourceDeployments, targetDeplotments)
	if len_sourceDeployments != len_targetDeplotments {
		fmt.Printf("NOTICE NOT EQUAL NUMBER OF DEPLOYMENTS!!!\n")
		compareDeploymentsByName(sourceDeployments, targetDeplotments,
			"- First cluster has deployment %s, but it's not in the second cluster\n")
		compareDeploymentsByName(targetDeplotments, sourceDeployments,
			"- Second cluster has deployment %s, but it's not in the first cluster\n")
	} else {
		deepDeployFirstSecondCompare(targetDeplotments, sourceDeployments, show_replicas)
	}
}

func compareDeploymentsByName(first_deployments, second_deployments *v1.DeploymentList, message_heading string) {
	for _, d := range first_deployments.Items {
		exists := false
		for _, b := range second_deployments.Items {
			if b.Name == d.Name {
				exists = true
			}
		}
		if exists == false {
			fmt.Printf(strings.Replace(message_heading, "%s", d.Name, -1))
		}
	}
}

func deepDeployFirstSecondCompare(sourceDeployments, targetDeplotments *v1.DeploymentList, show_replicas *bool) {

}

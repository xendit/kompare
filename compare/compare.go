package compare

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-test/deep"

	v1 "k8s.io/api/apps/v1"
)

// CompareNumberOfDeployments compares the number of deployments in the source and target clusters.
// It takes two DeploymentList objects as input.
// It prints the number of deployments in the source and target clusters.
// It returns the number of deployments in the source and target clusters.
func CompareNumberOfDeployments(sourceDeployments, targetDeplotments *v1.DeploymentList) (int, int) {
	// Print quantity of deployments
	lenSourceDeployments := len(sourceDeployments.Items)
	lenTargetDeplotments := len(targetDeplotments.Items)
	fmt.Printf("There are %d Deployments(apps) in the source cluster and %d in the target cluster\n",
		lenSourceDeployments, lenTargetDeplotments)
	// If deployment quantities in both clusters are different, find those different apps in a later function
	return lenSourceDeployments, lenTargetDeplotments
}

// IterateDeploymentsSimpleDiff iterates over the deployments in the source and target clusters.
// It calls CompareNumberOfDeployments to get the number of deployments in each cluster.
// If the number of deployments is not equal, it identifies deployments that are only present in one cluster but not the other.
// It prints a notice about the unequal number of deployments.
// It returns two lists: onlyInSource and onlyInTarget, which contain the names of deployments that are only present in the source or target cluster, respectively.
func IterateDeploymentsSimpleDiff(sourceDeployments, targetDeplotments *v1.DeploymentList) ([]string, []string) {
	lenSourceDeployments, lenTargetDeplotments := CompareNumberOfDeployments(sourceDeployments, targetDeplotments)
	if lenSourceDeployments != lenTargetDeplotments {
		var onlyInSource, onlyInTarget []string

		fmt.Printf("NOTICE NOT EQUAL NUMBER OF DEPLOYMENTS!!!\n")
		onlyInSource = compareDeploymentsByName(sourceDeployments, targetDeplotments,
			"- Source cluster has deployment %s, but it's not in the target cluster\n")
		onlyInTarget = compareDeploymentsByName(targetDeplotments, sourceDeployments,
			"- target cluster has deployment %s, but it's not in the source cluster\n")
		return onlyInSource, onlyInTarget
	}
	return nil, nil
}

// compareDeploymentsByName compares two lists of deployments by name and returns a list of names that exist in the first list but not in the second list.
//
// Parameters:
//
//	first_deployments: List of deployments to compare against the second list
//	second_deployments: List of deployments to compare against the first list
//	message_heading: A string that will be used to print a message when a deployment in the first list is not found in the second list
//
// Returns:
//
//	diffNameList: List of names of deployments that exist in the first list but not in the second list
func compareDeploymentsByName(first_deployments, second_deployments *v1.DeploymentList, message_heading string) []string {
	var diffNameList []string
	for _, d := range first_deployments.Items {
		exists := false
		for _, b := range second_deployments.Items {
			if b.Name == d.Name {
				exists = true
			}
		}
		if exists == false {
			fmt.Printf(strings.Replace(message_heading, "%s", d.Name, -1))
			diffNameList = append(diffNameList, d.Name)
		}
	}
	return diffNameList
}

// DeepDeploySourceTargetCompare compares the important parts of the manifest of deployments from the source and target clusters.
//
// Parameters:
//
//	sourceDeployments: List of deployments from the source cluster
//	targetDeployments: List of deployments from the target cluster
//
// Returns:
//
//	diffSourceTarget: List of DiffWithName structs containing the deployments that exist in both clusters but have differences in their specifications.
func DeepDeploySourceTargetCompare(sourceDeployments, targetDeployments *v1.DeploymentList) []DiffWithName {
	var tmpDiff DiffWithName
	var diffSourceTarget []DiffWithName
	for _, d := range sourceDeployments.Items {
		for _, b := range targetDeployments.Items {
			if b.Name == d.Name {
				fmt.Println("Comparing " + b.Name + " on both source and target cluster.")
				if !reflect.DeepEqual(d.Spec.Template.Spec, b.Spec.Template.Spec) {
					fmt.Println("Deployment " + b.Name + " exists on both clusters, but it's different")
					if diff := deep.Equal(d.Spec.Template.Spec, b.Spec.Template.Spec); diff != nil {
						fmt.Println("Diff:")
						fmt.Println(diff)
						tmpDiff.Name = b.Name
						tmpDiff.Namespace = b.Namespace
						tmpDiff.Diff = diff
						diffSourceTarget = append(diffSourceTarget, tmpDiff)
					}
				}
				fmt.Println("Done.")
			}
		}
	}
	return diffSourceTarget
}

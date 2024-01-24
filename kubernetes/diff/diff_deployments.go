package diff

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/go-test/deep"
	v1 "k8s.io/api/apps/v1"
)

// CompareNumberOfDeployments compares the number of deployments in the source and target clusters.
// It takes two DeploymentList objects as input.
// It prints the number of deployments in the source and target clusters.
// It returns the number of deployments in the source and target clusters.
func CompareNumberOfDeployments(sourceDeployments, targetDeployments *v1.DeploymentList) (int, int) {
	// Print quantity of deployments
	lenSourceDeployments := len(sourceDeployments.Items)
	lenTargetDeployments := len(targetDeployments.Items)
	fmt.Printf("There are %d Deployments(apps) in the source cluster and %d in the target cluster\n",
		lenSourceDeployments, lenTargetDeployments)

	// If deployment quantities in both clusters are different, find those different apps in a later function
	return lenSourceDeployments, lenTargetDeployments
}

// IterateDeploymentsSimpleDiff iterates over the deployments in the source and target clusters.
// It calls CompareNumberOfDeployments to get the number of deployments in each cluster.
// If the number of deployments is not equal, it identifies deployments that are only present in one cluster but not the other.
// It prints a notice about the unequal number of deployments.
// It returns two lists: onlyInSource and onlyInTarget, which contain the names of deployments that are only present in the source or target cluster, respectively.
func IterateDeploymentsSimpleDiff(sourceDeployments, targetDeployments *v1.DeploymentList) ([]string, []string) {
	lenSourceDeployments, lenTargetDeployments := CompareNumberOfDeployments(sourceDeployments, targetDeployments)
	if lenSourceDeployments != lenTargetDeployments {
		var onlyInSource, onlyInTarget []string

		fmt.Printf("NOTICE NOT EQUAL NUMBER OF DEPLOYMENTS!!!\n")
		onlyInSource = CompareDeploymentsByName(sourceDeployments, targetDeployments)
		onlyInTarget = CompareDeploymentsByName(targetDeployments, sourceDeployments)
		return onlyInSource, onlyInTarget
	}
	return nil, nil
}

// compareDeploymentsByName compares two lists of deployments by name and returns a list of names that exist in the first list but not in the second list.
//
// Parameters:
//
//	sourceDeployments: List of deployments to compare against the second list
//	targetDeployments: List of deployments to compare against the first list
//
// Returns:
//
//	diffNameList: List of names of deployments that exist in the first list but not in the second list
func CompareDeploymentsByName(sourceDeployments, targetDeployments *v1.DeploymentList) []string {
	var diffNameList []string

	// remap the deployments to a list of names
	sourceDeploymentNames := make([]string, len(sourceDeployments.Items))
	for i, d := range sourceDeployments.Items {
		sourceDeploymentNames[i] = d.Name
	}

	targetDeploymentNames := make([]string, len(targetDeployments.Items))
	for i, d := range targetDeployments.Items {
		targetDeploymentNames[i] = d.Name
	}

	for _, d := range sourceDeployments.Items {

		// search for the deployment name in the target cluster
		idx := sort.SearchStrings(targetDeploymentNames, d.Name)

		// if not found...
		if idx >= len(targetDeploymentNames) {
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

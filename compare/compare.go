package compare

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-test/deep"

	v1 "k8s.io/api/apps/v1"
)

// source vs target or old vs new or First vs Second
func CompareNumberOfDeployments(sourceDeployments, targetDeplotments *v1.DeploymentList) (int, int) {
	//Print quantity of deployments
	lenSourceDeployments := len(sourceDeployments.Items)
	lenTargetDeplotments := len(targetDeplotments.Items)
	fmt.Printf("There are %d Deployments(apps) in the source cluster and %d in the target cluster\n",
		lenSourceDeployments, lenTargetDeplotments)
	// IF deployment quantities in both clusters are different find those diferent apps in a later function
	return lenSourceDeployments, lenTargetDeplotments
}

func IterateSimpleDiff(sourceDeployments, targetDeplotments *v1.DeploymentList) ([]string, []string) {
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

// Compare the important parts of the manifest like images, configmaps, variables defined, optionally replica numbers, etc.
func DeepDeploySourceTargetCompare(sourceDeployments, targetDeplotments *v1.DeploymentList) []DiffWithName {
	var tmpDiff DiffWithName
	var diffSourceTarget []DiffWithName
	for _, d := range sourceDeployments.Items {
		for _, b := range targetDeplotments.Items {
			if b.Name == d.Name {
				fmt.Println("Comparing " + b.Name + " On both source and target cluster.")
				if !reflect.DeepEqual(d.Spec.Template.Spec, b.Spec.Template.Spec) {
					fmt.Println("Deployment " + b.Name + "  exists on both clusters, but it's different")
					if diff := deep.Equal(d.Spec.Template.Spec, b.Spec.Template.Spec); diff != nil {
						fmt.Println("Diff:")
						fmt.Println(diff)
						tmpDiff.Name = b.Name
						tmpDiff.NameSpace = b.Namespace
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

package compare

import (
	"fmt"
	"reflect"
	"strings"

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

// TODO; compare the important parts of the manifest like images, configmaps, variables defined, optionally replica numbers, etc.
func DeepDeploySourceTargetCompare(sourceDeployments, targetDeplotments *v1.DeploymentList) {
	fmt.Println("Deep comparing deployments in foruce and target")
	for _, d := range sourceDeployments.Items {
		for _, b := range targetDeplotments.Items {
			if b.Name == d.Name {
				if !reflect.DeepEqual(d.Spec.Template.Spec, b.Spec.Template.Spec) {
					fmt.Println("Deployment %s exists on both clusters, but it's different", b.Name)
				}
			}
		}
	}
	fmt.Println("We came here!")
}

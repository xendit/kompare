package main

import (
	"fmt"
	"kompare/cli"
	"kompare/compare"
	"kompare/connect"
	"kompare/query"
	"kompare/tools"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// Parse CLI arguments
	args := cli.PaserReader()
	if args.Err != nil {
		err := fmt.Errorf("error parsing arguments: %v", args.Err)
		panic(err)
	}

	// Connect to source cluster
	clientsetToSource, err := connect.ConnectToSource(args.SourceClusterContext, &args.KubeconfigFile)
	if err != nil {
		err = fmt.Errorf("error connecting to source cluster: %v", err)
		panic(err)
	}

	// Connect to target cluster
	clientsetToTarget, err := connect.ContextSwitch(args.TargetClusterContext, &args.KubeconfigFile)
	if err != nil {
		err = fmt.Errorf("error switching context: %v", err)
		panic(err)
	}

	// Determine namespace argument type
	var sourceNameSpacesList *v1.NamespaceList
	var sourceNameSpace *v1.Namespace
	namespaceArgType := DetectNamespacePattern(args.NamespaceName)
	switch namespaceArgType {
	case "specific":
		fmt.Println("Using", args.NamespaceName, "namespace")
		sourceNameSpace, err = query.GetNamespace(clientsetToSource, args.NamespaceName)
		if err != nil {
			err = fmt.Errorf("error listing namespaces: %v", err)
			panic(err)
		}
		sourceNameSpacesList = &v1.NamespaceList{Items: []v1.Namespace{*sourceNameSpace}}
	case "wildcard":
		sourceNameSpacesList, err = query.ListNameSpaces(clientsetToSource)
		if err != nil {
			err = fmt.Errorf("error listing namespaces: %v", err)
			panic(err)
		}
		sourceNameSpacesList = filterNamespaces(sourceNameSpacesList, args.NamespaceName)
	case "empty":
		iterateGoglabObjects(clientsetToSource, clientsetToTarget, args)
		sourceNameSpacesList, err = query.ListNameSpaces(clientsetToSource)
		if err != nil {
			err = fmt.Errorf("error listing namespaces: %v", err)
			panic(err)
		}
	}

	// Iterate over namespaces
	iterateNamespaces(sourceNameSpacesList, clientsetToSource, clientsetToTarget, args)

	fmt.Println("Finished all comparison works!")
}

func iterateGoglabObjects(clientsetToSource, clientsetToTarget *kubernetes.Clientset, args cli.ArgumentsReceivedValidated) bool {
	// Flag to track if any comparison was performed
	comparisonPerformed := false

	// Compare objects based on include list
	if args.Include != nil {
		includeObjects := []string{"namespace", "crd", "clusterrole", "clusterrolebinding"}
		for _, objectType := range includeObjects {
			if tools.IsInList(objectType, args.Include) {
				switch objectType {
				case "namespace":
					_, err := compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, args)
					if err != nil {
						err = fmt.Errorf("error comparing Namespaces: %v", err)
						panic(err)
					}
				case "crd":
					_, err := compare.CompareCRDs(args.TargetClusterContext, args.KubeconfigFile, args)
					if err != nil {
						err = fmt.Errorf("error comparing CRDs: %v", err)
						panic(err)
					}
				case "clusterrole":
					_, err := compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, args)
					if err != nil {
						err = fmt.Errorf("error comparing Cluster Role: %v", err)
						panic(err)
					}
				case "clusterrolebinding":
					_, err := compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, args)
					if err != nil {
						err = fmt.Errorf("error comparing Cluster Role: %v", err)
						panic(err)
					}
				}
				comparisonPerformed = true
			}
		}
	}

	// Compare objects based on exclude list
	if args.Exclude != nil {
		excludeObjects := []string{"namespace", "crd", "clusterrole", "clusterrolebinding"}
		for _, objectType := range excludeObjects {
			if !tools.IsInList(objectType, args.Exclude) {
				switch objectType {
				case "namespace":
					_, err := compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, args)
					if err != nil {
						err = fmt.Errorf("Error comparing Namspace: %v", err)
						panic(err)
					}
				case "crd":
					_, err := compare.CompareCRDs(args.TargetClusterContext, args.KubeconfigFile, args)
					if err != nil {
						err = fmt.Errorf("Error comparing CRDs: %v", err)
						panic(err)
					}
				case "clusterrole":
					_, err := compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, args)
					if err != nil {
						err = fmt.Errorf("Error comparing Cluster Role: %v", err)
						panic(err)
					}
				case "clusterrolebinding":
					_, err := compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, args)
					if err != nil {
						err = fmt.Errorf("Error comparing Cluster Role Binding: %v", err)
						panic(err)
					}
				}
				comparisonPerformed = true
			}
		}
	}

	// If no include or exclude lists are provided, perform default comparisons
	if args.Include == nil && args.Exclude == nil {
		_, err := compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, args)
		if err != nil {
			err = fmt.Errorf("error comparing Namespaces: %v", err)
			panic(err)
		}
		_, err = compare.CompareCRDs(args.TargetClusterContext, args.KubeconfigFile, args)
		if err != nil {
			err = fmt.Errorf("error comparing CRDs: %v", err)
			panic(err)
		}
		_, err = compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, args)
		if err != nil {
			err = fmt.Errorf("error comparing Cluster Roles: %v", err)
			panic(err)
		}
		_, err = compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, args)
		if err != nil {
			err = fmt.Errorf("error comparing Cluster Role Bindings: %v", err)
			panic(err)
		}
		comparisonPerformed = true
	}

	// Print completion message if any comparison was performed
	if comparisonPerformed {
		fmt.Println("Done comparing Kubernetes global objects.")
	}
	return comparisonPerformed
}

func compareAllResourcesInNamespace(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespace string, TheArgs cli.ArgumentsReceivedValidated) {
	fmt.Printf("Looping on Namespace: %s\n", namespace)
	// Compare all resources for the namespace
	resources := []string{"deployment", "ingress", "service", "serviceaccount", "configmap", "secret", "role", "rolebinding", "hpa", "cronjob", "networkpolicy"}

	// Create a title case converter for English
	titleCase := cases.Title(language.English)

	for _, resource := range resources {
		titleResource := titleCase.String(resource)
		fmt.Printf("%s\n", titleResource)
		compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
		fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)
	}

	fmt.Printf("... Done with all resources in ns: %s.\n", namespace)
}

func compareResourcesByLists(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespace string, TheArgs cli.ArgumentsReceivedValidated) {
	fmt.Printf("Looping namespace: %s\n", namespace)

	includeResources := TheArgs.Include
	excludeResources := TheArgs.Exclude

	// Create a title case converter for English
	titleCase := cases.Title(language.English)

	// Define all resources
	resources := []string{"deployment", "ingress", "service", "serviceaccout", "configmap", "secret", "role", "rolebinding", "networkpolicy", "hpa", "cronjob"}

	// Compare resources based on include list
	for _, resource := range includeResources {
		titleResource := titleCase.String(resource)
		fmt.Printf("%s\n", titleResource)
		compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
		fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)

	}

	// Compare resources based on exclude list
	if excludeResources != nil {
		for _, resource := range resources {
			// Check if resource is not in the exclude list
			if !tools.IsInList(resource, excludeResources) {
				titleResource := titleCase.String(resource)
				fmt.Printf("%s\n", titleResource)
				compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
				fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)
			}
		}
	}
}

func compareResource(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespace, resource string, TheArgs cli.ArgumentsReceivedValidated) {
	switch resource {
	case "deployment":
		_, err := compare.CompareDeployments(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Deployments: %v", err)
			panic(err)
		}
	case "ingress":
		_, err := compare.CompareIngresses(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Ingresses: %v", err)
			panic(err)
		}
	case "service":
		_, err := compare.CompareServices(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Services: %v", err)
			panic(err)
		}
	case "serviceaccount":
		_, err := compare.CompareServiceAccounts(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Service Accounts: %v", err)
			panic(err)
		}
	case "configmap":
		_, err := compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Config Maps: %v", err)
			panic(err)
		}
	case "secret":
		_, err := compare.CompareSecrets(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Secrets: %v", err)
			panic(err)
		}
	case "role":
		_, err := compare.CompareRoles(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Roles: %v", err)
			panic(err)
		}
	case "rolebinding":
		_, err := compare.CompareRoleBindings(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Role Bindings: %v", err)
			panic(err)
		}
	case "hpa":
		_, err := compare.CompareHPAs(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Horizontal Pod Autoscalers: %v", err)
			panic(err)
		}
	case "cronjob":
		_, err := compare.CompareCronJobs(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Cron Jobs: %v", err)
			panic(err)
		}
	case "networkpolicy":
		_, err := compare.CompareNetworkPolicies(clientsetToSource, clientsetToTarget, namespace, TheArgs)
		if err != nil {
			err = fmt.Errorf("error comparing Network Policies: %v", err)
			panic(err)
		}
	}
}

func iterateNamespaces(sourceNameSpacesList *v1.NamespaceList, clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) {
	// Check if include or exclude lists are provided, or if no specific lists are provided
	if TheArgs.Include == nil && TheArgs.Exclude == nil {
		// If no include or exclude lists are provided, compare all resources for each namespace
		for _, ns := range sourceNameSpacesList.Items {
			compareAllResourcesInNamespace(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
		}
	} else {
		// Compare resources based on include or exclude lists
		resources := []string{"deployment", "ingress", "service", "serviceaccount", "configmap", "secret", "role", "rolebinding", "hpa", "cronjob", "networkpolicy"}
		if tools.AreAnyInLists(TheArgs.Include, resources) || tools.AreAnyInLists(TheArgs.Exclude, resources) {
			for _, ns := range sourceNameSpacesList.Items {
				compareResourcesByLists(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
			}
		} else {
			fmt.Println("No namespaced resources to compare")
		}
	}
}

// filterNamespaces filters namespaces based on the wildcard pattern
func filterNamespaces(namespaces *v1.NamespaceList, pattern string) *v1.NamespaceList {
	matchingNamespaces := v1.NamespaceList{
		Items: []v1.Namespace{},
	}
	for _, ns := range namespaces.Items {
		if matchWildcard(ns.Name, pattern) {
			matchingNamespaces.Items = append(matchingNamespaces.Items, ns)
		}
	}
	return &matchingNamespaces
}

// matchWildcard checks if a string matches the wildcard pattern
func matchWildcard(s, pattern string) bool {
	match, err := filepath.Match(pattern, s)
	if err != nil {
		return false
	}
	return match
}

func DetectNamespacePattern(pattern string) string {
	if pattern == "" {
		return "empty"
	} else if strings.Contains(pattern, "*") {
		return "wildcard"
	} else {
		return "specific"
	}
}

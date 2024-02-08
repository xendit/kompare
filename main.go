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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// Getting and pasinf CLI argumets.
	TheArgs := cli.PaserReader()
	if TheArgs.Err != nil {
		fmt.Printf("Error parsing arguments: %v\n", TheArgs.Err)
		return
	}
	clientsetToSource, err := connect.ConnectToSource(TheArgs.SourceClusterContext, &TheArgs.KubeconfigFile)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}
	clientsetToTarget, err := connect.ContextSwitch(TheArgs.TargetClusterContext, &TheArgs.KubeconfigFile)
	if err != nil {
		fmt.Printf("Error switching context: %v\n", err)
		return
	}
	// End getting CLI arguments.
	// Global resources.
	var sourceNameSpacesList *v1.NamespaceList
	var sourceNameSpace *v1.Namespace
	NameSpaceArgContent := DetectNamespacePattern(TheArgs.NamespaceName)
	if NameSpaceArgContent == "specific" {
		fmt.Println("Using ", TheArgs.NamespaceName, " namespace")
		sourceNameSpace = &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: TheArgs.NamespaceName}}
		sourceNameSpacesList = &v1.NamespaceList{Items: []v1.Namespace{*sourceNameSpace}}
	} else if NameSpaceArgContent == "wildcard" {
		sourceNameSpacesList, err = query.ListNameSpaces(clientsetToSource)
		sourceNameSpacesList = filterNamespaces(sourceNameSpacesList, TheArgs.NamespaceName)
	} else if NameSpaceArgContent == "empty" {
		iterateGoglabObjects(clientsetToSource, clientsetToTarget, TheArgs)
		sourceNameSpacesList, err = query.ListNameSpaces(clientsetToSource)
		if err != nil {
			fmt.Printf("Error getting namespace list: %v\n", err)
			return
		}
		sourceNameSpace = nil
	}
	//Iterate each namespace.
	iterateNamespaces(sourceNameSpacesList, clientsetToSource, clientsetToTarget, TheArgs)
	fmt.Println("Finished all comparison works!")
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
		for _, ns := range sourceNameSpacesList.Items {
			compareResourcesByLists(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
		}
	}
}

func compareAllResourcesInNamespace(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespace string, TheArgs cli.ArgumentsReceivedValidated) {
	fmt.Printf("Looping on Namespace: %s\n", namespace)

	// Compare all resources for the namespace
	resources := []string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding", "hpa", "cronjob"}
	for _, resource := range resources {
		fmt.Printf("%s\n", strings.Title(resource))
		compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
		fmt.Printf("Finished %s for namespace: %s\n", strings.Title(resource), namespace)
	}

	fmt.Printf("... Done with all resources in ns: %s.\n", namespace)
}

func compareResourcesByLists(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespace string, TheArgs cli.ArgumentsReceivedValidated) {
	fmt.Printf("Looping on NS: %s\n", namespace)

	includeResources := TheArgs.Include
	excludeResources := TheArgs.Exclude

	// Compare resources based on include list
	if includeResources != nil {
		for _, resource := range includeResources {
			fmt.Printf("%s\n", strings.Title(resource))
			compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
			fmt.Printf("Finished %s for namespace: %s\n", strings.Title(resource), namespace)
		}
	}

	// Compare resources based on exclude list
	if excludeResources != nil {
		allResources := []string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding"}
		for _, resource := range allResources {
			if !tools.IsInList(resource, excludeResources) {
				fmt.Printf("%s\n", strings.Title(resource))
				compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
				fmt.Printf("Finished %s for namespace: %s\n", strings.Title(resource), namespace)
			}
		}
	}

	fmt.Printf("... Done with all resources in ns: %s.\n", namespace)
}

func compareResource(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespace, resource string, TheArgs cli.ArgumentsReceivedValidated) {
	switch resource {
	case "deployment":
		compare.CompareDeployments(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "ingress":
		compare.CompareIngresses(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "service":
		compare.CompareServices(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "sa":
		compare.CompareServiceAccounts(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "configmap":
		compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "secret":
		compare.CompareSecrets(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "role":
		compare.CompareRoles(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "rolebinding":
		compare.CompareRoleBindings(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "hpa":
		compare.CompareHPAs(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	case "cronjob":
		compare.CompareCronJobs(clientsetToSource, clientsetToTarget, namespace, TheArgs)
	}
}

func iterateGoglabObjects(clientsetToSource, clientsetToTarget *kubernetes.Clientset, args cli.ArgumentsReceivedValidated) {
	// Flag to track if any comparison was performed
	comparisonPerformed := false

	// Compare objects based on include list
	if args.Include != nil {
		includeObjects := []string{"namespace", "crd", "clusterrole", "clusterrolebinding"}
		for _, objectType := range includeObjects {
			if tools.IsInList(objectType, args.Include) {
				switch objectType {
				case "namespace":
					compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, args)
				case "crd":
					compare.CompareCRDs(args.TargetClusterContext, args.KubeconfigFile, args)
				case "clusterrole":
					compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, args)
				case "clusterrolebinding":
					compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, args)
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
					compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, args)
				case "crd":
					compare.CompareCRDs(args.TargetClusterContext, args.KubeconfigFile, args)
				case "clusterrole":
					compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, args)
				case "clusterrolebinding":
					compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, args)
				}
				comparisonPerformed = true
			}
		}
	}

	// If no include or exclude lists are provided, perform default comparisons
	if args.Include == nil && args.Exclude == nil {
		compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, args)
		compare.CompareCRDs(args.TargetClusterContext, args.KubeconfigFile, args)
		compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, args)
		compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, args)
		comparisonPerformed = true
	}

	// Print completion message if any comparison was performed
	if comparisonPerformed {
		fmt.Println("Done comparing Kubernetes global objects.")
	}
}

// filterNamespaces filters namespaces based on the wildcard pattern
func filterNamespaces(namespaces *v1.NamespaceList, pattern string) *v1.NamespaceList {
	var matchingNamespaces v1.NamespaceList
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

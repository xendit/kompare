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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// Parse CLI arguments
	args := cli.PaserReader()
	if args.Err != nil {
		fmt.Printf("Error parsing arguments: %v\n", args.Err)
		return
	}

	// Connect to source cluster
	clientsetToSource, err := connect.ConnectToSource(args.SourceClusterContext, &args.KubeconfigFile)
	if err != nil {
		fmt.Printf("Error connecting to source cluster: %v\n", err)
		return
	}

	// Connect to target cluster
	clientsetToTarget, err := connect.ContextSwitch(args.TargetClusterContext, &args.KubeconfigFile)
	if err != nil {
		fmt.Printf("Error switching context: %v\n", err)
		return
	}

	// Determine namespace argument type
	var sourceNameSpacesList *v1.NamespaceList
	var sourceNameSpace *v1.Namespace
	namespaceArgType := DetectNamespacePattern(args.NamespaceName)
	switch namespaceArgType {
	case "specific":
		fmt.Println("Using", args.NamespaceName, "namespace")
		sourceNameSpace = &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: args.NamespaceName}}
		sourceNameSpacesList = &v1.NamespaceList{Items: []v1.Namespace{*sourceNameSpace}}
	case "wildcard":
		sourceNameSpacesList, err = query.ListNameSpaces(clientsetToSource)
		if err != nil {
			fmt.Printf("Error listing namespaces: %v\n", err)
			return
		}
		sourceNameSpacesList = filterNamespaces(sourceNameSpacesList, args.NamespaceName)
	case "empty":
		iterateGoglabObjects(clientsetToSource, clientsetToTarget, args)
		sourceNameSpacesList, err = query.ListNameSpaces(clientsetToSource)
		if err != nil {
			fmt.Printf("Error listing namespaces: %v\n", err)
			return
		}
		sourceNameSpace = nil
	}

	// Iterate over namespaces
	iterateNamespaces(sourceNameSpacesList, clientsetToSource, clientsetToTarget, args)

	fmt.Println("Finished all comparison works!")
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
	fmt.Printf("Looping on NS: %s\n", namespace)

	includeResources := TheArgs.Include
	excludeResources := TheArgs.Exclude

	// Create a title case converter for English
	titleCase := cases.Title(language.English)

	// Define all resources
	allResources := []string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding"}

	// Compare resources based on include list
	if includeResources != nil {
		for _, resource := range includeResources {
			titleResource := titleCase.String(resource)
			fmt.Printf("%s\n", titleResource)
			compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
			fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)

		}
	}

	// Compare resources based on exclude list
	if excludeResources != nil {
		for _, resource := range allResources {
			// Check if resource is not in the exclude list
			if !tools.IsInList(resource, excludeResources) {
				titleResource := titleCase.String(resource)
				fmt.Printf("%s\n", titleResource)
				compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
				fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)
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

	fmt.Printf("... Done with all resources in ns: %s.\n", namespace)
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
	fmt.Printf("Looping on NS: %s\n", namespace)

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

	// Create a title case converter for English
	titleCase := cases.Title(language.English)

	// Define all resources
	allResources := []string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding"}

	// Compare resources based on include list
	if includeResources != nil {
		for _, resource := range includeResources {
			titleResource := titleCase.String(resource)
			fmt.Printf("%s\n", titleResource)
			compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
			fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)
		}
	}

	// Compare resources based on exclude list
	if excludeResources != nil {
		for _, resource := range allResources {
			// Check if resource is not in the exclude list
			if !tools.IsInList(resource, excludeResources) {
				titleResource := titleCase.String(resource)
				fmt.Printf("%s\n", titleResource)
				compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
				fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)
			}
		}
	}

	fmt.Printf("... Done with all resources in ns: %s.\n", namespace)
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
	fmt.Printf("Looping on NS: %s\n", namespace)

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

	// Create a title case converter for English
	titleCase := cases.Title(language.English)

	// Define all resources
	allResources := []string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding"}

	// Compare resources based on include list
	if includeResources != nil {
		for _, resource := range includeResources {
			titleResource := titleCase.String(resource)
			fmt.Printf("%s\n", titleResource)
			compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
			fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)
		}
	}

	// Compare resources based on exclude list
	if excludeResources != nil {
		for _, resource := range allResources {
			// Check if resource is not in the exclude list
			if !tools.IsInList(resource, excludeResources) {
				titleResource := titleCase.String(resource)
				fmt.Printf("%s\n", titleResource)
				compareResource(clientsetToSource, clientsetToTarget, namespace, resource, TheArgs)
				fmt.Printf("Finished %s for namespace: %s\n", titleResource, namespace)
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

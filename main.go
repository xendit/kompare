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

func iterateNamespaces(sourceNameSpacesList *v1.NamespaceList, clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) {
	// Comparing resources per namespace (Namespaced resources).
	if (TheArgs.Include == nil && TheArgs.Exclude == nil) || tools.AreAnyInLists([]string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding"}, TheArgs.Include) || tools.AreAnyInLists([]string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding"}, TheArgs.Exclude) {
		for _, ns := range sourceNameSpacesList.Items {
			fmt.Printf("Looping on NS: %s\n", ns.Name)
			if TheArgs.Exclude == nil && TheArgs.Include == nil {
				fmt.Println("Deployments")
				compare.CompareDeployments(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished deployments for namespace: ", ns.Name)
				fmt.Println("Services")
				compare.CompareServices(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Services for namespace: ", ns.Name)
				fmt.Println("Service Accounts")
				compare.CompareServiceAccounts(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Services Accounts for namespace: ", ns.Name)
				fmt.Println("Secrets")
				compare.CompareSecrets(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Secrets for namespace: ", ns.Name)
				fmt.Println("Config Maps (CM)")
				compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Config Maps (CM) for namespace: ", ns.Name)
				fmt.Println("Roles (RBAC)")
				compare.CompareRoles(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Roles (RBAC) for namespace: ", ns.Name)
				fmt.Println("Role Bindings (RBAC)")
				compare.CompareRoleBindings(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Role Bindings (RBAC) for namespace: ", ns.Name)
				fmt.Println("Ingress")
				compare.CompareIngresses(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Ingreess for namespace: ", ns.Name)
				fmt.Println("HPAs")
				compare.CompareHPAs(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished HPAs for namespace: ", ns.Name)
				fmt.Println("Cron Jobs")
				compare.CompareCronJobs(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
				fmt.Println("Finished Cron Jobs for namespace: ", ns.Name)
			}
			if TheArgs.Exclude != nil && tools.AreAnyInLists([]string{"deployment", "ingress", "service", "sa", "configmap", "secret", "role", "rolebinding"}, TheArgs.Exclude) {
				if tools.IsInList("deployment", TheArgs.Exclude) == false {
					fmt.Println("Deployments")
					compare.CompareDeployments(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished deployments for namespace: ", ns.Name)
				}
				// End Deployment
				if tools.IsInList("service", TheArgs.Exclude) == false {
					// - Services (Spec, Metadata.Annotations, Metadata.Labels )
					fmt.Println("Services")
					compare.CompareServices(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Services for namespace: ", ns.Name)
				}
				// End services
				if tools.IsInList("sa", TheArgs.Exclude) == false {
					// - Service accounts (Metadata.Annotations, Metadata.Labels)
					fmt.Println("Service Accounts")
					compare.CompareServiceAccounts(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Services Accounts for namespace: ", ns.Name)
				}
				// End Service accounts
				if tools.IsInList("secret", TheArgs.Exclude) == false {
					// - Secrets (Type, Data?)
					fmt.Println("Secrets")
					compare.CompareSecrets(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Secrets for namespace: ", ns.Name)
				}
				// End Secrets
				// - Config Maps (criteria)
				if tools.IsInList("configmap", TheArgs.Exclude) == false {
					fmt.Println("Config Maps (CM)")
					// compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
					// compare.GenericCompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
					compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Config Maps (CM) for namespace: ", ns.Name)
				}
				// End Config maps
				if tools.IsInList("role", TheArgs.Exclude) == false {
					// - Roles
					fmt.Println("Roles (RBAC)")
					compare.CompareRoles(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Roles (RBAC) for namespace: ", ns.Name)
				}
				// End Roles
				if tools.IsInList("role", TheArgs.Exclude) == false {
					// - Role Bindings
					fmt.Println("Role Bindings (RBAC)")
					compare.CompareRoleBindings(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Role Bindings (RBAC) for namespace: ", ns.Name)
				}
			}
			if TheArgs.Exclude == nil && TheArgs.Include != nil {
				if tools.IsInList("deployment", TheArgs.Include) == true {
					fmt.Println("Deployments")
					compare.CompareDeployments(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished deployments for namespace: ", ns.Name)
				}
				// End Deployment
				if tools.IsInList("service", TheArgs.Include) == true {
					// - Services (Spec, Metadata.Annotations, Metadata.Labels )
					fmt.Println("Services")
					compare.CompareServices(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Services for namespace: ", ns.Name)
				}
				// End services
				if tools.IsInList("sa", TheArgs.Include) == true {
					// - Service accounts (Metadata.Annotations, Metadata.Labels)
					fmt.Println("Service Accounts")
					compare.CompareServiceAccounts(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Services Accounts for namespace: ", ns.Name)
				}
				// End Service accounts
				if tools.IsInList("secret", TheArgs.Include) == true {
					// - Secrets (Type, Data?)
					fmt.Println("Secrets")
					compare.CompareSecrets(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Secrets for namespace: ", ns.Name)
				}
				// End Secrets
				// - Config Maps (criteria)
				if tools.IsInList("configmap", TheArgs.Include) == true {
					fmt.Println("Config Maps (CM)")
					// compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
					// compare.GenericCompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
					compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Config Maps (CM) for namespace: ", ns.Name)
				}
				// End Config maps
				if tools.IsInList("role", TheArgs.Include) == true {
					// - Roles
					fmt.Println("Roles (RBAC)")
					compare.CompareRoles(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Roles (RBAC) for namespace: ", ns.Name)
				}
				// End Roles
				if tools.IsInList("role", TheArgs.Include) == true {
					// - Role Bindings
					fmt.Println("Role Bindings (RBAC)")
					compare.CompareRoleBindings(clientsetToSource, clientsetToTarget, ns.Name, TheArgs)
					fmt.Println("Finished Role Bindings (RBAC) for namespace: ", ns.Name)
				}
			}
			// End Role Bindings
			fmt.Printf("... Done with all resources in ns: %s.\n", ns.Name)
		}
	}
}

func iterateGoglabObjects(clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) {
	// Comparing namespaces.
	// notice that this fuction returns a diff to be used if we use tests instead of CLI
	doSomething := false
	if TheArgs.Include != nil &&
		tools.AreAnyInLists([]string{"namespace", "crd", "clusterrole", "clusterrolebinding"}, TheArgs.Include) {
		if tools.IsInList("namespace", TheArgs.Include) == true {
			compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, TheArgs)
		}
		if tools.IsInList("crd", TheArgs.Include) == true {
			compare.CompareCRDs(TheArgs.TargetClusterContext, TheArgs.KubeconfigFile, TheArgs)
		}
		if tools.IsInList("clusterrole", TheArgs.Include) == true {
			compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, TheArgs)
		}
		if tools.IsInList("clusterrolebinding", TheArgs.Include) == true {
			compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, TheArgs)
		}
		doSomething = true
	}
	if TheArgs.Exclude != nil &&
		tools.AreAnyInLists([]string{"namespace", "crd", "clusterrole", "clusterrolebinding"}, TheArgs.Exclude) {
		if tools.IsInList("namespace", TheArgs.Exclude) == false {
			compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, TheArgs)
		}
		if tools.IsInList("crd", TheArgs.Exclude) == false {
			compare.CompareCRDs(TheArgs.TargetClusterContext, TheArgs.KubeconfigFile, TheArgs)
		}
		if tools.IsInList("clusterrole", TheArgs.Exclude) == false {
			compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, TheArgs)
		}
		if tools.IsInList("clusterrolebinding", TheArgs.Exclude) == false {
			compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, TheArgs)
		}
		doSomething = true
	}
	if TheArgs.Include == nil && TheArgs.Exclude == nil {
		compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, TheArgs)
		compare.CompareCRDs(TheArgs.TargetClusterContext, TheArgs.KubeconfigFile, TheArgs)
		compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, TheArgs)
		compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, TheArgs)
		doSomething = true
	}
	if doSomething {
		fmt.Println("Done comparing Kuberentes global objects.")
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

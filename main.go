package main

import (
	"fmt"
	"kompare/cli"
	"kompare/compare"
	"kompare/connect"
	"kompare/query"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// Getting and pasinf CLI argumets.
	// var TheArgs cli.ArgumentsReceived
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
	var sourceNameSpacesList *v1.NamespaceList
	var sourceNameSpace *v1.Namespace

	if TheArgs.NamespaceName != "" {
		fmt.Println("Using ", TheArgs.NamespaceName, " namespace")
		sourceNameSpace = &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: TheArgs.NamespaceName}}
		sourceNameSpacesList = &v1.NamespaceList{Items: []v1.Namespace{*sourceNameSpace}}
	} else {
		// Comparing namespaces.
		// notice that this fuction returns a diff to be used if we use tests instead of CLI
		compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, &TheArgs.VerboseDiffs)
		// End Namespaces comparison
		// Comparing CRDs
		compare.CompareCRDs(TheArgs.TargetClusterContext, TheArgs.KubeconfigFile, &TheArgs.VerboseDiffs)
		// End comparing CRDs
		// - Cluster Roles (criteria?)
		compare.CompareClusterRoles(clientsetToSource, clientsetToTarget, &TheArgs.VerboseDiffs)
		// End Cluster Roles
		// - clusterrolebindingss (criteria?)
		compare.CompareClusterRoleBindings(clientsetToSource, clientsetToTarget, &TheArgs.VerboseDiffs)
		// End Cluster Role Bindings
		sourceNameSpacesList, err = query.ListNameSpaces(clientsetToSource)

		if err != nil {
			fmt.Printf("Error getting namespace list: %v\n", err)
			return
		}
		sourceNameSpace = nil
	}
	iterateNamespaces(sourceNameSpacesList, clientsetToSource, clientsetToTarget, TheArgs)
	// - Ingress (Needed?)
	// Features (goot to have)
	// - save comparison to file.
	// - Compare file to target again.
	fmt.Println("Finished!")

}

func iterateNamespaces(sourceNameSpacesList *v1.NamespaceList, clientsetToSource, clientsetToTarget *kubernetes.Clientset, TheArgs cli.ArgumentsReceivedValidated) {
	// Comparing resources per namespace (Namespaced resources).
	for _, ns := range sourceNameSpacesList.Items {
		fmt.Printf("Looping on NS: %s", ns.Name)
		// - Deployment (Spec.Template.Spec & ?)
		fmt.Println("Deployments")
		compare.CompareDeployments(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
		fmt.Println("Finished deployments for namespace: ", ns.Name)
		// End Deployment
		// - Services (Spec, Metadata.Annotations, Metadata.Labels )
		fmt.Println("Services")
		compare.CompareServices(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
		fmt.Println("Finished Services for namespace: ", ns.Name)
		// End services
		// - Service accounts (Metadata.Annotations, Metadata.Labels)
		fmt.Println("Service Accounts")
		compare.CompareServiceAccounts(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
		fmt.Println("Finished Services Accounts for namespace: ", ns.Name)
		// End Service accounts
		// - Secrets (Type, Data?)
		fmt.Println("Secrets")
		compare.CompareSecrets(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
		fmt.Println("Finished Secrets for namespace: ", ns.Name)
		// End Secrets
		// - Config Maps (criteria)
		fmt.Println("Config Maps (CM)")
		compare.CompareConfigMaps(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
		fmt.Println("Finished Config Maps (CM) for namespace: ", ns.Name)
		// End Config maps
		// - Roles
		fmt.Println("Roles (RBAC)")
		compare.CompareRoles(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
		fmt.Println("Finished Roles (RBAC) for namespace: ", ns.Name)
		// End Roles
		// - Role Bindings
		fmt.Println("Role Bindings (RBAC)")
		compare.CompareRoleBindings(clientsetToSource, clientsetToTarget, ns.Name, &TheArgs.VerboseDiffs)
		fmt.Println("Finished Role Bindings (RBAC) for namespace: ", ns.Name)
		// End Role Bindings
		fmt.Printf("... Done with all resources in ns: %s.\n", ns.Name)
	}
}

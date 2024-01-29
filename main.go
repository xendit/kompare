package main

import (
	"fmt"
	"kompare/cli"
	"kompare/compare"
	"kompare/connect"
	"kompare/query"
)

func main() {
	// Getting and pasinf CLI argumets.
	var TheArgs cli.ArgumentsReceived
	TheArgs = cli.PaserReader()
	if TheArgs.Err != nil {
		fmt.Printf("Error parsing arguments: %v\n", TheArgs.Err)
		return
	}
	configFile, strSourceClusterContext, strTargetClusterContext, strNamespaceName, boolverboseDiffs, err :=
		cli.ValidateParametersFromParserArgs(TheArgs)
	if strNamespaceName != "default" {
		fmt.Println("Using ", strNamespaceName, " namespace")
	}
	if err != nil {
		fmt.Printf("Error validating command line Arguments: %v\n", err)
		return
	}
	// End getting CLI arguments.
	// Comparing namespaces.

	clientsetToSource, err := connect.ConnectToSource(strSourceClusterContext, &configFile)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return
	}
	clientsetToTarget, err := connect.ContextSwitch(strTargetClusterContext, &configFile)
	if err != nil {
		fmt.Printf("Error switching context: %v\n", err)
		return
	}
	// notice that this fuction returns a diff to be used if we use tests instead of CLI
	// fmt.Println(compareNameSpaces(clientsetToSource, clientsetToTarget))
	compare.CompareNameSpaces(clientsetToSource, clientsetToTarget, boolverboseDiffs)
	// End Namespaces comparison
	// Comparing CRDs
	// fmt.Println(compareCRDs(strTargetClusterContext, configFile))
	compare.CompareCRDs(strTargetClusterContext, configFile, boolverboseDiffs)
	// End comparing CRDs
	// PENDING
	// - roles (criteria?)
	// - clusterroles (criteria?)
	// - rolebindings (criteria?)
	// - clusterrolebindingss (criteria?)
	sourceNameSpacesList, err := query.ListNameSpaces(clientsetToSource)
	if err != nil {
		fmt.Printf("Error getting namespace list: %v\n", err)
		return
	}
	// Comparing resources per namespace (Namespaced resources).
	for _, ns := range sourceNameSpacesList.Items {
		fmt.Printf("Looping on NS: %s", ns.Name)
		// - Deployment (Spec.Template.Spec & ?)
		fmt.Println("Deployments")
		compare.CompareDeployments(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
		fmt.Println("Finished deployments for namespace: ", ns.Name)
		// End Deployment
		// - Services (Spec, Metadata.Annotations, Metadata.Labels )
		fmt.Println("Services")
		compare.CompareServices(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
		fmt.Println("Finished Services for namespace: ", ns.Name)
		// End services
		// - Service accounts (Metadata.Annotations, Metadata.Labels)
		fmt.Println("Service Accounts")
		compare.CompareServiceAccounts(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
		fmt.Println("Finished Services Accounts for namespace: ", ns.Name)
		// End Service accounts
		// - Secrets (Type, Data?)
		fmt.Println("Secrets")
		compare.CompareSecrets(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
		fmt.Println("Finished Secrets for namespace: ", ns.Name)
		// End Secrets
		fmt.Printf("... Done with all resources in ns: %s.\n", ns.Name)
	}

	// - Secrets (Type, Data?)
	// - Ingress (Needed?)
	// Features (goot to have)
	// - save comparison to file.
	// - Compare file to target again.
	fmt.Println("Finished!")
}

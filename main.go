package main

import (
	"fmt"
	"kompare/cli"
	"kompare/compare"
	"kompare/connect"
	"kompare/query"

	"k8s.io/client-go/kubernetes"
)

func ConnectToSource(strSourceClusterContext string, configFile *string) (*kubernetes.Clientset, error) {
	var clientsetToSource *kubernetes.Clientset
	var err error
	if strSourceClusterContext != "" {
		clientsetToSource, err = connect.ContextSwitch(strSourceClusterContext, configFile)
	} else {
		clientsetToSource, err = connect.ConnectNow(configFile)
	}
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		return nil, err
	}
	return clientsetToSource, nil
}

func main() {
	// Getting and pasinf CLI argumets.
	kubeconfigFile, sourceClusterContext, targetClusterContext, namespaceName, verboseDiffs, err := cli.PaserReader()
	if err != nil {
		fmt.Printf("Error parsing arguments: %v\n", err)
		return
	}
	configFile, strSourceClusterContext, strTargetClusterContext, strNamespaceName, boolverboseDiffs, err :=
		cli.ValidateParametersFromParserArgs(kubeconfigFile, sourceClusterContext, targetClusterContext, namespaceName, verboseDiffs)
	if strNamespaceName != "default" {
		fmt.Println("Using ", strNamespaceName, " namespace")
	}
	if err != nil {
		fmt.Printf("Error validating command line Arguments: %v\n", err)
		return
	}
	// End getting CLI arguments.
	// Comparing namespaces.

	clientsetToSource, err := ConnectToSource(strSourceClusterContext, &configFile)
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
	compareNameSpaces(clientsetToSource, clientsetToTarget, boolverboseDiffs)
	// End Namespaces comparison
	// Comparing CRDs
	// fmt.Println(compareCRDs(strTargetClusterContext, configFile))
	compareCRDs(strTargetClusterContext, configFile, boolverboseDiffs)
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
		compareDeployments(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
		fmt.Println("Finished deployments for namespace: ", ns.Name)
		// End Deployment
		// - Services (Spec, Metadata.Annotations, Metadata.Labels )
		fmt.Println("Services")
		compareServices(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
		fmt.Println("Finished Services for namespace: ", ns.Name)
		// End services
		// - Service accounts (Metadata.Annotations, Metadata.Labels)
		fmt.Println("Service Accounts")
		compareServiceAccounts(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
		fmt.Println("Finished Services Accounts for namespace: ", ns.Name)
		// End Service accounts
		// - Secrets (Type, Data?)
		fmt.Println("Secrets")
		compareSecrets(clientsetToSource, clientsetToTarget, ns.Name, boolverboseDiffs)
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

// Compare CRDs using generic functions from module "compare"
func compareCRDs(targetContext, configFile string, boolverboseDiffs *bool) ([]compare.DiffWithName, error) {
	var TheDiff []compare.DiffWithName
	sourceCRDs, err := query.ListCRDs("", configFile)
	if err != nil {
		fmt.Printf("Error getting CRDs list: %v\n", err)
		return TheDiff, err
	}
	targetCRDs, err := query.ListCRDs(targetContext, configFile)
	if err != nil {
		fmt.Printf("Error getting CRDs list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Spec", "Name"}
	return compareVerboseVSNonVerbose(sourceCRDs, targetCRDs, diffCriteria, boolverboseDiffs)
}

// Compare actual namespaces comparison using generic functions from module "compare"
func compareNameSpaces(clientsetToSource, clientsetToTarget *kubernetes.Clientset, boolverboseDiffs *bool) ([]compare.DiffWithName, error) {
	var TheDiff []compare.DiffWithName
	sourceNameSpacesList, err := query.ListNameSpaces(clientsetToSource)
	if err != nil {
		fmt.Printf("Error getting namespace list: %v\n", err)
		return TheDiff, err
	}
	targetNameSpacesList, err := query.ListNameSpaces(clientsetToTarget)
	if err != nil {
		fmt.Printf("Error getting namespace list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Spec", "Name", "Status.Phase"}
	return compareVerboseVSNonVerbose(sourceNameSpacesList, targetNameSpacesList, diffCriteria, boolverboseDiffs)
}

// compare deployments for a namespace
func compareDeployments(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]compare.DiffWithName, error) {
	var TheDiff []compare.DiffWithName
	sourceDeployments, err := query.ListK8sDeployments(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	targetDeplotments, err := query.ListK8sDeployments(clientsetToTarget, namespaceName)
	if err != nil {
		fmt.Printf("Error getting deployments list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Spec.Template.Spec", "Name"}
	return compareVerboseVSNonVerbose(sourceDeployments, targetDeplotments, diffCriteria, boolverboseDiffs)
}

func compareServices(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]compare.DiffWithName, error) {
	var TheDiff []compare.DiffWithName
	sourceServices, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	targetServices, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting services list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Spec", "Name"}
	return compareVerboseVSNonVerbose(sourceServices, targetServices, diffCriteria, boolverboseDiffs)
}

func compareServiceAccounts(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]compare.DiffWithName, error) {
	var TheDiff []compare.DiffWithName
	sourceServiceAccounts, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting service accounts list: %v\n", err)
		return TheDiff, err
	}
	targetServiceAccounts, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting service accounts list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Annotations", "Name"}
	return compareVerboseVSNonVerbose(sourceServiceAccounts, targetServiceAccounts, diffCriteria, boolverboseDiffs)
}
func compareSecrets(clientsetToSource, clientsetToTarget *kubernetes.Clientset, namespaceName string, boolverboseDiffs *bool) ([]compare.DiffWithName, error) {
	var TheDiff []compare.DiffWithName
	sourceSecrets, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting secrets list: %v\n", err)
		return TheDiff, err
	}
	targetSecrets, err := query.ListServices(clientsetToSource, namespaceName)
	if err != nil {
		fmt.Printf("Error getting secrets list: %v\n", err)
		return TheDiff, err
	}
	diffCriteria := []string{"Annotations", "Name"}
	return compareVerboseVSNonVerbose(sourceSecrets, targetSecrets, diffCriteria, boolverboseDiffs)
}

func compareVerboseVSNonVerbose(sourceNameSpacesList, targetNameSpacesList interface{}, diffCriteria []string, boolverboseDiffs *bool) ([]compare.DiffWithName, error) {
	if *boolverboseDiffs {
		TheDiff, err := compare.ShowResourceComparison(sourceNameSpacesList, targetNameSpacesList, diffCriteria)
		fmt.Println(compare.FormatDiffHumanReadable(TheDiff))
		return TheDiff, err
	} else {
		return compare.ShowResourceComparison(sourceNameSpacesList, targetNameSpacesList, diffCriteria)
	}
}

// TODO
// compare globally:
// - DRDs (criteria?) <- Alpha done.
// - Same Namespaces exist in both clusters. <- Alpha done.
// - roles (criteria?)
// - clusterroles (criteria?)
// - rolebindings (criteria?)
// - clusterrolebindingss (criteria?)
// Compare per namespace
// - Deployment (Spec.Template.Spec & ?)
// - Services (Spec, Metadata.Annotations, Metadata.Labels )
// - Service accounts (Metadata.Annotations, Metadata.Labels)
// - Secrets (Type, Data?)
// - Ingress (Needed?)
// Features (goot to have)
// - save comparison to file.
// - Compare file to target again.
// - Service specific particulars to compare; e.g.: when a type of object can have multiple ways
// of defining structure and we need to check some of those that are not always present.
// -

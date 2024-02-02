package cli

import (
	"fmt"
	"kompare/tools"
	"os"
	"path"

	"github.com/akamensky/argparse"
)

type ArgumentsReceived struct {
	KubeconfigFile, SourceClusterContext, TargetClusterContext, NamespaceName, Include, Exclude *string
	VerboseDiffs                                                                                *bool
	Err                                                                                         error
}
type ArgumentsReceivedValidated struct {
	KubeconfigFile, SourceClusterContext, TargetClusterContext, NamespaceName string
	Include, Exclude                                                          []string
	VerboseDiffs                                                              bool
	Err                                                                       error
}

func PaserReader() ArgumentsReceivedValidated {
	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	kubeconfigFile := parser.String("c", "conf", &argparse.Options{Required: false, Help: "Path to the clusters kubeconfig; assume ~/.kube/config if not provided"})
	// Create string flag for clusters. Keep present that the order -f and -s is very important.
	sourceClusterContext := parser.String("s", "src", &argparse.Options{Required: false, Help: "The Source cluster's context. Origin cluster in the comparison (LHS-left hand side)"})
	targetClusterContext := parser.String("d", "dst", &argparse.Options{Required: true, Help: "*The target cluster's context (Required). Cluster used as destination or consequent (RHS - Right hand side)"})
	verboseDiffs := parser.Flag("v", "verbose", &argparse.Options{Help: "Just show me all the diffs too. Notice: the output might be LONG!"})
	IncludeK8sObjects := parser.String("i", "include", &argparse.Options{Help: "List of kubernetes objects names to include, this should be an element or a comma separated list."})
	Excludek8sObjects := parser.String("e", "exclude", &argparse.Options{Help: "List of kubernetes objects to include, this should be an element or a comma separated list."})
	// pass namespace.
	namespaceName := parser.String("n", "namespace", &argparse.Options{Help: "Namespace that needs to be copied. defaults to 'default' namespace"})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return ArgumentsReceivedValidated{
			KubeconfigFile:       "",
			SourceClusterContext: "",
			TargetClusterContext: "",
			NamespaceName:        "",
			Include:              []string{""},
			Exclude:              []string{""},
			VerboseDiffs:         *verboseDiffs,
			Err:                  err}
	}
	TheArgs := ArgumentsReceived{
		KubeconfigFile:       kubeconfigFile,
		SourceClusterContext: sourceClusterContext,
		TargetClusterContext: targetClusterContext,
		NamespaceName:        namespaceName,
		Include:              IncludeK8sObjects,
		Exclude:              Excludek8sObjects,
		VerboseDiffs:         verboseDiffs,
		Err:                  err}
	ArgumentsReceivedValidated := ValidateParametersFromParserArgs(TheArgs)
	return ArgumentsReceivedValidated

}

func ValidateParametersFromParserArgs(TheArgs ArgumentsReceived) ArgumentsReceivedValidated {
	var strSourceClusterContext, strTargetClusterContext, strNamespaceName string
	strSourceClusterContext = *TheArgs.SourceClusterContext
	strTargetClusterContext = *TheArgs.TargetClusterContext
	if strSourceClusterContext == "" {
		fmt.Println("We will use current kubeconfig context as 'source cluster'.")
	} else {
		fmt.Printf("We will use %s kubeconfig context as 'source cluster' or 'origin cluster'.\n", strSourceClusterContext)
	}
	fmt.Printf("We will use %s kubeconfig context as 'target cluster'.\n", strTargetClusterContext)

	strNamespaceName = *TheArgs.NamespaceName
	configFile := ""
	if *TheArgs.KubeconfigFile != "" {
		configFile = *TheArgs.KubeconfigFile
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting the home dir: %v\n", err)
			return ArgumentsReceivedValidated{
				KubeconfigFile: "", SourceClusterContext: "",
				TargetClusterContext: "", NamespaceName: "",
				Include: []string{""}, Exclude: []string{""},
				VerboseDiffs: *TheArgs.VerboseDiffs, Err: nil}
		}
		configFile = path.Join(homeDir, ".kube", "config")
	}
	invalidInclude, includeStr := ValidateKubernetesObjects(tools.ParseCommaSeparateList(*TheArgs.Include))
	invalidExclude, excludeStr := ValidateKubernetesObjects(tools.ParseCommaSeparateList(*TheArgs.Exclude))
	if invalidInclude != nil {
		fmt.Print("You passed some invalid kubernetes object to incldue as a parameter: ", invalidInclude)
		fmt.Println(". The program will try to execute anyways and ignore this")
	}
	if invalidExclude != nil {
		fmt.Println("You passed some invalid kubernetes object to exclude as a parameter: ", invalidInclude)
		fmt.Println(". The program will try to execute anyways and ignore this")
	}
	return ArgumentsReceivedValidated{
		KubeconfigFile:       configFile,
		SourceClusterContext: strSourceClusterContext,
		TargetClusterContext: strTargetClusterContext,
		NamespaceName:        strNamespaceName,
		Include:              includeStr,
		Exclude:              excludeStr,
		VerboseDiffs:         *TheArgs.VerboseDiffs,
		Err:                  nil}
}

// ValidateKubernetesObjects validates the given list of Kubernetes object names
// against a list of valid object names and their aliases
// It returns two slices: invalidObjects and validObjects
func ValidateKubernetesObjects(objects []string) ([]string, []string) {
	validObjects := map[string][]string{
		"deployment":         {"deployment", "deployments", "deploy"},
		"ingress":            {"ingress", "ing"},
		"service":            {"service", "svc", "services"},
		"sa":                 {"sa", "serviceaccount", "serviceaccounts"},
		"configmap":          {"configmap", "configmaps", "cm"},
		"secret":             {"secret", "secrets"},
		"namespace":          {"namespace", "ns", "namespaces"},
		"role":               {"role", "roles"},
		"rolebinding":        {"rolebinding", "rolebindings"},
		"clusterrole":        {"clusterrole", "clusterroles"},
		"clusterrolebinding": {"clusterrolebinding", "clusterrolebindings"},
		"crd":                {"crd", "crds", "customresourcedefinition", "customresourcedefinitions"},
		// Add more valid objects and their aliases as needed
	}

	var invalidObjects []string
	var validObjectsStr []string

	for _, obj := range objects {
		found := false
		for standardName, aliases := range validObjects {
			for _, alias := range aliases {
				if obj == alias {
					found = true
					validObjectsStr = append(validObjectsStr, standardName)
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			invalidObjects = append(invalidObjects, obj)
		}
	}

	return invalidObjects, validObjectsStr
}

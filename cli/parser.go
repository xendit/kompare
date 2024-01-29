package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/akamensky/argparse"
)

type ArgumentsReceived struct {
	KubeconfigFile, SourceClusterContext, TargetClusterContext, NamespaceName *string
	VerboseDiffs                                                              *bool
	Err                                                                       error
}

func PaserReader() ArgumentsReceived {
	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	kubeconfigFile := parser.String("c", "conf", &argparse.Options{Required: false, Help: "Path to the clusters kubeconfig; assume ~/.kube/config if not provided"})
	// Create string flag for clusters. Keep present that the order -f and -s is very important.
	sourceClusterContext := parser.String("s", "src", &argparse.Options{Required: false, Help: "The Source cluster's context"})
	targetClusterContext := parser.String("d", "dst", &argparse.Options{Required: true, Help: "*The target cluster's context (Required)"})
	verboseDiffs := parser.Flag("v", "verbose", &argparse.Options{Help: "Just show me all the diffs too. Notice: the output might be LONG!"})
	// pass namespace.
	namespaceName := parser.String("n", "namespace", &argparse.Options{Help: "Namespace that needs to be copied. defaults to 'default' namespace"})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return ArgumentsReceived{KubeconfigFile: nil, SourceClusterContext: nil, TargetClusterContext: nil, NamespaceName: nil, VerboseDiffs: nil, Err: err}
	}
	return ArgumentsReceived{
		KubeconfigFile:       kubeconfigFile,
		SourceClusterContext: sourceClusterContext,
		TargetClusterContext: targetClusterContext,
		NamespaceName:        namespaceName,
		VerboseDiffs:         verboseDiffs,
		Err:                  err}

}

func ValidateParametersFromParserArgs(TheArgs ArgumentsReceived) (string, string, string, string, *bool, error) {
	var strSourceClusterContext, strTargetClusterContext, strNamespaceName string
	if *TheArgs.SourceClusterContext != "" {
		strSourceClusterContext = *TheArgs.SourceClusterContext
	} else {
		strSourceClusterContext = ""
	}
	strTargetClusterContext = *TheArgs.TargetClusterContext
	strNamespaceName = *TheArgs.NamespaceName
	configFile := ""
	if *TheArgs.KubeconfigFile != "" {
		configFile = *TheArgs.KubeconfigFile
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting the home dir: %v\n", err)
			return "", "", "", "", TheArgs.VerboseDiffs, TheArgs.Err
		}
		configFile = path.Join(homeDir, ".kube", "config")
	}
	if strNamespaceName == "" {
		strNamespaceName = "default"
	}

	return configFile, strSourceClusterContext, strTargetClusterContext, strNamespaceName, TheArgs.VerboseDiffs, nil
}

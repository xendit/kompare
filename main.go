package main

import (
	"fmt"
	"kompare/connect"
	"kompare/query"
)

func main() {
	configFile := "/Users/abel.guzman/.kube/config"
	x := connect.ConnectNow(&configFile)
	fmt.Println(x)
	fmt.Println("Just testing here!")
	// If you need to switch context
	// y := connect.ContextSwitxh("arn:aws:eks:ap-southeast-1:705506614808:cluster/trident-playground-0", &configFile)
	fmt.Println(query.ListK8sDeployments(x, "default"))
	fmt.Println(query.ListNameSpaces(x))

}

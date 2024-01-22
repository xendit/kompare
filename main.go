package main

import (
	"kompare/compare"
	"kompare/connect"
	"kompare/query"
)

func main() {
	configFile := "/Users/abel.guzman/.kube/config"
	x := connect.ConnectNow(&configFile)
	// fmt.Println(x)
	xx := query.ListK8sDeployments(x, "default")
	// If you need to switch context
	// y := connect.ContextSwitxh("arn:aws:eks:ap-southeast-1:705506614808:cluster/trident-playground-0", &configFile)
	y := connect.ContextSwitxh("arn:aws:eks:ap-northeast-2:273217691745:cluster/prod-kr-nihao-eks", &configFile)
	yy := query.ListK8sDeployments(y, "default")
	// fmt.Println(query.ListK8sDeployments(x, "default"))
	// fmt.Println(query.ListNameSpaces(x))
	compare.IterateSimpleDiff(xx, yy)
	compare.DeepDeploySourceTargetCompare(xx, yy)
}

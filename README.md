# Kompare

Kompare is a Go CLI runner to compare two clusters. This software compares two kubernetes cluster using kubeconfig to connect to them and compare existing objects in the two clusters based on flexible criteria passed via command line.

So you need to compare two kubernetes clusters. They are supposed to be completely or partially equal or equivalent. You can find tools that actually do some comparison, but we could not find one that we could give it comparison criteria for all the kubernetes objects we needed to compare. Criteria like ignore Kuberentes resource definition subtypes that are always different like UID, Dates, and other objects content that will never match even when they are the same. To give other examples Filter and compare ONLY namespaces whose name ends with '*-pci' or cimpare only the deployments and config maps and secrets. This is just some of the issues you might face while comparing two kubernetes lcusters. That is the reason why we decided to create this tool. 

## Why do we need Kompare
This CLI tool has been created in the context of having to compare two clusters to determine if they are different so they can be interchangeable or replace each other. Enterprices often prefer to keep a few k8s clusters for the same job or run upgrade this way. The practical/real work we use the tool for is to compare a source cluster that is currently in production with a new cluster that we intend to put in prod to replace it to to work side by side. 
**Notice:** Therefore the source cluster "tends" to be considered the source of truth for the comparison.

## Key terms & help to use the tool.
1. Source cluster refers to the origin cluster. If this was a number comparison then the source cluster would be the "left-hand side" (LHS) or "antecedent".
2. Target cluster refers to the destination cluster. If this was a number comparison then the source cluster would be the "right-hand side" (RHS) or the "second number" or  "consequent".
3. The too has a help that should be self explanatory. you can get to this help by executing the tool with the -h flag.
```
go run main.go -h
usage: print [-h|--help] [-c|--conf "<value>"] [-s|--src "<value>"] -t|--target
             "<value>" [-v|--verbose] [-i|--include "<value>"] [-e|--exclude
             "<value>"] [-n|--namespace "<value>"] [-f|--filter "<value>"]

             Prints provided string to stdout

Arguments:

  -h  --help       Print help information
  -c  --conf       Path to the clusters kubeconfig; assume ~/.kube/config if
                   not provided
  -s  --src        The Source cluster's context. Origin cluster in the
                   comparison (LHS-left hand side)
  -t  --target     *The target cluster's context (Required). Cluster used as
                   destination or consequent (RHS - Right hand side)
  -v  --verbose    -v lists the differences and -vv just shows all the diffs
                   too.
  -i  --include    List of kubernetes objects names to include, this should be
                   an element or a comma separated list.
  -e  --exclude    List of kubernetes objects to include, this should be an
                   element or a comma separated list.
  -n  --namespace  Namespace that needs to be copied. defaults to 'default'
                   namespace. The option also accepts wilcard matching of
                   namespace. E.G.: '*-pci' would match any namespace that ends
                   with -pci. Notice that the '' might be required in some
                   consoles like iterm
  -f  --filter     Filter what parts of the object I want to compare. must be
                   used together with -i option to apply to that type of
                   objects
```
So far only the -t option is required.

## Getting Started

This is a project currently in dev, the best way to use it it to do a go build or to just run it in the console with something like:
```
go run main.go -t some-target-context.
```
### example command:
Compare current contect to "MySecondContext-Cluster" use the -v or verbose option, to see the actual diffs. Then -e option is to exclude in this case services ,roles & role bindings for the namespace called velero.
```
go run main.go -t MySecondContext-Cluster -v -n velero -e svc,role,rolebinding
We will use current kubeconfig context as 'source cluster'.
We will use arn:aws:eks:ap-southeast-3:705506614808:cluster/xnd-jk-stg-aws-2 kubeconfig context as 'target cluster'.
Using  velero  namespace
Looping on NS: velero
Deployments
******************************************************************************************************
Done compering source cluster versus target cluster's  Deployment in the list
Done compering target cluster versus source cluster's  Deployment in the list
Object Name: velero
Namespace: velero
Differences:
- Containers.slice[0].Image: velero/velero:v1.12.2 != velero/velero:v1.9.2
- Containers.slice[0].Args.slice[1]: --uploader-type=restic != <no value>
- Containers.slice[0].LivenessProbe: v1.Probe != <nil pointer>
- Containers.slice[0].ReadinessProbe: v1.Probe != <nil pointer>
- TerminationGracePeriodSeconds: 3600 != 30


Finished deployments for namespace:  velero
Service Accounts
***************************************************************************************************
Done compering source cluster versus target cluster's  Service in the list
Done compering target cluster versus source cluster's  Service in the list

Finished Services Accounts for namespace:  velero
Secrets
***************************************************************************************************
Done compering source cluster versus target cluster's  Service in the list
Done compering target cluster versus source cluster's  Service in the list

Finished Secrets for namespace:  velero
Config Maps (CM)
Finished Config Maps (CM) for namespace:  velero
... Done with all resources in ns: velero.
Finished!
```
**Notice:** The software will assume you want to use your current context as source cluster to compare from (-s or --source to set a different source context). Then the -t option is required and its the destination cluster. 

### Prerequisites

It's a golang program and has been tested on Mac Silicon Procesors, but should work on other architectures and Operating systems as well. If you need to compile install golang preferably `go 1.21.6` version or higher.

## Contributing

PRs, testing and opening issues are all welcome! We are currently busy with the following:
1. Add more & better test cases.
2. Improve visualization, the sumary view should be better.
3. Improve the visualizations for the diffs in the -vv mode.
4. Issues and bug fixes.


# kompare

Go CLI runner to compare two clusters. This software compares two kubernetes cluster using kubeconfig to connect to them and compare existing objects in the two clusters.
## Why do we need kompare
This CLI tool has been created in the context of having to compare two clusters to determine if they are different so they can be interchangeable or replace each other. Enterprices often prefer to keep a few k8s clusters for the same job or run upgrade this way. The practical/real work we use the tool for is to compare a source cluster that is currently in production with a new cluster that we intend to put in prod. Therefore the source cluster "tends" to be considered the source of truth for the comparison.

## Key terms & help to use the tool.
1. Source cluster refers to the origin cluster. If this was a number comparison then the source cluster would be the "left-hand side" (LHS) or "antecedent".
2. Target cluster refers to the destination cluster. If this was a number comparison then the source cluster would be the "right-hand side" (RHS) or the "second number" or  "consequent".
3. The too has a help that should be self explanatory. you can get to this help by executing the tool with the -h flag.
```
go run main.go -h
usage: print [-h|--help] [-c|--conf "<value>"] [-s|--src "<value>"] -t|--target
             "<value>" [-v|--verbose] [-i|--include "<value>"] [-e|--exclude
             "<value>"] [-n|--namespace "<value>"]

             Prints provided string to stdout

Arguments:

  -h  --help       Print help information
  -c  --conf       Path to the clusters kubeconfig; assume ~/.kube/config if
                   not provided
  -s  --src        The Source cluster's context. Origin cluster in the
                   comparison (LHS-left hand side)
  -t  --target     *The target cluster's context (Required). Cluster used as
                   destination or consequent (RHS - Right hand side)
  -v  --verbose    Just show me all the diffs too. Notice: the output might be
                   LONG!
  -i  --include    List of kubernetes objects names to include, this should be
                   an element or a comma separated list.
  -e  --exclude    List of kubernetes objects to include, this should be an
                   element or a comma separated list.
  -n  --namespace  Namespace that needs to be copied. defaults to 'default'
                   namespace
```
So far only the -t option is required.

## Ownership
Team:   `xendit/infrastructure` 
Slack Handle: `@troops-infrastructure`

## Getting Started

This is a project currently in dev, the best way to use it it to do a go build or to just run it in the console with something like:
```
go run main.go -t some target context.
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
**Notice:** The software will assume you want to use your current context as source cluster to compare from (-s or --source to set a different context). Then the -t option is required and its the destination cluster. 

### Prerequisites

It's a golang program and has been tested on Mac Silicon Procesors, but should work on other architectures and Operating systems as well. If you need to compile install golang preferably `go 1.21.6` version or higher.

End with an example of getting some data out of the system or using it for a little demo.

## Contributing

We still do not have formal contribution procedures, but if you open a PR or issue for this project it will be reviewed and responded too.

# Roadmap

## Features (goot to have)
1. save comparison to file.
2. Compare file to target again.
3. (Done) Filter by "Kubernetes resource specifications type" e.g. Specs, Name, Annotations, etc. 
4. Customize filters for what to compare within a specific type of object.
5. Code Tests (unit testing).
6. run the cluster comparison as a test suite.
7. Sumary output version. All good vs There are some issues and a few lines of information about it.
8. Docs.
9. Make it do the works in parallel (go routines).
10. parser for different types of diff objects.
11. (Done) Bug source on left, target on the right for diffs.

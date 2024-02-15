# Kompare

Kompare is a Go CLI runner for comparing two clusters. This software compares two Kubernetes clusters using kubeconfig to connect to them and compare existing objects in the two clusters based on flexible criteria passed via the command line.

When comparing two Kubernetes clusters, you may encounter challenges such as ignoring Kubernetes resource definition subtypes like UID and dates, or filtering and comparing only specific types of objects. Existing tools may not offer the flexibility to define such criteria, which led us to create Kompare.

## Why Do We Need Kompare

This CLI tool is designed to compare two clusters to determine if they are different or if they can be interchangeable. Enterprises often maintain multiple Kubernetes clusters for redundancy or to facilitate upgrades. We primarily use this tool to compare a production cluster (source cluster) with a new cluster intended for production use or to run side by side.

**Notice:** The source cluster is typically considered the source of truth for the comparison.

## Key Terms & Usage

1. **Source cluster:** The original cluster, analogous to the "left-hand side" (LHS) in a number comparison.
2. **Target cluster:** The destination cluster, analogous to the "right-hand side" (RHS) or the "second number" in a number comparison.
3. The tool includes a help option that can be accessed using the `-h` flag.

```sh
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

This project is currently in development. To use it, either build it with Go or run it directly from the console:

```
go run main.go -t some-target-context.
```
see [this introductory post](https://medium/post/once/public) for mode details. 
### Example Command

To compare the current context with "MySecondContext-Cluster" and view the differences using the `-v` (verbose) option:

```
go run main.go -t MySecondContext-Cluster -v -n velero -e svc,role,rolebinding
We will use current kubeconfig context as 'source cluster'.
We will use MySecondContext-Cluster kubeconfig context as 'target cluster'.
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
**Notice:** The software assumes the current context as the source cluster by default (use `-s` or `--source` to set a different source context). The `-t` option specifies the destination cluster.

**Notice:** Therefore, the source cluster is typically considered the source of truth for the comparison.

### Prerequisites

Kompare is written in Go and has been tested on Mac Silicon Processors. It should work on other architectures and operating systems as well. Ensure you have Go 1.21.6 or higher installed.

## Contributing

Contributions such as pull requests, testing, and issue reporting are welcome. Current focus areas include:

1. Adding more and better test cases.
2. Improving visualization, particularly the summary view.
3. Enhancing visualizations for differences in verbose mode.
4. Addressing issues and bugs.


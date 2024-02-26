# Kompare

Kompare is a Go CLI runner for comparing two clusters. This software compares two Kubernetes clusters using kubeconfig to connect to them and compare existing objects in the two clusters based on flexible criteria passed via the command line.

When comparing two Kubernetes clusters, you may encounter challenges such as ignoring Kubernetes resource definition subtypes like UID and dates, or filtering and comparing only specific types of objects. Existing tools may not offer the flexibility to define such criteria, which led us to create Kompare.

![Kompare: Simplifying Kubernetes Cluster Comparison](https://miro.medium.com/v2/resize:fit:1400/1*oOPoArcHhU26oM0iUuGjAA.png)

## Why Do We Need Kompare

This CLI tool is designed to compare two clusters to determine if they are different or if they can be interchangeable. Enterprises often maintain multiple Kubernetes clusters for redundancy or to facilitate upgrades. We primarily use this tool to compare a production cluster (source cluster) with a new cluster intended for production use or to run side by side.

**Notice:** The source cluster is typically considered the source of truth for the comparison.

## Key Terms & Usage

1. **Source cluster:** The original cluster, analogous to the "left-hand side" (LHS) in a number comparison.
2. **Target cluster:** The destination cluster, analogous to the "right-hand side" (RHS) or the "second number" in a number comparison.

## Getting Started
You can run it by building out directly with go if you have golang installed:
```
$ git clone https://github.com/xendit/kompare.git
$ cd kompare
$ go build -o kompare .
$ ./kompare -h
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
$
```

To run directly with golang if you have it installed:
```
$ git clone https://github.com/xendit/kompare.git
$ cd kompare
$ go run main.go -h
(...)
```

### Using the software to compare two kubernetes clusters:
Notice: The only option requires is -t.

This project is currently in development. To use it, either build it with Go or run it directly from the console:

```
./kompare -t some-target-context.
```
The command with go would be `go run main.go -t some-target-context`.

see [this introductory post](https://blog.xendit.engineer/kompare-simplifying-kubernetes-cluster-comparison-ced2792716d9) for mode details. 
### Example Command

To compare the current context with "MySecondContext-Cluster", and view the differences using the `-v` (verbose see the help `-vv`) option, also in this case we are using `namespace kube-system` and we only want the `deployments`:

```
./kompare -t MySecondContext-Cluster  -v -n kube-system -i deploy
We will use current kubeconfig context as 'source cluster'.
We will use MySecondContext-Cluster kubeconfig context as 'target cluster'.
Using kube-system namespace
Looping namespace: kube-system
Deployment
******************************************************************************************************
- First cluster has Deployment in the list: blackbox-controller, but it's not in the second cluster

******************************************************************************************************
- Second cluster has Deployment in the list: aws-load-balancer-controller, but it's not in the first cluster

******************************************************************************************************
Finished Deployment for namespace: kube-system
Finished all comparison works!
```

Another example, we use velero for backups and want to see if both clusters have the same version deployed:
```
./kompare -t MySecondContext-Cluster -vv -n velero -i deploy
We will use current kubeconfig context as 'source cluster'.
We will use mycluster kubeconfig context as 'target cluster'.
Using velero namespace
Looping namespace: velero
Deployment
******************************************************************************************************
Done compering source cluster versus target cluster's  Deployment in the list
Done compering target cluster versus source cluster's  Deployment in the list
No differences found; Object Name velero, Kubernetes resource definition type Name, Namespace velero
Kubernetes resource definition type: Spec.Template.Spec
Object Name: velero
Namespace: velero
Differences:
- Containers.slice[0].Image: velero/velero:v1.12.2 != velero/velero:v1.9.2
- Containers.slice[0].Args.slice[1]: --uploader-type=restic != <no value>
- Containers.slice[0].LivenessProbe: v1.Probe != <nil pointer>
- Containers.slice[0].ReadinessProbe: v1.Probe != <nil pointer>
- TerminationGracePeriodSeconds: 3600 != 30


Finished Deployment for namespace: velero
Finished all comparison works!
```

**Notice:** The software assumes the current context as the source cluster by default (use `-s` or `--source` to set a different source context). The `-t` option specifies the destination/target cluster in your comparison. If it was number comparison -s is LHS and -t is RHS.

**Notice:** The source cluster is typically considered the source of truth for the comparison in Kompare.
## Prerequisites
### for use
For now you need to build the binary.
### Prerequisites for development

Kompare is written in Go and has been tested on Mac Silicon Processors. It should work on other architectures and operating systems as well. Ensure you have Go 1.21.6 or higher installed.

## Contributing

Contributions such as pull requests, testing, and issue reporting are welcome. Current focus areas include:

1. Adding more and better test cases.
2. Improving visualization, particularly the summary view.
3. Enhancing visualizations for differences in verbose mode.
4. Addressing issues and bugs.


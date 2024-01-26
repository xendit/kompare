# kompare

Go library and runner for running tests against Fawkes clusters.

## Ownership
Team:   `xendit/infrastructure` 

Slack Channel:    

Slack Handle: `@troops-infrastructure`

## Getting Started

Prepare a set of instructions to get a copy of the project up and running on a local machine, for both development and testing. 

There will be a separate deployment section for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them. This includes the explanation of configuration values that the service might need in order to run.

```
Give examples
```

### Installing

A step by step series of examples that tell you how to get a development env running.

Say what the step will be

```
Give the example
```

And repeat.

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo.

### Testing

Explain how to run the automated tests for this system.

```
Give an example
```

## Deployment

Add guides about how to deploy this on a live system.

## Contributing

Add guides about the process of contributing to the service.

## Additional Resource

You can link reading materials related to the service on this section. This may include but not limited to the API documentation, Confluence document, articles, etc.

# TODO
// compare globally:
// - DRDs (criteria?) <- Alpha done.
// - Same Namespaces exist in both clusters. <- Alpha done.
// - roles (criteria?)
// - clusterroles (criteria?)
// - rolebindings (criteria?)
// - clusterrolebindingss (criteria?)
// Compare per namespace
// - Deployment (Spec.Template.Spec & ?) <- Alpha done.
// - Services (Spec, Metadata.Annotations, Metadata.Labels ) <- Alpha done.
// - Service accounts (Metadata.Annotations, Metadata.Labels) <- Alpha done.
// - Secrets (Type, Data?) <- Alpha done.
// - Ingress (Needed?)
// Features (goot to have)
// - Verbose compare vs simple comparison
// - save comparison to file.
// - Compare file to target again.
// - Service specific particulars to compare; e.g.: when a type of object can have multiple ways
// of defining structure and we need to check some of those that are not always present.
// Documentation

# Resource booking operator
> Control and schedule cloud resources using custom Kubernetes operator

Manage your compute instances, or explore this general example on how a booking system for cloud computing resources could work, using Kubernetes custom resource operator.

## How it works
The custom resource operator provides a friendly interface to manage cloud resources through bookings.

We start by grouping our cloud instances under a common resource tag name, and then creating a custom resource with that tag on our cluster. Once we have resources, we can manage their state through bookings that have a resource name, start, and end time.

Example manifests can be found in the [config/samples] directory.

## Quick start

To play with the operator against a default local cluster, we first need to install the custom resource definitions:

```
make install
```

Start the operator:
```
make run
```
## The details

TODO: Emphasize and link to docs.

## Development
Kubebuilder is hard dependency of the project, so one of the best guides to extending and playing with this codebase is the [Kubebuilder book](https://book.kubebuilder.io/).

## Roadmap

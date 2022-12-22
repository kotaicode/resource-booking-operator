# Resource booking operator
> Control and schedule cloud resources using custom Kubernetes operator

Manage your compute instances or explore this as a general example on how a booking system for cloud computing resources could work, using Kubernetes custom resources.

## How it works
The custom resource operator provides a friendly interface to manage cloud resources through bookings.

We start by grouping our cloud instances under a common tag name, and then creating a custom resource with that tag on our cluster. Once we have resources, we can manage their state through bookings that have a resource name, start, and end time.

Example manifests can be found in the [config/samples](config/samples) directory.

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

📘 For more details on how to use the operator, we highly recommend [checking out the documentation](https://github.com/kotaicode/resource-booking-operator/edit/readme-to-docs/TODO.md).

## Development
Kubebuilder is a hard development dependency of the project, so one of the best guides to extending and playing with this codebase is the [Kubebuilder book](https://book.kubebuilder.io/).

## Roadmap

# Resource booking operator
> Control and schedule cloud resources using custom Kubernetes operator

Manage your compute instances, or explore this general example on how a booking system for cloud computing resources could work, using Kubernetes custom resource operator.

## How it works
The custom resource operator provides friendly interface to manage cloud resources through bookings.  

You start by grouping your cloud instances under a common resource tagname, and then creating a custom resource with that tag on your cluster. Once you have resources, you can manage their state through bookings that have a tag name, start, and end time.

Example yamls can be seen in the `config/samples` directory or in the [usage section](#usage).

## Quick start

To play with the operator against a default local cluster, we first need to install the custom resource definitions:

```
make install
```

## Usage

Start the operator.
```
make run
```
TODO: Emphasize and link to docs.

## Development
Kubebuilder is hard dependency of the project, so one of the best guides to extending and playing with this codebase is the [Kubebuilder book](https://book.kubebuilder.io/).

## Roadmap

# Resource booking operator
> Additional information or tagline

A brief description of your project, what it is used for and how does life get
awesome when someone starts to use it.

### Prerequisites
- Docker
- Minikube/Kubernetes cluster
- Kubebuilder

## Getting started

To play with the operator against a default local cluster, we first need to install the custom resource definitions:

```
make install
```

Next, refer to [Usage](#usage) and [Development](#development).

## How it works
The custom resource operator provides friendly interface to manage cloud resources through bookings.  

You start by grouping your cloud instances under a common resource tagname, and then creating a custom resource with that tag on your cluster. Once you have resources, you can manage their state through bookings that have a tag name, start, and end time.

Example yamls can be seen in the `config/samples` directory or in the [usage section](#usage).

## Usage

TODO: AWS client config
TODO: How to tag instances to make them visible as a resource

Start the operator.
```
make run
```

There are two custom resources - resources and bookings. Their template yamls which can be applied directly to the cluster are in the `config/samples` directory:

<details><summary>config/samples/manager_v1_resource.yaml</summary>
<p>

```yaml
apiVersion: manager.kotaico.de/v1
kind: Resource
metadata:
  labels:
    app.kubernetes.io/name: resource
    app.kubernetes.io/instance: web
    app.kubernetes.io/part-of: resource-booking-operator
    app.kuberentes.io/managed-by: kustomize
    app.kubernetes.io/created-by: resource-booking-operator
  name: web
spec:
  tag: web
  booked: false
```

</p>
</details>


<details><summary>config/samples/manager_v1_booking.yaml</summary>
<p>

```yaml
apiVersion: manager.kotaico.de/v1
kind: Booking
metadata:
  labels:
    app.kubernetes.io/name: booking
    app.kubernetes.io/instance: booking-sample
    app.kubernetes.io/part-of: resource-booking-operator
    app.kuberentes.io/managed-by: kustomize
    app.kubernetes.io/created-by: resource-booking-operator
  name: booking-sample
spec:
  resource_name: web
  start_at: 2006-01-02T15:04:05Z
  end_at: 2006-01-02T15:04:05Z
```

</p>
</details>

To start managing instances, you must change the template yamls according to the instance tag you want to manage and the time slot for which they will be booked, and finally - apply both of them to the cluster:

```
kubectl apply -f config/samples/manager_v1_resource.yaml
kubectl apply -f config/samples/manager_v1_booking.yaml
```

## Development
Kubebuilder is hard dependency of the project, so one of the best guides to extending and playing with this codebase is the [Kubebuilder book](https://book.kubebuilder.io/).

## Roadmap

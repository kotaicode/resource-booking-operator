# Resource booking operator
> Additional information or tagline

A brief description of your project, what it is used for and how does life get
awesome when someone starts to use it.

### Prerequisites
- Docker
- Minikube

## Getting started

TODO: Setting minikube or target cluster

```
make install
```

## How it works

TODO: Provides friendly interface to manage cloud resources through bookings.

## Usage

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

TODO: Auto discovery feature. Make command that runs the app with flag that discovers resources and creates them in the cluster.

## Development

TODO: Kubebuilder docs, make and make manifests commands.

## Roadmap

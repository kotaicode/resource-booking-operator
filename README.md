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

We start by creating the resources we want to manage. A hard prerequisite to that is to set up your cloud service credentials and tag the instances accordingly. More details can be found in the [extended documentation](https://kotaico.de/resource-booking-operator-docs/integrations/ec2/tagging-instances.html).

Since this is a quick start, we can ignore the manual creation of the cloud resource manifests and just use a custom resource we made for that purpose.
### Create a resource monitor
Their name hints at their purpose. Resource monitors sync on an interval with your cloud service and create new resources on the cluster for you. This makes the tagging process a bit more streamlined — You can tag instances in a way that they are visible to the operator and expect to see the corresponding new resources being created on the cluster in 2 minutes.

The only field that monitors require as of now is a type:
```yaml
apiVersion: manager.kotaico.de/v1
kind: ResourceMonitor
metadata:
  labels:
    app.kubernetes.io/name: resourcemonitor
    app.kubernetes.io/instance: ec2
    app.kubernetes.io/part-of: resource-booking-operator
    app.kuberentes.io/managed-by: kustomize
    app.kubernetes.io/created-by: resource-booking-operator
  name: ec2
spec:
  type: ec2
```

Which we then create in the cluster with:
```yaml
kubectl apply -f manager_v1_resourcemonitor.yaml
```

Once we see resources on the cluster through `kubectl get resource`, we can move on to managing their state with bookings.

### Create a booking
Bookings control resources. They dictate when a resource should be started or stopped. Say, on a given night, we want to start 2 instances that are grouped by a common tag — analytics. We prepare a manifest for the analytics resource that starts at 10 in the afternoon and stops at 10 minutes to midnight:

```yaml
apiVersion: manager.kotaico.de/v1
kind: Booking
metadata:
  labels:
    app.kubernetes.io/name: booking
    app.kubernetes.io/instance: analytics-jan01
    app.kubernetes.io/part-of: resource-booking-operator
    app.kuberentes.io/managed-by: kustomize
    app.kubernetes.io/created-by: resource-booking-operator
  name: analytics-jan01
spec:
  resource_name: analytics
  start_at: 2023-01-01T20:00:00Z
  end_at: 2023-01-01T23:50:00Z
  user_id: cd39ad8bc3

```

```yaml
kubectl apply -f manager_v1_booking.yaml
```

### Watch for changes
Once you create a booking, you can track their effect with:
```
kubectl get resources,bookings
```
Once the local cluster time hits the start time of the booking, you'll see the instances from this resource spinning up, and the booking status moving to being in progress. When the end time comes — the resource instances will be shut down and the booking will be marked as finished.
```
NAME                                   LOCKED BY   LOCKED UNTIL   INSTANCES   RUNNING   STATUS
analytics                                                         2           0         STOPPED

NAME                                      START                  END                    STATUS
analytics-jan01                           2023-01-01T20:00:00Z   2023-01-01T23:50:00Z   FINISHED
```


## The details

📘 For more details on how to use the operator, we highly recommend [checking out the documentation](https://kotaico.de/resource-booking-operator/).

## Tests
Tests are still in progress and work only in an environment with configured AWS credentials.  

First, make sure to install `envtest` with:
```
make envtest
```

Then running the tests is as simple as:
```
make test
```

## Development
Kubebuilder is a hard development dependency of the project, so one of the best guides to extending and playing with this codebase is the [Kubebuilder book](https://book.kubebuilder.io/).

## Roadmap
- Recurring bookings
- Support AWS RDS

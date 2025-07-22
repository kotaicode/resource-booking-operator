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

---

# Helm Chart

A Helm chart for deploying the Resource Booking Operator to Kubernetes clusters.

## Overview

The Resource Booking Operator is a Kubernetes operator that manages resource bookings and reservations. It provides custom resources for managing bookings, resources, and resource monitors.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

## Installation

### Add the Helm repository

```bash
helm repo add resource-booking-operator https://kotaicode.github.io/resource-booking-operator
helm repo update
```

### Install the chart

```bash
# Install with default values
helm install resource-booking-operator resource-booking-operator/resource-booking-operator

# Install with custom values
helm install resource-booking-operator resource-booking-operator/resource-booking-operator \
  --values custom-values.yaml
```

### Install from local chart

```bash
# Clone the repository
git clone https://github.com/kotaicode/resource-booking-operator.git
cd resource-booking-operator

# Install from local chart
helm install resource-booking-operator charts/resource-booking-operator/
```

## Configuration

The following table lists the configurable parameters of the resource-booking-operator chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `global.imageRegistry` | Global Docker image registry | `""` |
| `global.imagePullPolicy` | Global image pull policy | `IfNotPresent` |
| `global.imagePullSecrets` | Global image pull secrets | `[]` |
| `operator.image.repository` | Operator image repository | `controller` |
| `operator.image.tag` | Operator image tag | `latest` |
| `operator.image.pullPolicy` | Operator image pull policy | `IfNotPresent` |
| `operator.deployment.replicas` | Number of operator replicas | `1` |
| `operator.deployment.resources.limits.cpu` | CPU resource limits | `500m` |
| `operator.deployment.resources.limits.memory` | Memory resource limits | `128Mi` |
| `operator.deployment.resources.requests.cpu` | CPU resource requests | `10m` |
| `operator.deployment.resources.requests.memory` | Memory resource requests | `64Mi` |
| `rbac.create` | Create RBAC resources | `true` |
| `rbac.createClusterRoleBinding` | Create cluster role binding | `true` |
| `rbac.serviceAccount.create` | Create service account | `true` |
| `rbac.serviceAccount.name` | Service account name | `controller-manager` |
| `crd.install` | Install CRDs | `true` |
| `crd.validation` | Enable CRD validation | `true` |
| `metrics.enabled` | Enable metrics | `true` |
| `metrics.serviceMonitor.enabled` | Create service monitor | `false` |
| `authProxy.enabled` | Enable auth proxy | `true` |
| `authProxy.image.repository` | Auth proxy image repository | `gcr.io/kubebuilder/kube-rbac-proxy` |
| `authProxy.image.tag` | Auth proxy image tag | `v0.13.0` |
| `namespace.create` | Create namespace | `true` |
| `namespace.name` | Namespace name | `system` |

### Example custom values

```yaml
# custom-values.yaml
operator:
  image:
    repository: my-registry/controller
    tag: v1.0.0
  deployment:
    replicas: 2
    resources:
      limits:
        cpu: 1000m
        memory: 256Mi
      requests:
        cpu: 100m
        memory: 128Mi

rbac:
  create: true
  serviceAccount:
    create: true
    name: resource-booking-operator

metrics:
  enabled: true
  serviceMonitor:
    enabled: true

authProxy:
  enabled: true
  resources:
    limits:
      cpu: 200m
      memory: 64Mi
    requests:
      cpu: 10m
      memory: 32Mi
```

## Custom Resources

The operator creates the following custom resources:

### Resource

```yaml
apiVersion: manager.kotaico.de/v1
kind: Resource
metadata:
  name: example-resource
spec:
  type: "ec2"
  tag: "production"
  booked_by: "user@example.com"
  booked_until: "2024-12-31T23:59:59Z"
```

### Booking

```yaml
apiVersion: manager.kotaico.de/v1
kind: Booking
metadata:
  name: example-booking
spec:
  resource_name: "example-resource"
  user_id: "user@example.com"
  start_at: "2024-01-01T09:00:00Z"
  end_at: "2024-01-01T17:00:00Z"
  notifications:
    - type: "email"
      recipient: "user@example.com"
```

### ResourceMonitor

```yaml
apiVersion: manager.kotaico.de/v1
kind: ResourceMonitor
metadata:
  name: example-monitor
spec:
  type: "ec2"
```

## Additional RBAC Roles

The chart includes additional RBAC roles for fine-grained access control:

- `booking-editor`: Full access to bookings
- `booking-viewer`: Read-only access to bookings
- `resource-editor`: Full access to resources
- `resource-viewer`: Read-only access to resources
- `resourcemonitor-editor`: Full access to resource monitors
- `resourcemonitor-viewer`: Read-only access to resource monitors

These roles can be enabled/disabled via the `additionalRoles` section in values.yaml.

## Monitoring

The operator exposes metrics on port 8080 (or 8443 when auth proxy is enabled). You can configure Prometheus ServiceMonitor for monitoring:

```yaml
metrics:
  enabled: true
  serviceMonitor:
    enabled: true
    interval: 30s
    scrapeTimeout: 10s
```

## Security

The operator runs with the following security features:

- Non-root user execution
- Dropped capabilities
- No privilege escalation
- RBAC authorization via auth proxy (when enabled)

## Troubleshooting

### Check operator status

```bash
kubectl get pods -n system -l control-plane=controller-manager
```

### View operator logs

```bash
kubectl logs -n system -l control-plane=controller-manager -c manager
```

### Check CRD installation

```bash
kubectl get crd | grep manager.kotaico.de
```

### Test operator functionality

```bash
helm test resource-booking-operator
```

## Uninstallation

```bash
helm uninstall resource-booking-operator
```

**Note**: CRDs are not automatically removed. To remove them:

```bash
kubectl delete crd resources.manager.kotaico.de
kubectl delete crd bookings.manager.kotaico.de
kubectl delete crd resourcemonitors.manager.kotaico.de
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

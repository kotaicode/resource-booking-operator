# Resource Booking Operator Helm Chart

This Helm chart deploys the Resource Booking Operator, a Kubernetes operator for managing AWS resource bookings and scheduling.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- AWS credentials configured (via IAM roles, environment variables, or mounted secrets)

## Installation

### Add the Helm repository (if published)
```bash
helm repo add resource-booking-operator https://charts.kotaico.de
helm repo update
```

### Install the chart
```bash
# Install with default values
helm install resource-booking-operator ./helm

# Install with custom values
helm install resource-booking-operator ./helm -f custom-values.yaml

# Install in a specific namespace
helm install resource-booking-operator ./helm --namespace my-namespace --create-namespace
```

## Configuration

The following table lists the configurable parameters of the resource-booking-operator chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `operator.image.repository` | Container image repository | `kotaicode/resource-booking-operator` |
| `operator.image.tag` | Container image tag | `latest` |
| `operator.image.pullPolicy` | Container image pull policy | `IfNotPresent` |
| `operator.replicas` | Number of operator replicas | `1` |
| `operator.namespace` | Namespace to deploy the operator | `system` |
| `operator.resources.limits.cpu` | CPU resource limit | `500m` |
| `operator.resources.limits.memory` | Memory resource limit | `128Mi` |
| `operator.resources.requests.cpu` | CPU resource request | `10m` |
| `operator.resources.requests.memory` | Memory resource request | `64Mi` |
| `operator.metrics.enabled` | Enable metrics endpoint | `true` |
| `operator.metrics.port` | Metrics port | `8080` |
| `operator.healthProbe.enabled` | Enable health probes | `true` |
| `operator.healthProbe.port` | Health probe port | `8081` |
| `operator.leaderElection.enabled` | Enable leader election | `false` |
| `operator.webhook.enabled` | Enable webhook server | `false` |
| `operator.webhook.port` | Webhook port | `9443` |
| `crd.install` | Install Custom Resource Definitions | `true` |
| `serviceMonitor.enabled` | Enable Prometheus ServiceMonitor | `false` |
| `hpa.enabled` | Enable Horizontal Pod Autoscaler | `false` |
| `podDisruptionBudget.enabled` | Enable Pod Disruption Budget | `false` |
| `ingress.enabled` | Enable Ingress | `false` |
| `networkPolicy.enabled` | Enable Network Policy | `false` |

## Custom Resource Definitions

The chart installs the following CRDs:

- **Resource**: Represents AWS resources (EC2, RDS, etc.)
- **ResourceMonitor**: Monitors AWS resources on a schedule
- **Booking**: Represents a resource booking with start/end times
- **BookingScheduler**: Automatically creates bookings on a schedule

## Usage Examples

### Create an EC2 Resource
```yaml
apiVersion: manager.kotaico.de/v1
kind: Resource
metadata:
  name: my-ec2-instance
spec:
  type: ec2
  name: production-server
  region: eu-central-1
  tags:
    Environment: production
    Team: devops
```

### Create a Resource Monitor
```yaml
apiVersion: manager.kotaico.de/v1
kind: ResourceMonitor
metadata:
  name: ec2-monitor
spec:
  type: ec2
  schedule: "*/5 * * * *"  # Every 5 minutes
  filters:
    Environment: production
```

### Create a Booking
```yaml
apiVersion: manager.kotaico.de/v1
kind: Booking
metadata:
  name: maintenance-window
spec:
  resourceName: my-ec2-instance
  startTime: "2024-01-15T02:00:00Z"
  endTime: "2024-01-15T04:00:00Z"
  user: "maintenance-team"
```

### Create a Booking Scheduler
```yaml
apiVersion: manager.kotaico.de/v1
kind: BookingScheduler
metadata:
  name: daily-backup
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  template:
    resourceName: my-ec2-instance
    startTime: "2024-01-15T02:00:00Z"
    endTime: "2024-01-15T03:00:00Z"
    user: "backup-system"
```

## Monitoring

### Metrics
If metrics are enabled, you can access them via:
```bash
kubectl port-forward -n system svc/resource-booking-operator-metrics 8080:8080
curl http://localhost:8080/metrics
```

### Health Checks
Health checks are available at:
```bash
kubectl port-forward -n system svc/resource-booking-operator-health 8081:8081
curl http://localhost:8081/healthz
curl http://localhost:8081/readyz
```

## Upgrading

```bash
helm upgrade resource-booking-operator ./helm
```

## Uninstalling

```bash
helm uninstall resource-booking-operator
```

**Note**: CRDs are not automatically removed. To remove them:
```bash
kubectl delete crd resources.manager.kotaico.de
kubectl delete crd resourcemonitors.manager.kotaico.de
kubectl delete crd bookings.manager.kotaico.de
kubectl delete crd bookingschedulers.manager.kotaico.de
```

## Troubleshooting

### Check operator status
```bash
kubectl get pods -n system -l app.kubernetes.io/name=resource-booking-operator
kubectl logs -n system -l app.kubernetes.io/name=resource-booking-operator
```

### Check CRDs
```bash
kubectl get crd | grep manager.kotaico.de
```

### Check operator logs
```bash
kubectl logs -n system deployment/resource-booking-operator
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test the chart
5. Submit a pull request

## License

This chart is licensed under the Apache License 2.0. 
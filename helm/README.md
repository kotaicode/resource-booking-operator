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

## Examples

### BookingScheduler with Time-based Scheduling

```yaml
apiVersion: manager.kotaico.de/v1
kind: BookingScheduler
metadata:
  name: daily-backup
spec:
  end_time: "13:00"
  resource_name: resource-booking-operator
  schedule: "00 06 * * 1-5"
  start_time: "10:00"
  user_id: tom
```

**Smart Scheduling Behavior:**
- **Before start time**: Resource starts at the scheduled start time
- **Within time window**: Resource starts immediately (e.g., if it's 11:30 AM and window is 10:00 AM - 1:00 PM)
- **After end time**: Resource starts tomorrow at the scheduled start time
- **Automatic stop**: Resource automatically stops when end time is reached

### BookingScheduler with Timestamp-based Scheduling

```yaml
apiVersion: manager.kotaico.de/v1
kind: BookingScheduler
metadata:
  name: specific-time-booking
spec:
  end_at: "2025-08-18T13:00:00+02:00"
  resource_name: resource-booking-operator
  schedule: "00 06 * * 1-5"
  start_at: "2025-08-18T10:00:00+02:00"
  user_id: tom
```

**Timestamp-based Behavior:**
- Uses exact timestamps for start and end times
- Perfect for one-time or specific date scheduling
- Resources start and stop at the exact specified times

## Usage Examples

### Create an EC2 Resource
```
```
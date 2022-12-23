// Package clients contains common logic and data structures for the supported cloud providers.
package clients

import "github.com/kotaicode/resource-booking-operator/clients/ec2"

const (
    TypeEC2 string = "ec2"
)

const (
    StatusStopped string = "STOPPED"
    StatusRunning string = "RUNNING"
    StatusPending string = "PENDING"
)

// ResourceStatus holds the status summary of the resource.
type ResourceStatus struct {
    Available, Running int
}

// CloudResource provides generic Resource interface. A Resource is a group of instances which
// can be started or stopped. The interface also requires a method to list their statuses.
type CloudResource interface {
    Start() error
    Stop() error
    Status() (ResourceStatus, error)
}

// ResourceFactory generates structs that abide by the CloudResource interface.
// The returned struct can start, stop, and list instances. Each new integration needso to be added to this factory function.
func ResourceFactory(resType, tag string) CloudResource {
    var resource CloudResource

    switch resType {
    case TypeEC2:
        resource = &ec2.Resource{NameTag: tag}
    }

    return resource
}

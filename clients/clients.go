// Package clients contains common logic and data structures for the supported cloud providers.
package clients

import (
	"errors"
	"path/filepath"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
func ResourceFactory(resType, tag string) (CloudResource, error) {
	var resource CloudResource

	switch resType {
	case TypeEC2:
		resource = &EC2Resource{NameTag: tag}
	default:
		return nil, errors.New("Resource type not found")
	}

	return resource, nil
}

// TODO: Bring back the EC2 back in this package and make it a struct or something?

func GetK8sClient() (client.Client, error) {
	// TODO Inspire from the "using k8s client outside of cluster"
	// TODO Test the kubeconfig flag here?
	cfg, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}

	// TODO Try with and without that
	scheme := runtime.NewScheme()
	utilruntime.Must(managerv1.AddToScheme(scheme))

	clientOpts := client.Options{Scheme: scheme}

	c, err := client.New(cfg, clientOpts)
	if err != nil {
		return nil, err
	}

	return c, nil
}

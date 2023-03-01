// Package clients contains common logic and data structures for the supported cloud providers.
package clients

import (
	"errors"
	"flag"
	"os"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/cache"
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

// ResourceStatusOutput holds the status summary of the resource.
type ResourceStatusOutput struct {
	Available, Running    int
	LockedBy, LockedUntil string
}

// ResourceStartInput stores data that is used for book-keeping during the starting of the resource
type ResourceStartInput struct {
	UID, EndAt string
}

// ResourceStartInput stores data that is used for book-keeping during the stopping of the resource
type ResourceStopInput struct {
	UID string
}

// ClientCache holds the client and cache objects.
type ClientCache struct {
	Client client.Client
	Cache  cache.Cache
}

// CloudResource provides generic Resource interface. A Resource is a group of instances which
// can be started or stopped. The interface also requires a method to list their statuses.
type CloudResource interface {
	Start(startInput ResourceStartInput) error
	Stop(stopInput ResourceStopInput) error
	Status() (ResourceStatusOutput, error)
}

type ResourceMonitor interface {
	GetNewResources(clusterResources map[string]bool) ([]string, error)
}

var kubeconfig string

func init() {
	// Check for kube config. Priority is:
	// 1. --kubeconfig flag.
	// 2. KUBECONFIG env variable
	// 3. Default = /.kube/config
	kubeflag := flag.Lookup("kubeconfig")
	if kubeflag != nil {
		kubeconfig = flag.Lookup("kubeconfig").Value.String()
	}

	if kubeconfig == "" {
		kubeconfig, _ = os.LookupEnv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = "/.kube/config"
		}
	}
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

// MonitorFactory generates structs that abide by the ResourceMonitor interface.
// The returned struct can get new resources of the specified type. Each new integration needso to be added to this factory function.
func MonitorFactory(monitorType string) (ResourceMonitor, error) {
	var resourceMonitor ResourceMonitor

	switch monitorType {
	case TypeEC2:
		resourceMonitor = &EC2Monitor{Type: monitorType}
	default:
		return nil, errors.New("Monitor type not found")
	}

	return resourceMonitor, nil
}

// GetClient returns a ready to use kubernetes client and cache.
func GetClient() (ClientCache, error) {
	var clientCache ClientCache
	var err error
	var config *rest.Config

	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(managerv1.AddToScheme(scheme))
	clientOpts := client.Options{Scheme: scheme}

	clientCache.Client, err = client.New(config, clientOpts)
	if err != nil {
		return clientCache, err
	}

	clientCache.Cache, err = cache.New(config, cache.Options{Namespace: "default", Scheme: scheme})
	if err != nil {
		return clientCache, err
	}

	return clientCache, nil
}

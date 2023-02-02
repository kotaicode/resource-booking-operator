// Package ec2 is a thin wrapper around the EC2 Client.
// Along with that it provides a few additions, like grabbing the instance IDs that have a predefined tag which we use to identify resources.
package clients

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	// Integer representation of the instance status
	statusPending  int64 = 0
	statusRunning  int64 = 16
	statusStopping int64 = 64
	statusStopped  int64 = 80

	// DefaultTagKey is used to store the tag which marks the instance as managed by the operator
	defaultTagKey         string = "resource-booking/application"
	resourceMonitorTagKey string = "resource-booking/managed"
)

type EC2Monitor struct {
	Type string
}

// Resource represents a collection of EC2 instances grouped by a common "resource-booking/application" tag.
type EC2Resource struct {
	NameTag string
}

var mySession *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var ec2Client *ec2.EC2 = ec2.New(mySession)

// Start makes a call through the EC2 client to start resource instances by their IDs.
func (r *EC2Resource) Start() error {
	instanceIds, err := r.getInstanceIds(r.NameTag)
	if err != nil {
		return err
	}

	_, err = ec2Client.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return err
	}

	return nil
}

// Stop makes a call through the EC2 client to stop the instances that belong to the resource.
func (r *EC2Resource) Stop() error {
	instanceIds, err := r.getInstanceIds(r.NameTag)
	if err != nil {
		return err
	}

	_, err = ec2Client.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return err
	}

	return nil
}

// Status returns the current summary of a given resource instance statuses.
// It makes a call through the EC2 client with a given set of instance IDs and summarises their status (active vs running).
func (r *EC2Resource) Status() (ResourceStatus, error) {
	includeAll := true
	var rst ResourceStatus

	instanceIds, err := r.getInstanceIds(r.NameTag)
	if err != nil {
		return rst, err
	}

	resp, err := ec2Client.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: &includeAll,
		InstanceIds:         instanceIds,
	})
	if err != nil {
		return rst, err
	}

	// EC2 API will return all instances if we send an empty instance IDs list. Handle that.
	if len(instanceIds) > 0 {
		for _, inst := range resp.InstanceStatuses {
			rst.Available++
			if *inst.InstanceState.Code == statusRunning {
				rst.Running++
			}
		}
	}

	return rst, nil
}

// getInstanceIds returns instance IDs from a given name tag. Basically wrap the EC2 call with a filter of our default tag identificator.
func (r *EC2Resource) getInstanceIds(nameTag string) ([]*string, error) {
	var instanceIds []*string

	// Prepare filters
	tagKey := fmt.Sprintf("tag:%s", defaultTagKey)
	nameFilter := &ec2.Filter{
		Name:   &tagKey,
		Values: []*string{&nameTag},
	}

	resp, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{nameFilter},
	})
	if err != nil {
		return nil, err
	}

	for _, reserv := range resp.Reservations {
		for _, inst := range reserv.Instances {
			instanceIds = append(instanceIds, inst.InstanceId)
		}
	}

	return instanceIds, nil
}

// sliceToMap returns map, used to convert slice to map
func sliceToMap(slice []string) map[string]bool { //TODO: give a better name
	m := make(map[string]bool)
	for _, s := range slice {
		m[s] = true
	}
	return m
}

// diff returns a slice of strings, used to compare two maps.
func diff(m1, m2 map[string]bool) []string { //TODO: give a better name
	slice := make([]string, 0, len(m1))
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			slice = append(slice, k)
		}
	}
	return slice
}

// GetUniqueTags returns a slice of unique tags.
// It makes a call through the EC2 client to get all the unique tags in the cluster.
func GetUniqueTags() ([]string, error) {
	var uniqueTags []string

	// Prepare filters
	tagKey := "tag:" + resourceMonitorTagKey
	tagValue := "true"
	nameFilter := &ec2.Filter{
		Name:   &tagKey,
		Values: []*string{&tagValue},
	}
	tagMap := make(map[string]bool)
	resourceBookingInstances, err := ec2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{nameFilter},
	})

	if err != nil {
		return nil, err
	}

	for _, reservation := range resourceBookingInstances.Reservations {
		for _, instance := range reservation.Instances {
			resourceBookingTags := instance.Tags
			for _, v := range resourceBookingTags {
				if *v.Key == defaultTagKey {
					tagMap[*v.Value] = true
				}
			}
		}
	}

	for tag := range tagMap {
		uniqueTags = append(uniqueTags, tag)
	}
	return uniqueTags, nil

}

func (m *EC2Monitor) GetNewResources(clusterResources []string) ([]string, error) {
	uniqueTags, err := GetUniqueTags()
	if err != nil {
		return nil, err
	}

	uniqueTagsMap := sliceToMap(uniqueTags)
	clusterResourcesMap := sliceToMap(clusterResources)

	slice1 := diff(uniqueTagsMap, clusterResourcesMap) //TODO: give a better name
	slice2 := diff(clusterResourcesMap, uniqueTagsMap) //TODO: give a better name
	nonMatchingTags := append(slice1, slice2...)

	return nonMatchingTags, nil
}

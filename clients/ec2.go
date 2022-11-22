// Package ec2 is a thin wrapper around the EC2 Client. Along with that it provides a few additions, like grabbing the instance IDs that have a predefined tag that we use to identify resources.
package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// StatusPending holds the EC2 integer representation of the pending status
const StatusPending int64 = 0

// StatusRunning holds the EC2 integer representation of the running status
const StatusRunning int64 = 16

// StatusStopping holds the EC2 integer representation of the stopping status
const StatusStopping int64 = 64

// StatusStopped holds the EC2 integer representation of the stopping status
const StatusStopped int64 = 80

const defaultTagKey string = "resource-booking/application"

// Resource is a struct specifically created because of the methods assigned to it.
// Although a good TODO will be to check if it's a good practice to tie them to the repository struct which is almost identical
type Resource struct {
	NameTag    string
	IsArchived bool
}

// ResourceStatus holds the status summary of the resource
type ResourceStatus struct {
	InstanceStatusName string `json:"instance_status_name"`
	InstanceStatusCode int64  `json:"instance_status_code"`
}

var mySession *session.Session = session.Must(session.NewSession())

// EC2Client is the EC2 client session
var EC2Client *ec2.EC2 = ec2.New(mySession)

// Start makes a call through the EC2 client to start the instances from a given resource by their IDs
func (r *Resource) Start() error {
	instanceIds, _ := r.getInstanceIds(r.NameTag)
	_, err := EC2Client.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		// TODO: Log
		fmt.Println("TODO LOG", err.Error())
		return err
	}

	return nil
}

// Stop makes a call through the EC2 client to stop the instances that belong to the given resource
func (r *Resource) Stop() error {
	instanceIds, _ := r.getInstanceIds(r.NameTag)
	_, err := EC2Client.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		// TODO: Log
		return err
	}

	return nil
}

// Status returns the current summary of a given resource instance statuses. It makes a call through the EC2 client with a given set of instance IDs and summarises their status (active vs running)
func (r *Resource) Status() ([]ResourceStatus, error) {
	var rst []ResourceStatus

	if !r.IsArchived {
		includeAll := true
		instanceIds, _ := r.getInstanceIds(r.NameTag)
		resp, err := EC2Client.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
			IncludeAllInstances: &includeAll,
			InstanceIds:         instanceIds,
		})
		if err != nil {
			// TODO: Log
			fmt.Println("TODO LOG", err.Error())
			return rst, err
		}

		// EC2 API will return all instances if we send an empty instance IDs list. Handle that.
		if len(instanceIds) > 0 {
			for _, inst := range resp.InstanceStatuses {
				rst = append(rst, ResourceStatus{
					InstanceStatusName: *inst.InstanceState.Name,
					InstanceStatusCode: *inst.InstanceState.Code,
				})
			}
		}
	}

	return rst, nil
}

// Get instance IDs from a given name tag. Basically wrap the EC2 call with a filter of our default tag identificator.
func (r *Resource) getInstanceIds(nameTag string) ([]*string, error) {
	var instanceIds []*string

	tagKey := fmt.Sprintf("tag:%s", defaultTagKey)
	nameFilter := &ec2.Filter{
		Name:   &tagKey,
		Values: []*string{&nameTag},
	}

	resp, err := EC2Client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{nameFilter},
	})
	if err != nil {
		// TODO: Log
		return nil, err
	}

	for _, reserv := range resp.Reservations {
		for _, inst := range reserv.Instances {
			instanceIds = append(instanceIds, inst.InstanceId)
		}
	}

	return instanceIds, nil
}

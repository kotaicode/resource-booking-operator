// Package ec2 is a thin wrapper around the EC2 Client.
// Along with that it provides a few additions, like grabbing the instance IDs that have a predefined tag which we use to identify resources.
package clients

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
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

var (
	lockedByTag    string = "resource-booking/locked-by"
	lockedUntilTag string = "resource-booking/locked-until"
)

type EC2Monitor struct {
	Type string
}

// Resource represents a collection of EC2 instances grouped by a common "resource-booking/application" tag.
type EC2Resource struct {
	NameTag string
}

type instanceDetails struct {
	IDs  []*string
	Tags map[string]string
}

var mySession *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var ec2Client *ec2.EC2 = ec2.New(mySession)

func init() {
	if os.Getenv("AWS_ROLE_ARN") != "" {
		err := assumeRole(ec2Client, os.Getenv("AWS_ROLE_ARN"))
		if err != nil {
			panic(err)
		}
	}
}

// Start makes a call through the EC2 client to start resource instances by their IDs.
func (r *EC2Resource) Start(startInput ResourceStartInput) error {

	instances, err := r.getInstanceDetails(r.NameTag)
	if err != nil {
		return err
	}

	if _, err = r.canManage(startInput.UID, instances.Tags); err != nil {
		return err
	}

	_, err = ec2Client.StartInstances(&ec2.StartInstancesInput{
		InstanceIds: instances.IDs,
	})
	if err != nil {
		return err
	}

	err = r.lock(startInput.UID, startInput.EndAt, instances.IDs)
	if err != nil {
		return err
	}

	return nil
}

// Stop makes a call through the EC2 client to stop the instances that belong to the resource.
func (r *EC2Resource) Stop(stopInput ResourceStopInput) error {
	instances, err := r.getInstanceDetails(r.NameTag)
	if err != nil {
		return err
	}

	if _, err = r.canManage(stopInput.UID, instances.Tags); err != nil {
		return err
	}

	_, err = ec2Client.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: instances.IDs,
	})
	if err != nil {
		return err
	}

	err = r.unlock(instances.IDs)
	if err != nil {
		return err
	}

	return nil
}

// Status returns the current summary of a given resource instance statuses.
// It makes a call through the EC2 client with a given set of instance IDs and summarises their status (active vs running).
func (r *EC2Resource) Status() (ResourceStatusOutput, error) {
	includeAll := true
	var rst ResourceStatusOutput

	instances, err := r.getInstanceDetails(r.NameTag)
	if err != nil {
		return rst, err
	}

	resp, err := ec2Client.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: &includeAll,
		InstanceIds:         instances.IDs,
	})
	if err != nil {
		return rst, err
	}

	// EC2 API will return all instances if we send an empty instance IDs list. Handle that.
	if len(instances.IDs) > 0 {
		for _, inst := range resp.InstanceStatuses {
			rst.Available++
			if *inst.InstanceState.Code == statusRunning {
				rst.Running++
			}
		}
	}

	rst.LockedBy, rst.LockedUntil = instances.Tags[lockedByTag], instances.Tags[lockedUntilTag]

	return rst, nil
}

func (r *EC2Resource) canManage(uid string, instanceTags map[string]string) (bool, error) {
	if _, ok := instanceTags[lockedByTag]; ok && instanceTags[lockedUntilTag] != "" {
		d, err := time.Parse(time.RFC3339, instanceTags[lockedUntilTag])
		if err != nil {
			return false, err
		}

		if instanceTags[lockedByTag] != uid && time.Now().Before(d) {
			m := "Resource is locked by %s. The lock expires at %s."
			err = fmt.Errorf(m, instanceTags[lockedByTag], instanceTags[lockedUntilTag])
			return false, err
		}
	}

	return true, nil
}

// lock sets locking tags to the resource instances. Tags are:
// resource-booking/locked-by    - The identifier of the booking that owns the instance at this moment
// resource-booking/locked-until - Date time until the instance is available again. The endAt of the booking.
func (r *EC2Resource) lock(uid string, endAt string, instanceIDs []*string) error {
	_, err := ec2Client.CreateTags(&ec2.CreateTagsInput{
		Resources: instanceIDs,
		Tags: []*ec2.Tag{
			{Key: &lockedByTag, Value: &uid},
			{Key: &lockedUntilTag, Value: &endAt},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

// unlock removes the locking tags, freeing the resource to other users.
func (r *EC2Resource) unlock(instanceIDs []*string) error {
	_, err := ec2Client.DeleteTags(&ec2.DeleteTagsInput{
		Resources: instanceIDs,
		Tags: []*ec2.Tag{
			{Key: &lockedByTag},
			{Key: &lockedUntilTag},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

// getInstanceDetails returns instance IDs from a given name tag. Basically wrap the EC2 call with a filter of our default tag identificator.
func (r *EC2Resource) getInstanceDetails(nameTag string) (instanceDetails, error) {
	details := instanceDetails{Tags: make(map[string]string)}

	var instanceTagList []*ec2.Tag

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
		return details, err
	}

	for _, reserv := range resp.Reservations {
		for _, inst := range reserv.Instances {
			details.IDs = append(details.IDs, inst.InstanceId)
			instanceTagList = append(instanceTagList, inst.Tags...)
		}
	}

	// For now Don't care about the edge case where two instances might theoretically have different lock tags
	for _, v := range instanceTagList {
		if *v.Key == lockedByTag || *v.Key == lockedUntilTag {
			details.Tags[*v.Key] = *v.Value
		}

		if len(details.Tags) == 2 {
			break
		}
	}

	return details, nil
}

// GetNewResources compares the local cluster resources with the ones returned from EC2
// and gives back a list of resources that need to be created on the cluster.
func (m *EC2Monitor) GetNewResources(clusterResources map[string]bool) ([]string, error) {
	uniqueTags, err := GetUniqueTags()
	if err != nil {
		return nil, err
	}

	slice1, slice2 := setDiff(uniqueTags, clusterResources), setDiff(clusterResources, uniqueTags)
	nonMatchingTags := append(slice1, slice2...)

	return nonMatchingTags, nil
}

// GetUniqueTags makes a call through the EC2 client to collect all instance tags and returns a set of them
func GetUniqueTags() (map[string]bool, error) {
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

	return tagMap, nil
}

// setDiff returns the difference between two sets
func setDiff(m1, m2 map[string]bool) []string {
	slice := make([]string, 0, len(m1))
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			slice = append(slice, k)
		}
	}
	return slice
}

// assumeRole assumes a role and returns a new EC2 client with the new credentials
func assumeRole(cli *ec2.EC2, roleArn string) error {
	stsClient := sts.New(mySession)

	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("resource-booking-operator"),
	}

	resp, err := stsClient.AssumeRole(params)
	if err != nil {
		return err
	}

	creds := resp.Credentials
	cli.Config.Credentials = credentials.NewStaticCredentials(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)

	return nil
}

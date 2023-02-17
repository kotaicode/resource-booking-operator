package clients

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSResource struct {
	NameTag string
}

type RDSMonitor struct {
	Type string
}

type RDSInstanceDetails struct {
	IDs           []*string
	Tags          map[string]string
	ResourceNames []*string
}

var rdsSession *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var rdsClient = rds.New(rdsSession)

func (r *RDSResource) Start(startInput ResourceStartInput) error {
	instances, err := r.getRDSInstanceDetails(r.NameTag)
	if err != nil {
		return err
	}

	if _, err = r.canManageRDS(startInput.UID, instances.Tags); err != nil {
		return err
	}

	for _, DBInstance := range instances.IDs {
		_, err = rdsClient.StartDBInstance(&rds.StartDBInstanceInput{
			DBInstanceIdentifier: DBInstance,
		})
		if err != nil {
			return err
		}
	}

	err = r.lockRDS(startInput.UID, startInput.EndAt, instances.ResourceNames)
	if err != nil {
		return err
	}

	return nil
}

func (r *RDSResource) Stop(stopInput ResourceStopInput) error {
	instances, err := r.getRDSInstanceDetails(r.NameTag)
	if err != nil {
		return err
	}

	if _, err = r.canManageRDS(stopInput.UID, instances.Tags); err != nil {
		return err
	}
	for _, instance := range instances.IDs {
		_, err = rdsClient.StopDBInstance(&rds.StopDBInstanceInput{
			DBInstanceIdentifier: instance,
		})
		if err != nil {
			return err
		}
	}

	err = r.unlockRDS(instances.ResourceNames)
	if err != nil {
		return err
	}

	return nil
}

func (r *RDSResource) Status() (ResourceStatusOutput, error) {
	var rst ResourceStatusOutput
	details, err := r.getRDSInstanceDetails(r.NameTag)

	if err != nil {
		return rst, err
	}

	instances, err := r.getRDSInstancesByTag(r.NameTag)
	if err != nil {
		return rst, err
	}

	if err != nil {
		return rst, err
	}

	for _, inst := range instances {
		rst.Available++
		if *inst.DBInstanceStatus == "available" {
			rst.Running++
		}
	}
	rst.LockedBy, rst.LockedUntil = details.Tags[lockedByTag], details.Tags[lockedUntilTag]

	return rst, nil
}

func (r *RDSResource) getRDSInstancesByTag(nameTag string) ([]*rds.DBInstance, error) {

	// Retrieve the list of all DB instances
	instances, err := rdsClient.DescribeDBInstances(nil)
	if err != nil {
		return nil, err
	}

	// Filter the instances based on the specified tag key and value
	var filteredInstances []*rds.DBInstance
	for _, instance := range instances.DBInstances {
		input := &rds.ListTagsForResourceInput{
			ResourceName: instance.DBInstanceArn,
		}
		result, err := rdsClient.ListTagsForResource(input)
		if err != nil {
			return nil, err
		}
		for _, tag := range result.TagList {
			if *tag.Key == defaultTagKey && *tag.Value == nameTag {
				filteredInstances = append(filteredInstances, instance)
				break
			}
		}
	}
	return filteredInstances, nil
}

func (m *RDSMonitor) GetNewResources(clusterResources map[string]bool) ([]string, error) {
	uniqueTags, err := GetUniqueRDSTags()
	if err != nil {
		return nil, err
	}

	slice1, slice2 := setDiff(uniqueTags, clusterResources), setDiff(clusterResources, uniqueTags)
	nonMatchingTags := append(slice1, slice2...)

	return nonMatchingTags, nil
}

func GetUniqueRDSTags() (map[string]bool, error) {
	tagMap := make(map[string]bool)

	instances, err := rdsClient.DescribeDBInstances(nil)
	if err != nil {
		return nil, err
	}

	var filteredInstances []*rds.DBInstance
	for _, instance := range instances.DBInstances {
		input := &rds.ListTagsForResourceInput{
			ResourceName: instance.DBInstanceArn,
		}
		result, err := rdsClient.ListTagsForResource(input)
		if err != nil {
			return nil, err
		}
		for _, tag := range result.TagList {
			if *tag.Key == resourceMonitorTagKey && *tag.Value == "true" {
				filteredInstances = append(filteredInstances, instance)
				break
			}
		}
	}

	if err != nil {
		return nil, err
	}

	for _, instance := range filteredInstances {
		resourceBookingTags := instance.TagList
		for _, v := range resourceBookingTags {
			if *v.Key == defaultTagKey {
				tagMap[*v.Value] = true
			}
		}

	}
	return tagMap, nil
}

// lock sets locking tags to the resource instances. Tags are:
// resource-booking/locked-by    - The identifier of the booking that owns the instance at this moment
// resource-booking/locked-until - Date time until the instance is available again. The endAt of the booking.
func (r *RDSResource) lockRDS(uid string, endAt string, resourceNames []*string) error {
	for _, resourceName := range resourceNames {
		_, err := rdsClient.AddTagsToResource(&rds.AddTagsToResourceInput{
			ResourceName: resourceName,
			Tags: []*rds.Tag{
				{Key: &lockedByTag, Value: &uid},
				{Key: &lockedUntilTag, Value: &endAt},
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// unlock removes the locking tags, freeing the resource to other users.
func (r *RDSResource) unlockRDS(resourceNames []*string) error {
	for _, resourceName := range resourceNames {
		_, err := rdsClient.RemoveTagsFromResource(&rds.RemoveTagsFromResourceInput{
			ResourceName: resourceName,
			TagKeys:      []*string{&lockedByTag, &lockedUntilTag},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RDSResource) canManageRDS(uid string, instanceTags map[string]string) (bool, error) {
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

func (r *RDSResource) getRDSInstanceDetails(nameTag string) (RDSInstanceDetails, error) {
	details := RDSInstanceDetails{Tags: make(map[string]string)}

	var instanceTagList []*rds.Tag
	resp, err := r.getRDSInstancesByTag(nameTag)
	for _, instance := range resp {
		details.IDs = append(details.IDs, instance.DBInstanceIdentifier)
		details.ResourceNames = append(details.ResourceNames, instance.DBInstanceArn)
		instanceTagList = append(instanceTagList, instance.TagList...)
	}

	if err != nil {
		return details, err
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

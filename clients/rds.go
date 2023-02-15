package clients

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSResource struct {
	NameTag string
}

type RDSMonitor struct {
	Type string
}

var rdsSession *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var rdsClient = rds.New(rdsSession)

func (r *RDSResource) Start(startInput ResourceStartInput) error {
	DBInstanceIds, err := r.getRDSInstanceIdsByTag(r.NameTag)
	if err != nil {
		return err
	}

	for _, DBInstance := range DBInstanceIds {
		_, err = rdsClient.StartDBInstance(&rds.StartDBInstanceInput{
			DBInstanceIdentifier: DBInstance,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RDSResource) Stop(stopInput ResourceStopInput) error {
	DBInstanceIds, err := r.getRDSInstanceIdsByTag(r.NameTag)
	if err != nil {
		return err
	}
	for _, DBInstance := range DBInstanceIds {
		_, err = rdsClient.StopDBInstance(&rds.StopDBInstanceInput{
			DBInstanceIdentifier: DBInstance,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RDSResource) Status() (ResourceStatusOutput, error) {
	var rst ResourceStatusOutput

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

func (r *RDSResource) getRDSInstanceIdsByTag(nameTag string) ([]*string, error) {
	var DBInstanceIds []*string
	filteredInstances, err := r.getRDSInstancesByTag(nameTag)

	if err != nil {
		return nil, err
	}

	for _, instance := range filteredInstances {
		DBInstanceIds = append(DBInstanceIds, instance.DBInstanceIdentifier)
	}

	return DBInstanceIds, nil
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

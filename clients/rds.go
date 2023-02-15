package clients

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RDSResource struct {
	NameTag string
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

// func (r *RDSResource) getRDSInstanceIds(nameTag string) ([]*string, error) {
// 	var DBInstanceIds []*string
// 	// tagKey := fmt.Sprintf("tag:%s", defaultTagKey)
// 	// nameFilter := &rds.Filter{
// 	// 	Name:   &tagKey,
// 	// 	Values: []*string{&nameTag},
// 	// }

// 	resp, err := rdsClient.DescribeDBInstances(&rds.DescribeDBInstancesInput{
// 		Filters: []*rds.Filter{},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, DBInstance := range resp.DBInstances {
// 		DBInstanceIds = append(DBInstanceIds, DBInstance.DBInstanceIdentifier)
// 	}

// 	return DBInstanceIds, nil
// }

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

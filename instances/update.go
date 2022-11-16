package instances

// ScanAndUpdate grabs all instances that have a specific tag we use and creates a resource if a tag name is new to the system
// func ScanAndUpdate() error {
// 	resp, err := EC2Client.DescribeInstances(&ec2.DescribeInstancesInput{})
// 	if err != nil {
// 		return err
// 	}
//
// 	instanceInfo := generateInstanceInfo(resp)
// 	return updateResources(instanceInfo)
// }
//
// // generateInstanceInfo iterates over the instances and returns a combined map of name tag as key and formatted name tag as value
// func generateInstanceInfo(resp *ec2.DescribeInstancesOutput) map[string]string {
// 	instanceInfo := make(map[string]string)
//
// 	for _, reserv := range resp.Reservations {
// 		for _, inst := range reserv.Instances {
// 			for _, t := range inst.Tags {
// 				if *t.Key == defaultTagKey {
// 					nt := strings.ReplaceAll(*t.Value, "-", " ")
// 					instanceInfo[*t.Value] = cases.Title(language.English, cases.Compact).String(nt)
// 				}
// 			}
// 		}
// 	}
//
// 	return instanceInfo
// }
//
// // updateResources iterates over the resources and checks if a name tag is present in the map of name tags from EC2
// // then it creates resources or changes their archival flag based on that
// func updateResources(instanceInfo map[string]string) error {
// 	rcs, err := resources.GetResourceList()
// 	if err != nil {
// 		return err
// 	}
//
// 	for _, r := range rcs {
// 		_, ok := instanceInfo[r.NameTag]
// 		delete(instanceInfo, r.NameTag)
//
// 		if ok {
// 			if r.IsArchived {
// 				r.IsArchived = false
// 				err = resources.UpdateIsArchived(r)
// 				if err != nil {
// 					// TODO Log
// 					return err
// 				}
//
// 			}
// 			continue
// 		}
//
// 		r.IsArchived = true
// 		err = resources.UpdateIsArchived(r)
// 		if err != nil {
// 			// TODO Log
// 			return err
// 		}
// 	}
//
// 	for nameTag, nameTagDisplay := range instanceInfo {
// 		rce := resources.Resource{Name: nameTagDisplay, NameTag: nameTag}
// 		err = resources.CreateResource(rce)
// 		if err != nil {
// 			// TODO Log
// 			return err
// 		}
// 	}
//
// 	return nil
// }

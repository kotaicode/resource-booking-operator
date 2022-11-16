package instances

// StatusAll is a wrapper around Status. It returns a list of resource status summaries, which are used by the client app to display to users in a friendly way.
// func StatusAll() (map[int32]ResourceStatus, error) {
// 	statusList := make(map[int32]ResourceStatus)
//
// 	// TODO Here
// 	// rcs, _ := resources.GetResourceList()
// 	for _, resource := range rcs {
// 		rce := Resource{NameTag: resource.NameTag, IsArchived: resource.IsArchived}
//
// 		status, err := rce.Status()
// 		if err != nil {
// 			return statusList, err
// 		}
// 		statusList[resource.ID] = status
// 	}
//
// 	return statusList, nil
// }

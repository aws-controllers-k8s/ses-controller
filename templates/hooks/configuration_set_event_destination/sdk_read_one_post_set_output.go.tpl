	eventDestination := getEventDestination(ko, resp)
	if eventDestination == nil {
		return nil, ackerr.NotFound
	}
	ko.Spec.EventDestination = eventDestination

	if r.ko.Spec.EventDestination == nil || r.ko.Spec.EventDestination.Name == nil {
		return nil, ackerr.NotFound
	}

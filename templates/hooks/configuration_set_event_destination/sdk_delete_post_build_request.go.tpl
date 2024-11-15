	if eventDestination := r.ko.Spec.EventDestination; eventDestination != nil {
		input.EventDestinationName = eventDestination.Name
	}

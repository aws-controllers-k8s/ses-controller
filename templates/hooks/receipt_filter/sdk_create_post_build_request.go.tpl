	if input.Filter != nil {
		input.Filter.Name = desired.ko.Spec.Name
	} else {
		input.Filter = &svcsdk.ReceiptFilter{
			Name: desired.ko.Spec.Name,
		}
	}

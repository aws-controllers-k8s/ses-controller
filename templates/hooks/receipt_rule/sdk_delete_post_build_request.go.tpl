	if rule := r.ko.Spec.Rule; rule != nil {
		input.RuleName = rule.Name
	}

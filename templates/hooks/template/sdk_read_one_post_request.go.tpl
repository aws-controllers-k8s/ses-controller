	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == svcsdk.ErrCodeTemplateDoesNotExistException {
			rm.metrics.RecordAPICall("READ_ONE", "GetTemplate", err)
			return nil, ackerr.NotFound
		}
	}

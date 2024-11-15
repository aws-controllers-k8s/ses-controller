	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == svcsdk.ErrCodeCustomVerificationEmailTemplateDoesNotExistException {
			rm.metrics.RecordAPICall("READ_ONE", "GetCustomVerificationEmailTemplate", err)
			return nil, ackerr.NotFound
		}
	}

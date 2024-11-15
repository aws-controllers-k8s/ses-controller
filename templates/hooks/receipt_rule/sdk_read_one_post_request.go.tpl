	_ = resp
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == svcsdk.ErrCodeConfigurationSetDoesNotExistException {
			rm.metrics.RecordAPICall("READ_ONE", "DescribeReceiptRule", err)
			return nil, ackerr.NotFound
		}
	}

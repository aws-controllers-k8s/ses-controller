	_ = resp
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == svcsdk.ErrCodeRuleSetDoesNotExistException {
			rm.metrics.RecordAPICall("READ_ONE", "DescribeReceiptRuleSet", err)
			return nil, ackerr.NotFound
		}
	}

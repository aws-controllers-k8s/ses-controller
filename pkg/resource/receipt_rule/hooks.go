// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package receipt_rule

import (
	"context"

	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/ses"

	svcapitypes "github.com/aws-controllers-k8s/ses-controller/apis/v1alpha1"
)

func setReceiptRule(ko *svcapitypes.ReceiptRule, resp *svcsdk.DescribeReceiptRuleOutput) {
	ko.Spec.Rule = setResourceReceiptRule(resp.Rule)
}

func (rm *resourceManager) customFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if r.ko.Spec.Rule == nil || r.ko.Spec.Rule.Name == nil || r.ko.Spec.RuleSetName == nil {
		return nil, ackerr.NotFound
	}

	input := &svcsdk.DescribeReceiptRuleInput{
		RuleSetName: r.ko.Spec.RuleSetName,
	}
	if rule := r.ko.Spec.Rule; rule != nil {
		input.RuleName = rule.Name
	}

	var resp *svcsdk.DescribeReceiptRuleOutput
	resp, err = rm.sdkapi.DescribeReceiptRuleWithContext(ctx, input)

	rm.metrics.RecordAPICall("READ_ONE", "DescribeReceiptRule", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == svcsdk.ErrCodeRuleDoesNotExistException {
			return nil, ackerr.NotFound
		}
		if reqErr, ok := ackerr.AWSRequestFailure(err); ok && reqErr.StatusCode() == 404 {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	rm.setStatusDefaults(ko)
	setReceiptRule(ko, resp)

	return &resource{ko}, nil
}

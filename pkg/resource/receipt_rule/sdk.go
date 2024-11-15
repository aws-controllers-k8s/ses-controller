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

// Code generated by ack-generate. DO NOT EDIT.

package receipt_rule

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/ses"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/ses-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.SES{}
	_ = &svcapitypes.ReceiptRule{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	return rm.customFind(ctx, r)
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer func() {
		exit(err)
	}()
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateReceiptRuleOutput
	_ = resp
	resp, err = rm.sdkapi.CreateReceiptRuleWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateReceiptRule", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateReceiptRuleInput, error) {
	res := &svcsdk.CreateReceiptRuleInput{}

	if r.ko.Spec.After != nil {
		res.SetAfter(*r.ko.Spec.After)
	}
	if r.ko.Spec.Rule != nil {
		f1 := &svcsdk.ReceiptRule{}
		if r.ko.Spec.Rule.Actions != nil {
			f1f0 := []*svcsdk.ReceiptAction{}
			for _, f1f0iter := range r.ko.Spec.Rule.Actions {
				f1f0elem := &svcsdk.ReceiptAction{}
				if f1f0iter.AddHeaderAction != nil {
					f1f0elemf0 := &svcsdk.AddHeaderAction{}
					if f1f0iter.AddHeaderAction.HeaderName != nil {
						f1f0elemf0.SetHeaderName(*f1f0iter.AddHeaderAction.HeaderName)
					}
					if f1f0iter.AddHeaderAction.HeaderValue != nil {
						f1f0elemf0.SetHeaderValue(*f1f0iter.AddHeaderAction.HeaderValue)
					}
					f1f0elem.SetAddHeaderAction(f1f0elemf0)
				}
				if f1f0iter.BounceAction != nil {
					f1f0elemf1 := &svcsdk.BounceAction{}
					if f1f0iter.BounceAction.Message != nil {
						f1f0elemf1.SetMessage(*f1f0iter.BounceAction.Message)
					}
					if f1f0iter.BounceAction.Sender != nil {
						f1f0elemf1.SetSender(*f1f0iter.BounceAction.Sender)
					}
					if f1f0iter.BounceAction.SmtpReplyCode != nil {
						f1f0elemf1.SetSmtpReplyCode(*f1f0iter.BounceAction.SmtpReplyCode)
					}
					if f1f0iter.BounceAction.StatusCode != nil {
						f1f0elemf1.SetStatusCode(*f1f0iter.BounceAction.StatusCode)
					}
					if f1f0iter.BounceAction.TopicARN != nil {
						f1f0elemf1.SetTopicArn(*f1f0iter.BounceAction.TopicARN)
					}
					f1f0elem.SetBounceAction(f1f0elemf1)
				}
				if f1f0iter.LambdaAction != nil {
					f1f0elemf2 := &svcsdk.LambdaAction{}
					if f1f0iter.LambdaAction.FunctionARN != nil {
						f1f0elemf2.SetFunctionArn(*f1f0iter.LambdaAction.FunctionARN)
					}
					if f1f0iter.LambdaAction.InvocationType != nil {
						f1f0elemf2.SetInvocationType(*f1f0iter.LambdaAction.InvocationType)
					}
					if f1f0iter.LambdaAction.TopicARN != nil {
						f1f0elemf2.SetTopicArn(*f1f0iter.LambdaAction.TopicARN)
					}
					f1f0elem.SetLambdaAction(f1f0elemf2)
				}
				if f1f0iter.S3Action != nil {
					f1f0elemf3 := &svcsdk.S3Action{}
					if f1f0iter.S3Action.BucketName != nil {
						f1f0elemf3.SetBucketName(*f1f0iter.S3Action.BucketName)
					}
					if f1f0iter.S3Action.KMSKeyARN != nil {
						f1f0elemf3.SetKmsKeyArn(*f1f0iter.S3Action.KMSKeyARN)
					}
					if f1f0iter.S3Action.ObjectKeyPrefix != nil {
						f1f0elemf3.SetObjectKeyPrefix(*f1f0iter.S3Action.ObjectKeyPrefix)
					}
					if f1f0iter.S3Action.TopicARN != nil {
						f1f0elemf3.SetTopicArn(*f1f0iter.S3Action.TopicARN)
					}
					f1f0elem.SetS3Action(f1f0elemf3)
				}
				if f1f0iter.SNSAction != nil {
					f1f0elemf4 := &svcsdk.SNSAction{}
					if f1f0iter.SNSAction.Encoding != nil {
						f1f0elemf4.SetEncoding(*f1f0iter.SNSAction.Encoding)
					}
					if f1f0iter.SNSAction.TopicARN != nil {
						f1f0elemf4.SetTopicArn(*f1f0iter.SNSAction.TopicARN)
					}
					f1f0elem.SetSNSAction(f1f0elemf4)
				}
				if f1f0iter.StopAction != nil {
					f1f0elemf5 := &svcsdk.StopAction{}
					if f1f0iter.StopAction.Scope != nil {
						f1f0elemf5.SetScope(*f1f0iter.StopAction.Scope)
					}
					if f1f0iter.StopAction.TopicARN != nil {
						f1f0elemf5.SetTopicArn(*f1f0iter.StopAction.TopicARN)
					}
					f1f0elem.SetStopAction(f1f0elemf5)
				}
				if f1f0iter.WorkmailAction != nil {
					f1f0elemf6 := &svcsdk.WorkmailAction{}
					if f1f0iter.WorkmailAction.OrganizationARN != nil {
						f1f0elemf6.SetOrganizationArn(*f1f0iter.WorkmailAction.OrganizationARN)
					}
					if f1f0iter.WorkmailAction.TopicARN != nil {
						f1f0elemf6.SetTopicArn(*f1f0iter.WorkmailAction.TopicARN)
					}
					f1f0elem.SetWorkmailAction(f1f0elemf6)
				}
				f1f0 = append(f1f0, f1f0elem)
			}
			f1.SetActions(f1f0)
		}
		if r.ko.Spec.Rule.Enabled != nil {
			f1.SetEnabled(*r.ko.Spec.Rule.Enabled)
		}
		if r.ko.Spec.Rule.Name != nil {
			f1.SetName(*r.ko.Spec.Rule.Name)
		}
		if r.ko.Spec.Rule.Recipients != nil {
			f1f3 := []*string{}
			for _, f1f3iter := range r.ko.Spec.Rule.Recipients {
				var f1f3elem string
				f1f3elem = *f1f3iter
				f1f3 = append(f1f3, &f1f3elem)
			}
			f1.SetRecipients(f1f3)
		}
		if r.ko.Spec.Rule.ScanEnabled != nil {
			f1.SetScanEnabled(*r.ko.Spec.Rule.ScanEnabled)
		}
		if r.ko.Spec.Rule.TLSPolicy != nil {
			f1.SetTlsPolicy(*r.ko.Spec.Rule.TLSPolicy)
		}
		res.SetRule(f1)
	}
	if r.ko.Spec.RuleSetName != nil {
		res.SetRuleSetName(*r.ko.Spec.RuleSetName)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()
	if immutableFieldChanges := rm.getImmutableFieldChanges(delta); len(immutableFieldChanges) > 0 {
		msg := fmt.Sprintf("Immutable Spec fields have been modified: %s", strings.Join(immutableFieldChanges, ","))
		return nil, ackerr.NewTerminalError(fmt.Errorf(msg))
	}
	input, err := rm.newUpdateRequestPayload(ctx, desired, delta)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.UpdateReceiptRuleOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateReceiptRuleWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateReceiptRule", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (*svcsdk.UpdateReceiptRuleInput, error) {
	res := &svcsdk.UpdateReceiptRuleInput{}

	if r.ko.Spec.Rule != nil {
		f0 := &svcsdk.ReceiptRule{}
		if r.ko.Spec.Rule.Actions != nil {
			f0f0 := []*svcsdk.ReceiptAction{}
			for _, f0f0iter := range r.ko.Spec.Rule.Actions {
				f0f0elem := &svcsdk.ReceiptAction{}
				if f0f0iter.AddHeaderAction != nil {
					f0f0elemf0 := &svcsdk.AddHeaderAction{}
					if f0f0iter.AddHeaderAction.HeaderName != nil {
						f0f0elemf0.SetHeaderName(*f0f0iter.AddHeaderAction.HeaderName)
					}
					if f0f0iter.AddHeaderAction.HeaderValue != nil {
						f0f0elemf0.SetHeaderValue(*f0f0iter.AddHeaderAction.HeaderValue)
					}
					f0f0elem.SetAddHeaderAction(f0f0elemf0)
				}
				if f0f0iter.BounceAction != nil {
					f0f0elemf1 := &svcsdk.BounceAction{}
					if f0f0iter.BounceAction.Message != nil {
						f0f0elemf1.SetMessage(*f0f0iter.BounceAction.Message)
					}
					if f0f0iter.BounceAction.Sender != nil {
						f0f0elemf1.SetSender(*f0f0iter.BounceAction.Sender)
					}
					if f0f0iter.BounceAction.SmtpReplyCode != nil {
						f0f0elemf1.SetSmtpReplyCode(*f0f0iter.BounceAction.SmtpReplyCode)
					}
					if f0f0iter.BounceAction.StatusCode != nil {
						f0f0elemf1.SetStatusCode(*f0f0iter.BounceAction.StatusCode)
					}
					if f0f0iter.BounceAction.TopicARN != nil {
						f0f0elemf1.SetTopicArn(*f0f0iter.BounceAction.TopicARN)
					}
					f0f0elem.SetBounceAction(f0f0elemf1)
				}
				if f0f0iter.LambdaAction != nil {
					f0f0elemf2 := &svcsdk.LambdaAction{}
					if f0f0iter.LambdaAction.FunctionARN != nil {
						f0f0elemf2.SetFunctionArn(*f0f0iter.LambdaAction.FunctionARN)
					}
					if f0f0iter.LambdaAction.InvocationType != nil {
						f0f0elemf2.SetInvocationType(*f0f0iter.LambdaAction.InvocationType)
					}
					if f0f0iter.LambdaAction.TopicARN != nil {
						f0f0elemf2.SetTopicArn(*f0f0iter.LambdaAction.TopicARN)
					}
					f0f0elem.SetLambdaAction(f0f0elemf2)
				}
				if f0f0iter.S3Action != nil {
					f0f0elemf3 := &svcsdk.S3Action{}
					if f0f0iter.S3Action.BucketName != nil {
						f0f0elemf3.SetBucketName(*f0f0iter.S3Action.BucketName)
					}
					if f0f0iter.S3Action.KMSKeyARN != nil {
						f0f0elemf3.SetKmsKeyArn(*f0f0iter.S3Action.KMSKeyARN)
					}
					if f0f0iter.S3Action.ObjectKeyPrefix != nil {
						f0f0elemf3.SetObjectKeyPrefix(*f0f0iter.S3Action.ObjectKeyPrefix)
					}
					if f0f0iter.S3Action.TopicARN != nil {
						f0f0elemf3.SetTopicArn(*f0f0iter.S3Action.TopicARN)
					}
					f0f0elem.SetS3Action(f0f0elemf3)
				}
				if f0f0iter.SNSAction != nil {
					f0f0elemf4 := &svcsdk.SNSAction{}
					if f0f0iter.SNSAction.Encoding != nil {
						f0f0elemf4.SetEncoding(*f0f0iter.SNSAction.Encoding)
					}
					if f0f0iter.SNSAction.TopicARN != nil {
						f0f0elemf4.SetTopicArn(*f0f0iter.SNSAction.TopicARN)
					}
					f0f0elem.SetSNSAction(f0f0elemf4)
				}
				if f0f0iter.StopAction != nil {
					f0f0elemf5 := &svcsdk.StopAction{}
					if f0f0iter.StopAction.Scope != nil {
						f0f0elemf5.SetScope(*f0f0iter.StopAction.Scope)
					}
					if f0f0iter.StopAction.TopicARN != nil {
						f0f0elemf5.SetTopicArn(*f0f0iter.StopAction.TopicARN)
					}
					f0f0elem.SetStopAction(f0f0elemf5)
				}
				if f0f0iter.WorkmailAction != nil {
					f0f0elemf6 := &svcsdk.WorkmailAction{}
					if f0f0iter.WorkmailAction.OrganizationARN != nil {
						f0f0elemf6.SetOrganizationArn(*f0f0iter.WorkmailAction.OrganizationARN)
					}
					if f0f0iter.WorkmailAction.TopicARN != nil {
						f0f0elemf6.SetTopicArn(*f0f0iter.WorkmailAction.TopicARN)
					}
					f0f0elem.SetWorkmailAction(f0f0elemf6)
				}
				f0f0 = append(f0f0, f0f0elem)
			}
			f0.SetActions(f0f0)
		}
		if r.ko.Spec.Rule.Enabled != nil {
			f0.SetEnabled(*r.ko.Spec.Rule.Enabled)
		}
		if r.ko.Spec.Rule.Name != nil {
			f0.SetName(*r.ko.Spec.Rule.Name)
		}
		if r.ko.Spec.Rule.Recipients != nil {
			f0f3 := []*string{}
			for _, f0f3iter := range r.ko.Spec.Rule.Recipients {
				var f0f3elem string
				f0f3elem = *f0f3iter
				f0f3 = append(f0f3, &f0f3elem)
			}
			f0.SetRecipients(f0f3)
		}
		if r.ko.Spec.Rule.ScanEnabled != nil {
			f0.SetScanEnabled(*r.ko.Spec.Rule.ScanEnabled)
		}
		if r.ko.Spec.Rule.TLSPolicy != nil {
			f0.SetTlsPolicy(*r.ko.Spec.Rule.TLSPolicy)
		}
		res.SetRule(f0)
	}
	if r.ko.Spec.RuleSetName != nil {
		res.SetRuleSetName(*r.ko.Spec.RuleSetName)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	if rule := r.ko.Spec.Rule; rule != nil {
		input.RuleName = rule.Name
	}

	var resp *svcsdk.DeleteReceiptRuleOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteReceiptRuleWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteReceiptRule", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteReceiptRuleInput, error) {
	res := &svcsdk.DeleteReceiptRuleInput{}

	if r.ko.Spec.RuleSetName != nil {
		res.SetRuleSetName(*r.ko.Spec.RuleSetName)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.ReceiptRule,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "AlreadyExists",
		"InvalidLambdaFunction",
		"InvalidS3Configuration",
		"InvalidSnsTopic",
		"RuleSetDoesNotExist":
		return true
	default:
		return false
	}
}

// getImmutableFieldChanges returns list of immutable fields from the
func (rm *resourceManager) getImmutableFieldChanges(
	delta *ackcompare.Delta,
) []string {
	var fields []string
	if delta.DifferentAt("Spec.Rule.Name") {
		fields = append(fields, "Rule.Name")
	}
	if delta.DifferentAt("Spec.RuleSetName") {
		fields = append(fields, "RuleSetName")
	}

	return fields
}

// setReceiptRule sets a resource ReceiptRule type
// given the SDK type.
func setResourceReceiptRule(
	resp *svcsdk.ReceiptRule,
) *svcapitypes.ReceiptRule_SDK {
	res := &svcapitypes.ReceiptRule_SDK{}

	if resp.Actions != nil {
		resf0 := []*svcapitypes.ReceiptAction{}
		for _, resf0iter := range resp.Actions {
			resf0elem := &svcapitypes.ReceiptAction{}
			if resf0iter.AddHeaderAction != nil {
				resf0elemf0 := &svcapitypes.AddHeaderAction{}
				if resf0iter.AddHeaderAction.HeaderName != nil {
					resf0elemf0.HeaderName = resf0iter.AddHeaderAction.HeaderName
				}
				if resf0iter.AddHeaderAction.HeaderValue != nil {
					resf0elemf0.HeaderValue = resf0iter.AddHeaderAction.HeaderValue
				}
				resf0elem.AddHeaderAction = resf0elemf0
			}
			if resf0iter.BounceAction != nil {
				resf0elemf1 := &svcapitypes.BounceAction{}
				if resf0iter.BounceAction.Message != nil {
					resf0elemf1.Message = resf0iter.BounceAction.Message
				}
				if resf0iter.BounceAction.Sender != nil {
					resf0elemf1.Sender = resf0iter.BounceAction.Sender
				}
				if resf0iter.BounceAction.SmtpReplyCode != nil {
					resf0elemf1.SmtpReplyCode = resf0iter.BounceAction.SmtpReplyCode
				}
				if resf0iter.BounceAction.StatusCode != nil {
					resf0elemf1.StatusCode = resf0iter.BounceAction.StatusCode
				}
				if resf0iter.BounceAction.TopicArn != nil {
					resf0elemf1.TopicARN = resf0iter.BounceAction.TopicArn
				}
				resf0elem.BounceAction = resf0elemf1
			}
			if resf0iter.LambdaAction != nil {
				resf0elemf2 := &svcapitypes.LambdaAction{}
				if resf0iter.LambdaAction.FunctionArn != nil {
					resf0elemf2.FunctionARN = resf0iter.LambdaAction.FunctionArn
				}
				if resf0iter.LambdaAction.InvocationType != nil {
					resf0elemf2.InvocationType = resf0iter.LambdaAction.InvocationType
				}
				if resf0iter.LambdaAction.TopicArn != nil {
					resf0elemf2.TopicARN = resf0iter.LambdaAction.TopicArn
				}
				resf0elem.LambdaAction = resf0elemf2
			}
			if resf0iter.S3Action != nil {
				resf0elemf3 := &svcapitypes.S3Action{}
				if resf0iter.S3Action.BucketName != nil {
					resf0elemf3.BucketName = resf0iter.S3Action.BucketName
				}
				if resf0iter.S3Action.KmsKeyArn != nil {
					resf0elemf3.KMSKeyARN = resf0iter.S3Action.KmsKeyArn
				}
				if resf0iter.S3Action.ObjectKeyPrefix != nil {
					resf0elemf3.ObjectKeyPrefix = resf0iter.S3Action.ObjectKeyPrefix
				}
				if resf0iter.S3Action.TopicArn != nil {
					resf0elemf3.TopicARN = resf0iter.S3Action.TopicArn
				}
				resf0elem.S3Action = resf0elemf3
			}
			if resf0iter.SNSAction != nil {
				resf0elemf4 := &svcapitypes.SNSAction{}
				if resf0iter.SNSAction.Encoding != nil {
					resf0elemf4.Encoding = resf0iter.SNSAction.Encoding
				}
				if resf0iter.SNSAction.TopicArn != nil {
					resf0elemf4.TopicARN = resf0iter.SNSAction.TopicArn
				}
				resf0elem.SNSAction = resf0elemf4
			}
			if resf0iter.StopAction != nil {
				resf0elemf5 := &svcapitypes.StopAction{}
				if resf0iter.StopAction.Scope != nil {
					resf0elemf5.Scope = resf0iter.StopAction.Scope
				}
				if resf0iter.StopAction.TopicArn != nil {
					resf0elemf5.TopicARN = resf0iter.StopAction.TopicArn
				}
				resf0elem.StopAction = resf0elemf5
			}
			if resf0iter.WorkmailAction != nil {
				resf0elemf6 := &svcapitypes.WorkmailAction{}
				if resf0iter.WorkmailAction.OrganizationArn != nil {
					resf0elemf6.OrganizationARN = resf0iter.WorkmailAction.OrganizationArn
				}
				if resf0iter.WorkmailAction.TopicArn != nil {
					resf0elemf6.TopicARN = resf0iter.WorkmailAction.TopicArn
				}
				resf0elem.WorkmailAction = resf0elemf6
			}
			resf0 = append(resf0, resf0elem)
		}
		res.Actions = resf0
	}
	if resp.Enabled != nil {
		res.Enabled = resp.Enabled
	}
	if resp.Name != nil {
		res.Name = resp.Name
	}
	if resp.Recipients != nil {
		resf3 := []*string{}
		for _, resf3iter := range resp.Recipients {
			var resf3elem string
			resf3elem = *resf3iter
			resf3 = append(resf3, &resf3elem)
		}
		res.Recipients = resf3
	}
	if resp.ScanEnabled != nil {
		res.ScanEnabled = resp.ScanEnabled
	}
	if resp.TlsPolicy != nil {
		res.TLSPolicy = resp.TlsPolicy
	}

	return res
}

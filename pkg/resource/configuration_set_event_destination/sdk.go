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

package configuration_set_event_destination

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
	_ = &svcapitypes.ConfigurationSetEventDestination{}
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
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
	if r.ko.Spec.EventDestination == nil || r.ko.Spec.EventDestination.Name == nil {
		return nil, ackerr.NotFound
	}

	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}
	input.SetConfigurationSetAttributeNames(aws.StringSlice([]string{svcsdk.ConfigurationSetAttributeEventDestinations}))

	var resp *svcsdk.DescribeConfigurationSetOutput
	resp, err = rm.sdkapi.DescribeConfigurationSetWithContext(ctx, input)
	_ = resp
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == svcsdk.ErrCodeConfigurationSetDoesNotExistException {
			rm.metrics.RecordAPICall("READ_ONE", "DescribeConfigurationSet", err)
			return nil, ackerr.NotFound
		}
	}

	rm.metrics.RecordAPICall("READ_ONE", "DescribeConfigurationSet", err)
	if err != nil {
		if reqErr, ok := ackerr.AWSRequestFailure(err); ok && reqErr.StatusCode() == 404 {
			return nil, ackerr.NotFound
		}
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "UNKNOWN" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	rm.setStatusDefaults(ko)
	eventDestination := getEventDestination(ko, resp)
	if eventDestination == nil {
		return nil, ackerr.NotFound
	}
	ko.Spec.EventDestination = eventDestination

	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
	return r.ko.Spec.ConfigurationSetName == nil

}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.DescribeConfigurationSetInput, error) {
	res := &svcsdk.DescribeConfigurationSetInput{}

	if r.ko.Spec.ConfigurationSetName != nil {
		res.SetConfigurationSetName(*r.ko.Spec.ConfigurationSetName)
	}

	return res, nil
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

	var resp *svcsdk.CreateConfigurationSetEventDestinationOutput
	_ = resp
	resp, err = rm.sdkapi.CreateConfigurationSetEventDestinationWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateConfigurationSetEventDestination", err)
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
) (*svcsdk.CreateConfigurationSetEventDestinationInput, error) {
	res := &svcsdk.CreateConfigurationSetEventDestinationInput{}

	if r.ko.Spec.ConfigurationSetName != nil {
		res.SetConfigurationSetName(*r.ko.Spec.ConfigurationSetName)
	}
	if r.ko.Spec.EventDestination != nil {
		f1 := &svcsdk.EventDestination{}
		if r.ko.Spec.EventDestination.CloudWatchDestination != nil {
			f1f0 := &svcsdk.CloudWatchDestination{}
			if r.ko.Spec.EventDestination.CloudWatchDestination.DimensionConfigurations != nil {
				f1f0f0 := []*svcsdk.CloudWatchDimensionConfiguration{}
				for _, f1f0f0iter := range r.ko.Spec.EventDestination.CloudWatchDestination.DimensionConfigurations {
					f1f0f0elem := &svcsdk.CloudWatchDimensionConfiguration{}
					if f1f0f0iter.DefaultDimensionValue != nil {
						f1f0f0elem.SetDefaultDimensionValue(*f1f0f0iter.DefaultDimensionValue)
					}
					if f1f0f0iter.DimensionName != nil {
						f1f0f0elem.SetDimensionName(*f1f0f0iter.DimensionName)
					}
					if f1f0f0iter.DimensionValueSource != nil {
						f1f0f0elem.SetDimensionValueSource(*f1f0f0iter.DimensionValueSource)
					}
					f1f0f0 = append(f1f0f0, f1f0f0elem)
				}
				f1f0.SetDimensionConfigurations(f1f0f0)
			}
			f1.SetCloudWatchDestination(f1f0)
		}
		if r.ko.Spec.EventDestination.Enabled != nil {
			f1.SetEnabled(*r.ko.Spec.EventDestination.Enabled)
		}
		if r.ko.Spec.EventDestination.KinesisFirehoseDestination != nil {
			f1f2 := &svcsdk.KinesisFirehoseDestination{}
			if r.ko.Spec.EventDestination.KinesisFirehoseDestination.DeliveryStreamARN != nil {
				f1f2.SetDeliveryStreamARN(*r.ko.Spec.EventDestination.KinesisFirehoseDestination.DeliveryStreamARN)
			}
			if r.ko.Spec.EventDestination.KinesisFirehoseDestination.IAMRoleARN != nil {
				f1f2.SetIAMRoleARN(*r.ko.Spec.EventDestination.KinesisFirehoseDestination.IAMRoleARN)
			}
			f1.SetKinesisFirehoseDestination(f1f2)
		}
		if r.ko.Spec.EventDestination.MatchingEventTypes != nil {
			f1f3 := []*string{}
			for _, f1f3iter := range r.ko.Spec.EventDestination.MatchingEventTypes {
				var f1f3elem string
				f1f3elem = *f1f3iter
				f1f3 = append(f1f3, &f1f3elem)
			}
			f1.SetMatchingEventTypes(f1f3)
		}
		if r.ko.Spec.EventDestination.Name != nil {
			f1.SetName(*r.ko.Spec.EventDestination.Name)
		}
		if r.ko.Spec.EventDestination.SNSDestination != nil {
			f1f5 := &svcsdk.SNSDestination{}
			if r.ko.Spec.EventDestination.SNSDestination.TopicARN != nil {
				f1f5.SetTopicARN(*r.ko.Spec.EventDestination.SNSDestination.TopicARN)
			}
			f1.SetSNSDestination(f1f5)
		}
		res.SetEventDestination(f1)
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

	var resp *svcsdk.UpdateConfigurationSetEventDestinationOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateConfigurationSetEventDestinationWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateConfigurationSetEventDestination", err)
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
) (*svcsdk.UpdateConfigurationSetEventDestinationInput, error) {
	res := &svcsdk.UpdateConfigurationSetEventDestinationInput{}

	if r.ko.Spec.ConfigurationSetName != nil {
		res.SetConfigurationSetName(*r.ko.Spec.ConfigurationSetName)
	}
	if r.ko.Spec.EventDestination != nil {
		f1 := &svcsdk.EventDestination{}
		if r.ko.Spec.EventDestination.CloudWatchDestination != nil {
			f1f0 := &svcsdk.CloudWatchDestination{}
			if r.ko.Spec.EventDestination.CloudWatchDestination.DimensionConfigurations != nil {
				f1f0f0 := []*svcsdk.CloudWatchDimensionConfiguration{}
				for _, f1f0f0iter := range r.ko.Spec.EventDestination.CloudWatchDestination.DimensionConfigurations {
					f1f0f0elem := &svcsdk.CloudWatchDimensionConfiguration{}
					if f1f0f0iter.DefaultDimensionValue != nil {
						f1f0f0elem.SetDefaultDimensionValue(*f1f0f0iter.DefaultDimensionValue)
					}
					if f1f0f0iter.DimensionName != nil {
						f1f0f0elem.SetDimensionName(*f1f0f0iter.DimensionName)
					}
					if f1f0f0iter.DimensionValueSource != nil {
						f1f0f0elem.SetDimensionValueSource(*f1f0f0iter.DimensionValueSource)
					}
					f1f0f0 = append(f1f0f0, f1f0f0elem)
				}
				f1f0.SetDimensionConfigurations(f1f0f0)
			}
			f1.SetCloudWatchDestination(f1f0)
		}
		if r.ko.Spec.EventDestination.Enabled != nil {
			f1.SetEnabled(*r.ko.Spec.EventDestination.Enabled)
		}
		if r.ko.Spec.EventDestination.KinesisFirehoseDestination != nil {
			f1f2 := &svcsdk.KinesisFirehoseDestination{}
			if r.ko.Spec.EventDestination.KinesisFirehoseDestination.DeliveryStreamARN != nil {
				f1f2.SetDeliveryStreamARN(*r.ko.Spec.EventDestination.KinesisFirehoseDestination.DeliveryStreamARN)
			}
			if r.ko.Spec.EventDestination.KinesisFirehoseDestination.IAMRoleARN != nil {
				f1f2.SetIAMRoleARN(*r.ko.Spec.EventDestination.KinesisFirehoseDestination.IAMRoleARN)
			}
			f1.SetKinesisFirehoseDestination(f1f2)
		}
		if r.ko.Spec.EventDestination.MatchingEventTypes != nil {
			f1f3 := []*string{}
			for _, f1f3iter := range r.ko.Spec.EventDestination.MatchingEventTypes {
				var f1f3elem string
				f1f3elem = *f1f3iter
				f1f3 = append(f1f3, &f1f3elem)
			}
			f1.SetMatchingEventTypes(f1f3)
		}
		if r.ko.Spec.EventDestination.Name != nil {
			f1.SetName(*r.ko.Spec.EventDestination.Name)
		}
		if r.ko.Spec.EventDestination.SNSDestination != nil {
			f1f5 := &svcsdk.SNSDestination{}
			if r.ko.Spec.EventDestination.SNSDestination.TopicARN != nil {
				f1f5.SetTopicARN(*r.ko.Spec.EventDestination.SNSDestination.TopicARN)
			}
			f1.SetSNSDestination(f1f5)
		}
		res.SetEventDestination(f1)
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
	if eventDestination := r.ko.Spec.EventDestination; eventDestination != nil {
		input.EventDestinationName = eventDestination.Name
	}

	var resp *svcsdk.DeleteConfigurationSetEventDestinationOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteConfigurationSetEventDestinationWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteConfigurationSetEventDestination", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteConfigurationSetEventDestinationInput, error) {
	res := &svcsdk.DeleteConfigurationSetEventDestinationInput{}

	if r.ko.Spec.ConfigurationSetName != nil {
		res.SetConfigurationSetName(*r.ko.Spec.ConfigurationSetName)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.ConfigurationSetEventDestination,
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
	case "ConfigurationSetDoesNotExist",
		"EventDestinationAlreadyExists",
		"InvalidCloudWatchDestination",
		"InvalidFirehoseDestination",
		"InvalidSNSDestination",
		"LimitExceeded":
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
	if delta.DifferentAt("Spec.ConfigurationSetName") {
		fields = append(fields, "ConfigurationSetName")
	}
	if delta.DifferentAt("Spec.EventDestination.Name") {
		fields = append(fields, "EventDestination.Name")
	}

	return fields
}

// setEventDestination sets a resource EventDestination type
// given the SDK type.
func setResourceEventDestination(
	resp *svcsdk.EventDestination,
) *svcapitypes.EventDestination {
	res := &svcapitypes.EventDestination{}

	if resp.CloudWatchDestination != nil {
		resf0 := &svcapitypes.CloudWatchDestination{}
		if resp.CloudWatchDestination.DimensionConfigurations != nil {
			resf0f0 := []*svcapitypes.CloudWatchDimensionConfiguration{}
			for _, resf0f0iter := range resp.CloudWatchDestination.DimensionConfigurations {
				resf0f0elem := &svcapitypes.CloudWatchDimensionConfiguration{}
				if resf0f0iter.DefaultDimensionValue != nil {
					resf0f0elem.DefaultDimensionValue = resf0f0iter.DefaultDimensionValue
				}
				if resf0f0iter.DimensionName != nil {
					resf0f0elem.DimensionName = resf0f0iter.DimensionName
				}
				if resf0f0iter.DimensionValueSource != nil {
					resf0f0elem.DimensionValueSource = resf0f0iter.DimensionValueSource
				}
				resf0f0 = append(resf0f0, resf0f0elem)
			}
			resf0.DimensionConfigurations = resf0f0
		}
		res.CloudWatchDestination = resf0
	}
	if resp.Enabled != nil {
		res.Enabled = resp.Enabled
	}
	if resp.KinesisFirehoseDestination != nil {
		resf2 := &svcapitypes.KinesisFirehoseDestination{}
		if resp.KinesisFirehoseDestination.DeliveryStreamARN != nil {
			resf2.DeliveryStreamARN = resp.KinesisFirehoseDestination.DeliveryStreamARN
		}
		if resp.KinesisFirehoseDestination.IAMRoleARN != nil {
			resf2.IAMRoleARN = resp.KinesisFirehoseDestination.IAMRoleARN
		}
		res.KinesisFirehoseDestination = resf2
	}
	if resp.MatchingEventTypes != nil {
		resf3 := []*string{}
		for _, resf3iter := range resp.MatchingEventTypes {
			var resf3elem string
			resf3elem = *resf3iter
			resf3 = append(resf3, &resf3elem)
		}
		res.MatchingEventTypes = resf3
	}
	if resp.Name != nil {
		res.Name = resp.Name
	}
	if resp.SNSDestination != nil {
		resf5 := &svcapitypes.SNSDestination{}
		if resp.SNSDestination.TopicARN != nil {
			resf5.TopicARN = resp.SNSDestination.TopicARN
		}
		res.SNSDestination = resf5
	}

	return res
}
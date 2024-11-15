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

package configuration_set_event_destination

import (
	svcsdk "github.com/aws/aws-sdk-go/service/ses"

	svcapitypes "github.com/aws-controllers-k8s/ses-controller/apis/v1alpha1"
)

func getEventDestination(ko *svcapitypes.ConfigurationSetEventDestination, resp *svcsdk.DescribeConfigurationSetOutput) *svcapitypes.EventDestination {
	for _, eventDestination := range resp.EventDestinations {
		if *eventDestination.Name == *ko.Spec.EventDestination.Name {
			return setResourceEventDestination(resp.EventDestinations[0])
		}
	}
	return nil
}

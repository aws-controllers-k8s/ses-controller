# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the ConfigurationSetEventDestination resource
"""

import boto3
import logging
import time

import pytest
from functools import partial

from typing import Dict, Tuple
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest.k8s import condition
from e2e import service_marker, CRD_GROUP, CRD_VERSION, SERVICE_NAME, load_ses_resource
from e2e.common.waiter import wait_until_deleted, MAX_WAIT_FOR_SYNCED_MINUTES
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources

from .configuration_set_test import simple_configuration_set


EVENT_DESTINATION_RESOURCE_PLURAL = "configurationseteventdestinations"

@pytest.fixture(scope='module')
def ses_client():
    return boto3.client(SERVICE_NAME)

def get_event_destination(ses_client, configuration_set_name, event_destination_name):
    try:
        res = ses_client.describe_configuration_set(ConfigurationSetName=configuration_set_name, ConfigurationSetAttributeNames=['eventDestinations'])
        event_destinations = [event_destination for event_destination in res['EventDestinations'] if event_destination['Name'] == event_destination_name]
        return event_destinations[0] if event_destinations else None
    except ses_client.exceptions.ConfigurationSetDoesNotExistException:
        return None

@pytest.fixture
def simple_event_destination(simple_configuration_set, ses_client) -> Tuple[k8s.CustomResourceReference, Dict]:
    event_destination_name = random_suffix_name('simple-event-destination', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['EVENT_DESTINATION_NAME'] = event_destination_name
    replacements['SNS_TOPIC_ARN'] = get_bootstrap_resources().SNSTopic.arn

    (ref, configuration_set_cr) = simple_configuration_set
    replacements['CONFIGURATION_SET_REF_NAME'] = configuration_set_cr['spec']['name']

    resource_data = load_ses_resource(
        'event_destination_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, EVENT_DESTINATION_RESOURCE_PLURAL,
        event_destination_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=10)

    assert cr is not None
    assert cr['status'] is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(get_event_destination, ses_client, configuration_set_cr['spec']['name'], event_destination_name))


@service_marker
@pytest.mark.canary
class TestEventDestination:
    def test_create_update_event_destination(self, simple_event_destination, ses_client):
        (ref, cr) = simple_event_destination
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )
        configuration_set_name = cr['spec']['configurationSetRef']['from']['name']
        event_destination_name = cr['spec']['eventDestination']['name']
        get_current_event_destination = partial(get_event_destination, ses_client, configuration_set_name, event_destination_name)
        assert get_current_event_destination() is not None

        updates = {
            'spec': {
                'eventDestination': {
                    'matchingEventTypes': ['send', 'reject', 'bounce', 'complaint'],
                    'enabled': True,
                }
            }
        }
        k8s.patch_custom_resource(ref, updates)
        time.sleep(10)
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )

        event_destination = get_current_event_destination()
        updated_event_destination = updates['spec']['eventDestination']
        assert event_destination == {
            'Name': event_destination_name,
            'MatchingEventTypes': updated_event_destination['matchingEventTypes'],
            'Enabled': updated_event_destination['enabled'],
            'SNSDestination': {
                'TopicARN': cr['eventDestination']['snsDestination']['topicARN'],
            },
        }

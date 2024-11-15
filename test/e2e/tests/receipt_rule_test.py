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

"""Integration tests for the ReceiptRule resource
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

from .receipt_rule_set_test import simple_receipt_rule_set

RECEIPT_RULE_RESOURCE_PLURAL = "receiptrules"

@pytest.fixture(scope='module')
def ses_client():
    return boto3.client(SERVICE_NAME)

def get_receipt_rule(ses_client, rule_set_name, rule_name):
    try:
        return ses_client.describe_receipt_rule(RuleSetName=rule_set_name, RuleName=rule_name)
    except ses_client.exceptions.RuleDoesNotExistException:
        return None

@pytest.fixture
def simple_receipt_rule(simple_receipt_rule_set, ses_client) -> Tuple[k8s.CustomResourceReference, Dict]:
    receipt_rule_name = random_suffix_name('simple-receipt-rule', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['RECEIPT_RULE_NAME'] = receipt_rule_name

    (ref, receipt_rule_set_cr) = simple_receipt_rule_set
    replacements['RULE_SET_REF_NAME'] = receipt_rule_set_cr['spec']['ruleSetName']

    resource_data = load_ses_resource(
        'receipt_rule_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RECEIPT_RULE_RESOURCE_PLURAL,
        receipt_rule_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=10)

    assert cr is not None
    assert cr['status'] is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(get_receipt_rule, ses_client, receipt_rule_set_cr['spec']['ruleSetName'], receipt_rule_name))


@service_marker
@pytest.mark.canary
class TestReceiptRule:
    def test_create_update_receipt_rule(self, simple_receipt_rule, ses_client):
        (ref, cr) = simple_receipt_rule
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )
        get_current_receipt_rule = partial(get_receipt_rule, ses_client, cr['spec']['ruleSetRef']['from']['name'], cr['spec']['rule']['name'])
        assert get_current_receipt_rule() is not None

        updates = {
            'spec': {
                'rule': {
                    'enabled': False,
                    'recipients': ['updated2@example.com', 'updated@example.com'],
                    'scanEnabled': True,
                    'tlsPolicy': 'Require',
                },
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

        receipt_rule = get_current_receipt_rule()
        updated_receipt_rule = updates['spec']['rule']
        assert receipt_rule['Rule'] == {
            'Name': cr['spec']['rule']['name'],
            'Enabled': updated_receipt_rule['enabled'],
            'Recipients': updated_receipt_rule['recipients'],
            'ScanEnabled': updated_receipt_rule['scanEnabled'],
            'TlsPolicy': updated_receipt_rule['tlsPolicy'],
            'Actions': [],
        }

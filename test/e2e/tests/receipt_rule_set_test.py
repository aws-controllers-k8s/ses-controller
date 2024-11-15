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

"""Integration tests for the ReceiptRuleSet resource
"""

import boto3
import logging

import pytest
from functools import partial

from typing import Dict, Tuple
from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest.k8s import condition
from e2e import service_marker, CRD_GROUP, CRD_VERSION, SERVICE_NAME, load_ses_resource
from e2e.common.waiter import wait_until_deleted, MAX_WAIT_FOR_SYNCED_MINUTES
from e2e.replacement_values import REPLACEMENT_VALUES


RECEIPT_RULE_SET_RESOURCE_PLURAL = "receiptrulesets"

@pytest.fixture(scope='module')
def ses_client():
    return boto3.client(SERVICE_NAME)

def get_receipt_rule_set(ses_client, receipt_rule_set_name):
    receipt_rule_sets = [rule_set for rule_set in ses_client.list_receipt_rule_sets()['RuleSets'] if rule_set['Name'] == receipt_rule_set_name]
    return receipt_rule_sets[0] if receipt_rule_sets else None

@pytest.fixture
def simple_receipt_rule_set(ses_client) -> Tuple[k8s.CustomResourceReference, Dict]:
    receipt_rule_set_name = random_suffix_name('simple-receipt-rule-set', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['RECEIPT_RULE_SET_NAME'] = receipt_rule_set_name

    resource_data = load_ses_resource(
        'receipt_rule_set_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RECEIPT_RULE_SET_RESOURCE_PLURAL,
        receipt_rule_set_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=10)

    assert cr is not None
    assert cr['status'] is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(get_receipt_rule_set, ses_client, receipt_rule_set_name))


@service_marker
@pytest.mark.canary
class TestReceiptRuleSet:
    def test_create_receipt_rule_set(self, simple_receipt_rule_set, ses_client):
        (ref, cr) = simple_receipt_rule_set
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )
        assert get_receipt_rule_set(ses_client, cr['spec']['ruleSetName']) is not None

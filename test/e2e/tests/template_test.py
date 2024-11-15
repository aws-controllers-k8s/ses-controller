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

"""Integration tests for the Template resource
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


TEMPLATE_RESOURCE_PLURAL = "templates"

@pytest.fixture(scope='module')
def ses_client():
    return boto3.client(SERVICE_NAME)

def get_template(ses_client, template_name):
    try:
        return ses_client.get_template(TemplateName=template_name)
    except ses_client.exceptions.TemplateDoesNotExistException:
        return None

@pytest.fixture
def simple_template(ses_client) -> Tuple[k8s.CustomResourceReference, Dict]:
    template_name = random_suffix_name('simple-template', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['TEMPLATE_NAME'] = template_name

    resource_data = load_ses_resource(
        'template_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, TEMPLATE_RESOURCE_PLURAL,
        template_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=10)

    assert cr is not None
    assert cr['status'] is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(get_template, ses_client, template_name))


@service_marker
@pytest.mark.canary
class TestTemplate:
    def test_create_update_template(self, simple_template, ses_client):
        (ref, cr) = simple_template
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )
        assert get_template(ses_client, cr['spec']['name']) is not None

        updates = {
            'spec': {
                'textPart': 'updated-text',
                'subjectPart': 'updated-subject',
                'htmlPart': '<h5>Updated HTML</h5>',
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

        template = get_template(ses_client, cr['spec']['name'])
        updated_template = updates['spec']
        assert template['Template'] == {
            'TemplateName': cr['spec']['name'],
            'TextPart': updated_template['textPart'],
            'SubjectPart': updated_template['subjectPart'],
            'HtmlPart': updated_template['htmlPart'],
        }

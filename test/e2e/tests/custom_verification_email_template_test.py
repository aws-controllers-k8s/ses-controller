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

"""Integration tests for the CustomVerificationEmailTemplate resource
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

from .template_test import simple_template


CUSTOM_VERIFICATION_EMAIL_TEMPLATE_RESOURCE_PLURAL = "customverificationemailtemplates"

@pytest.fixture(scope='module')
def ses_client():
    return boto3.client(SERVICE_NAME)

def get_custom_verification_email_template(ses_client, template_name):
    try:
        return ses_client.get_custom_verification_email_template(TemplateName=template_name)
    except ses_client.exceptions.CustomVerificationEmailTemplateDoesNotExistException:
        return None

@pytest.fixture
def simple_custom_verification_email_template(simple_template, ses_client) -> Tuple[k8s.CustomResourceReference, Dict]:
    custom_verification_email_template_name = random_suffix_name('simple-verification-template', 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements['CUSTOM_VERIFICATION_EMAIL_TEMPLATE_NAME'] = custom_verification_email_template_name

    (ref, template_cr) = simple_template
    replacements['TEMPLATE_REF_NAME'] = template_cr['spec']['name']

    resource_data = load_ses_resource(
        'custom_verification_email_template_simple',
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, CUSTOM_VERIFICATION_EMAIL_TEMPLATE_RESOURCE_PLURAL,
        custom_verification_email_template_name, namespace='default',
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref, wait_periods=10)

    assert cr is not None
    assert cr['status'] is not None
    assert k8s.get_resource_exists(ref)

    yield ref, cr

    _, deleted = k8s.delete_custom_resource(ref, 3, 10)
    assert deleted
    wait_until_deleted(partial(get_custom_verification_email_template, ses_client, custom_verification_email_template_name))


@service_marker
@pytest.mark.canary
class TestCustomVerificationEmailTemplate:
    def test_create_update_custom_verification_email_template(self, simple_custom_verification_email_template, ses_client):
        (ref, cr) = simple_custom_verification_email_template
        assert k8s.wait_on_condition(
            ref,
            condition.CONDITION_TYPE_RESOURCE_SYNCED,
            'True',
            wait_periods=MAX_WAIT_FOR_SYNCED_MINUTES,
        )
        template_name = cr['spec']['templateRef']['from']['name']
        get_current_template = partial(get_custom_verification_email_template, ses_client, template_name)
        assert get_current_template() is not None

        updates = {
            'spec': {
                'templateSubject': 'updated-subject',
                'failureRedirectionURL': 'https://example.com/updated',
                'successRedirectionURL': 'https://example.com/updated',
                'templateContent': 'updated template content',
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

        template = get_current_template()
        template.pop('ResponseMetadata', None)
        updated_template = updates['spec']
        assert template == {
            'TemplateName': template_name,
            'TemplateSubject': updated_template['templateSubject'],
            'FromEmailAddress': cr['spec']['fromEmailAddress'],
            'FailureRedirectionURL': updated_template['failureRedirectionURL'],
            'SuccessRedirectionURL': updated_template['successRedirectionURL'],
            'TemplateContent': updated_template['templateContent'],
        }

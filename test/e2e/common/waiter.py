"""Utilities for working with SES resources"""

import datetime
import time
import typing

import pytest

DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS = 30
DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS = 15
MAX_WAIT_FOR_SYNCED_MINUTES = 1

GetResourceFunc = typing.NewType(
    'GetResourceFunc',
    typing.Callable[[], dict],
)


def wait_until_deleted(
        get_resource: GetResourceFunc,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_INTERVAL_SECONDS,
) -> None:
    """Waits until a resource is deleted from the SES API

    Usage:
        from e2e.common.waiter import wait_until_deleted

        wait_until_deleted(partial(ses_client.describe_configuration_set, **resource_query))

    Raises:
        pytest.fail upon timeout
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while get_resource() is not None:
        if datetime.datetime.now() >= timeout:
            pytest.fail('Timed out waiting for resource to be deleted in SES')
        time.sleep(interval_seconds)

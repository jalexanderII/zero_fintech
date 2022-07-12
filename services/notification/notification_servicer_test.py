import os

import pytest
from dotenv import load_dotenv
from pytest_mock import MockerFixture
from twilio.rest import Client

from gen.Python.core.users_pb2 import User
from gen.Python.notification.notification_pb2 import NotificationRequest, NOTIFICATION_DESTINATION_SMS
from services.notification.notification_servicer import NotificationService


load_dotenv()

TEST_USER_ID = 'test_id'


@pytest.fixture
def mock_notification_server() -> NotificationService:
    account_sid = os.getenv('TWILIO_TEST_ACCOUNT_SID')
    auth_token = os.getenv('TWILIO_TEST_AUTH_TOKEN')
    client = Client(account_sid, auth_token)
    sending_number = os.getenv('TWILIO_TEST_PHONE_NUMBER')
    return NotificationService(client=client, sending_number=sending_number)


@pytest.fixture
def patch__fetch_user(mocker: MockerFixture) -> None:
    mocker.patch.object(NotificationService, '_fetch_user').return_value = User(
        id=TEST_USER_ID,
        username='test_username',
        email='test.user@test.mail',
        password='test_pwd',
        phone_number='+19174427389'
    )


def test_notify(patch__fetch_user: None, mock_notification_server: NotificationService):
    notification_response = mock_notification_server.Notify(request=NotificationRequest(
        user_id=TEST_USER_ID,
        destination=NOTIFICATION_DESTINATION_SMS),
        message='Test message'
    )

    assert notification_response.successful, f"Error message was: {notification_response.error_message}"

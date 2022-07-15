import os

import pytest
from dotenv import load_dotenv
from twilio.rest import Client

from gen.Python.notification.notification_pb2 import SendSMSRequest
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


def test_notify(mock_notification_server: NotificationService):
    notification_response = mock_notification_server.SendSMS(request=SendSMSRequest(
        phone_number="+19174427389",
        message='Test message'
    ))

    assert notification_response.successful, f"Error message was: {notification_response.error_message}"

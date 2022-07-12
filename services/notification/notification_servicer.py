from attr import define
from twilio.rest import Client

from gen.Python.notification.notification_pb2 import SendSMSResponse, SendSMSRequest
from gen.Python.notification.notification_pb2_grpc import NotificationServicer


@define
class NotificationService(NotificationServicer):
    client: Client
    sending_number: str

    def SendSMS(self, request: SendSMSRequest, context=None) -> SendSMSResponse:
        message = self.client.messages.create(
            body=request.message,
            from_=self.sending_number,
            to=request.phone_number
        )
        return SendSMSResponse(successful=message.error_code is None, error_message=str(message.error_message))

import os

import grpc
from attr import define
from twilio.rest import Client

from gen.Python.core.core_pb2_grpc import CoreStub
from gen.Python.core.users_pb2 import GetUserRequest, User
from gen.Python.notification.notification_pb2 import NotificationRequest, NotificationResponse, \
    NOTIFICATION_DESTINATION_SMS
from gen.Python.notification.notification_pb2_grpc import NotificationServicer


@define
class NotificationService(NotificationServicer):
    client: Client
    sending_number: str
    core_client: CoreStub = CoreStub(grpc.insecure_channel(f'localhost:{os.getenv("CORE_SERVER_PORT")}'))

    def Notify(self, request: NotificationRequest, context) -> NotificationResponse:
        user: User = self._fetch_user(user_id=request.user_id)

        if request.destination == NOTIFICATION_DESTINATION_SMS:
            message = self.client.messages.create(
                body=request.message,
                from_=self.sending_number,
                to=user.phone_number
            )
            return NotificationResponse(successful=message.error_code is None, error_message=str(message.error_message))
        else:
            return NotificationResponse(successful=False, error_message=f"Unknown destination: {request.destination}")

    def _fetch_user(self, user_id: str) -> User:
        return self.core_client.GetUser(GetUserRequest(id=user_id))

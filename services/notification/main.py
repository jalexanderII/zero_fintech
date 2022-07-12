import logging
import os
import sys
from concurrent import futures

import grpc
from dotenv import load_dotenv
from twilio.rest import Client

from gen.Python.notification import notification_pb2_grpc as NotificationServicePB
from services.notification.notification_servicer import NotificationService

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger("Server")


def serve():
    logger.info("Load .env file")
    load_dotenv()

    logger.info("Initiate Twilio client and servicer")
    account_sid = os.getenv('TWILIO_ACCOUNT_SID')
    auth_token = os.getenv('TWILIO_AUTH_TOKEN')
    sending_number = os.getenv('TWILIO_PHONE_NUMBER')
    client = Client(account_sid, auth_token)
    servicer = NotificationService(client=client, sending_number=sending_number)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    NotificationServicePB.add_NotificationServicer_to_server(servicer, server)

    server.add_insecure_port("[::]:{}".format(os.getenv("NOTIFICATION_SERVER_PORT")))
    server.start()
    logger.info("Server running")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()






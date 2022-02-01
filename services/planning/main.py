import logging
import os
import sys
from concurrent import futures

import grpc
from dotenv import load_dotenv

from database.db import initiate_mongo_client
from gen.Python.core import core_pb2_grpc as coreClient, accounts_pb2 as Accounts
from gen.Python.planning import planning_pb2_grpc as PlanningServicePB
from server.planning_servicer import PlanningService

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger("Server")


def serve():
    logger.info("Load .env file")
    load_dotenv()

    logger.info("Initiate Mongo client and servicer")
    servicer = PlanningService(planning_collection=initiate_mongo_client())
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    PlanningServicePB.add_PlanningServicer_to_server(servicer, server)

    server.add_insecure_port("[::]:{}".format(os.getenv("PLANNING_SERVER_PORT")))
    server.start()
    logger.info("Server running")
    server.wait_for_termination()


def run_client():
    with grpc.insecure_channel(f'localhost:{os.getenv("CORE_SERVER_PORT")}') as channel:
        core_client_stub = coreClient.CoreStub(channel=channel)
        response = core_client_stub.ListAccounts(request=Accounts.ListAccountRequest())
        print(response)


if __name__ == "__main__":
    serve()
    # run_client()

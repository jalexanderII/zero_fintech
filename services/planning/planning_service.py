import logging
import os
import sys
from concurrent import futures

print('\n'.join(sys.path))

import grpc
from dotenv import load_dotenv
from pymongo.collection import Collection

from database.database import initiate_mongo_client
from gen.Python.core import core_pb2_grpc as coreClient, accounts_pb2 as Accounts
from gen.Python.planning import planning_pb2_grpc as PlanningServicePB
from server.server import PlanningServicer

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)]
)
logger = logging.getLogger('Server')


def serve():
    logger.info('Load .env file')
    load_dotenv()


    logger.info('Initiate Mongo client and servicer')
    planning_collection: Collection = initiate_mongo_client()

    servicer = PlanningServicer(planning_collection)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    PlanningServicePB.add_PlanningServicer_to_server(servicer, server)

    server.add_insecure_port("[::]:{}".format(os.getenv('PLANNING_SERVER_PORT')))
    logger.info('Server running')
    server.start()
    server.wait_for_termination()




def run_client():
    with grpc.insecure_channel('localhost:9090') as channel:
        core_client_stub = coreClient.CoreStub(channel=channel)
        response = core_client_stub.ListAccounts(request=Accounts.ListAccountRequest())
        print(response)


if __name__ == '__main__':
    serve()
    # run_client()

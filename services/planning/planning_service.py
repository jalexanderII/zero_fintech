import logging
from concurrent import futures
import grpc
from dotenv import load_dotenv
import os


import sys
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, 'gen')))

# from Python.planning import planning_pb2_grpc
from Python.planning import planning_pb2_grpc as PlanningServicePB
from Python.core import core_pb2_grpc as coreClient, accounts_pb2 as Accounts

from server.server import PlanningServicer
from database.database import initateMongoClient

def serve():
    # load .env file
    logging.info('Load .env file')
    load_dotenv()
    # initiate Mongo client and servicer
    logging.info('Initiate Mongo client and servicer')
    mongoClient = initateMongoClient()
    planningCollection = mongoClient[os.getenv('PLANNING_COLLECTION')]
    servicer = PlanningServicer(planningCollection=planningCollection)
    # start server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    PlanningServicePB.add_PlanningServicer_to_server(servicer, server)
    server.add_insecure_port('[::]:50051')
    server.start()
    logging.info('Server running')
    server.wait_for_termination()

def run_client():
    with grpc.insecure_channel('localhost:9090') as channel:
        coreClientStub = coreClient.CoreStub(channel=channel)
        response = coreClientStub.ListAccounts(request=Accounts.ListAccountRequest())
        print(response)

if __name__ == '__main__':
    logging.basicConfig()
    # serve()
    run_client()
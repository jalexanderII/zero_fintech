import logging
from concurrent import futures
import grpc
from dotenv import load_dotenv
import os


import sys
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, 'gen')))

from Python.planning import planning_pb2_grpc
from Python.core import core_pb2_grpc, accounts_pb2

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
    planning_pb2_grpc.add_PlanningServicer_to_server(servicer, server)
    server.add_insecure_port('[::]:50051')
    server.start()
    logging.info('Server running')
    server.wait_for_termination()

def run_client():
    with grpc.insecure_channel('localhost:9090') as channel:
        stub = core_pb2_grpc.CoreStub(channel=channel)
        response = stub.ListAccounts(request=accounts_pb2.ListAccountRequest())
        print(response)

if __name__ == '__main__':
    logging.basicConfig()
    serve()
    # run_client()
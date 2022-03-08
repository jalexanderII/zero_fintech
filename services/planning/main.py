import logging
import os
import sys
from concurrent import futures

import grpc
from dotenv import load_dotenv
from fastapi import FastAPI

from database.db import initiate_mongo_client
from gen.Python.planning import planning_pb2_grpc as PlanningServicePB
from server.planning_servicer import PlanningService
from app.api import dashboard

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


def create_application() -> FastAPI:
    application = FastAPI()
    application.include_router(dashboard.router, prefix='/dashboard')

    return application


app = create_application()


@app.on_event("startup")
async def startup_event():
    logger.info("Starting up...")   # TODO: or use different logger for app?
    # somehow init db


@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Shutting down...")


if __name__ == "__main__":
    serve()

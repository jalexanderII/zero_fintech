import os

from dotenv import load_dotenv

from server import PlanningServer
from database.database import initateMongoClient

def gen_server() -> PlanningServer:
    # load .env file
    load_dotenv()

    mongoClient = initateMongoClient()
    planningCollection = mongoClient[os.getenv('PLANNING_COLLECTION')]

    server = PlanningServer(planningCollection=planningCollection)

    return server

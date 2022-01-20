import os
from mongoengine import connect
import logging
from dotenv import load_dotenv


def initateMongoClient():
    load_dotenv()

    mongoClient = None
    try:
        mongoClient = connect(host=os.getenv('MONGOURIPY'), db=os.getenv('PLANNING_DB_NAME'))
    except:
        logging.error("Connecting to MongoDB failed", exc_info=True)

    return mongoClient
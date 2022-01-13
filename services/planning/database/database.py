import os
from mongoengine import connect
from dotenv import load_dotenv
import logging

def initateMongoClient():
    # load .env file
    load_dotenv()

    mongoClient = None
    try:
        mongoClient = connect(host=os.getenv('MONGOURI'), db=os.getenv('CORE_DB_NAME'))
    except:
        logging.error("Connecting to MongoDB failed", exc_info=True)

    return mongoClient
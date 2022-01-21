import os

from dotenv import load_dotenv
from pymongo import MongoClient
from pymongo.collection import Collection
from pymongo.database import Database


def initiate_mongo_client() -> Collection:
    load_dotenv()
    client = MongoClient(os.getenv('MONGOURI'))
    db: Database = client[os.getenv('PLANNING_DB_NAME')]
    return db[os.getenv('PLANNING_COLLECTION')]

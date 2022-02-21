import os

from dotenv import load_dotenv
from pymongo import MongoClient
from pymongo.collection import Collection


def initiate_mongo_client() -> Collection:
    load_dotenv()
    client = MongoClient(os.getenv("MONGOURI"))
    db = client[os.getenv("PLANNING_DB_NAME")]
    return db[os.getenv("PLANNING_COLLECTION")]


def initiate_mongo_test_client() -> Collection:
    load_dotenv()
    client = MongoClient(os.getenv("MONGOURI"))
    db = client[os.getenv("PLANNING_DB_NAME")]
    return db[f'{os.getenv("PLANNING_COLLECTION")}_TEST']

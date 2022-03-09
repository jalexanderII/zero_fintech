import os
from typing import Generator, Any

import grpc
from dotenv import load_dotenv
from motor.motor_asyncio import AsyncIOMotorCollection, AsyncIOMotorClient

from gen.Python.core.core_pb2_grpc import CoreStub


async def get_async_mongo_collection() -> Generator[AsyncIOMotorCollection, Any, None]:
    try:
        load_dotenv()
        client = AsyncIOMotorClient(os.getenv("MONGOURI"))
        db = client[os.getenv("PLANNING_DB_NAME")]
        yield db[os.getenv("PLANNING_COLLECTION")]
    finally:
        client.close()


def get_core_client() -> CoreStub:
    return CoreStub(
        grpc.insecure_channel(f'localhost:{os.getenv("CORE_SERVER_PORT")}')
    )

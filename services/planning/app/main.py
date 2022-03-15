import logging

from fastapi import FastAPI

from services.planning.app.api import dashboard


logger = logging.getLogger("uvicorn")


def create_application() -> FastAPI:
    application = FastAPI()
    application.include_router(dashboard.router, prefix='/dashboard')

    return application


app = create_application()


@app.on_event("startup")
async def startup_event():
    logger.info("Starting up...")
    # somehow init db


@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Shutting down...")

import datetime
import logging
import os
from typing import List

import grpc
from dotenv import load_dotenv

from gen.Python.common.payment_plan_pb2 import ListPaymentPlanRequest, PaymentPlan
from gen.Python.core.core_pb2_grpc import CoreStub
from gen.Python.planning.planning_pb2_grpc import PlanningStub

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)


load_dotenv()


core_client: CoreStub = CoreStub(grpc.insecure_channel(f'localhost:{os.getenv("CORE_SERVER_PORT")}'))
planning_client: PlanningStub = PlanningStub(grpc.insecure_channel(f'localhost:{os.getenv("PLANNING_SERVER_PORT")}'))


def run(event, context):
    current_time = datetime.datetime.now().time()
    name = context.function_name
    logger.info("Your cron function " + name + " ran at " + str(current_time))

    # TODO: write this in golang in core service

    # TODO: list active payment plans
    payment_plans: List[PaymentPlan] = planning_client.ListPaymentPlans(ListPaymentPlanRequest()).payment_plans

    # TODO: see if there is any upcoming in the next few days

    # TODO: check if users have enough in their balance to pay them

    # TODO: send notification to appropriate user

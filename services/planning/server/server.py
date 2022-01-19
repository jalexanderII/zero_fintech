from collections import defaultdict
from typing import OrderedDict, List
import numpy as np
import pandas as pd
import grpc
import json

import logging
from datetime import timedelta, datetime
from dotenv import load_dotenv

from google.protobuf.timestamp_pb2 import Timestamp
from pymongo.database import Database



import sys, os
# make gen/Python importable by import Python.X
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, os.pardir, 'gen')))
# from Python.core import accounts_pb2, core_pb2_grpc, payment_task_pb2, transactions_pb2
from Python.common.common_pb2 import (PlanType as PlanTypePB, PaymentFrequency as PaymentFrequencyPB, DELETE_STATUS as DELETE_STATUS_PB)
from Python.planning.planning_pb2_grpc import PlanningServicer
from Python.planning.planning_pb2 import CreatePaymentPlanRequest, CreatePaymentPlanResponse
from Python.planning.payment_plan_pb2 import (PaymentPlan as PaymentPlanPB, PaymentStatus as PaymentStatusPB, PaymentActionStatus as PaymentActionStatusPB, PaymentAction as PaymentActionPB, DeletePaymentPlanResponse, ListPaymentPlanResponse, UpdatePaymentPlanRequest, GetPaymentPlanRequest, ListPaymentPlanRequest, DeletePaymentPlanRequest)

# make ../database importable
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir)))
from database.models import (PaymentPlan as PaymentPlanDB, PaymentFrequency as PaymentFrequencyDB, PlanType as PlanTypeDB, PaymentStatus as PaymentStatusDB, PaymentAction as PaymentActionDB, PaymentActionStatus as PaymentActionStatusDB)

from utils import paymentPlanDBToPB, paymentPlanPBToDB, payment_actions_pb_to_db
from PaymentPlanBuilder import createPaymentPlan


class PlanningServicer(PlanningServicer):

    def __init__(self, planningCollection: Database) -> None:
        self.planningCollection = planningCollection

    def CreatePaymentPlan(self, request: CreatePaymentPlanRequest, context: grpc.ServicerContext) -> CreatePaymentPlanResponse:
        """ Creates PaymentPlan(s) for given request containing PaymentTasks.

        Args:
            request (CreatePaymentPlanRequest): Protobuf containing list of PaymentTasks
            context (grpc.ServicerContext): [description]
        Returns:
            CreatePaymentPlanResponse: List of possible PaymentPlans for given PaymentTasks
        """
        paymentPlanListPB = createPaymentPlan(list(request.payment_tasks))
        for paymentPlanPB in paymentPlanListPB:
            self.CreatePaymentPlanNoBuilding(paymentPlanPB)
        return CreatePaymentPlanRequest(payment_plans=paymentPlanListPB)

    def CreatePaymentPlanNoBuilding(self, paymentPlanPB: PaymentPlanPB) -> PaymentPlanPB:
        """ Saves a PaymentPlan into the database without creating it.

        Args:
            paymentPlanPB (PaymentPlanPB): PaymentPlan which shall be stored in database
        Returns:
            PaymentPlanPB: A protobuf of the saved PaymentPlan with payment_plan_id=id in database
        """
        paymentPlanDB = paymentPlanPBToDB(paymentPlanPB)
        paymentPlanDB.save()
        paymentPlanDB.PaymentPlanID = str(paymentPlanDB.id)
        paymentPlanDB.save()
        return paymentPlanDBToPB(paymentPlanDB)
    
    def GetPaymentPlan(self, request: GetPaymentPlanRequest, context) -> PaymentPlanPB:
        paymentPlanDB = PaymentPlanDB.objects.get(id=request.payment_plan_id)
        return paymentPlanDBToPB(paymentPlanDB)

    def ListPaymentPlans(self, request: ListPaymentPlanRequest, context) -> ListPaymentPlanResponse:
        paymentPlansDB = PaymentPlanDB.objects
        paymentPlansPB = []
        for pp in paymentPlansDB:
            paymentPlansPB.append(paymentPlanDBToPB(pp))
        return ListPaymentPlanResponse(payment_plans=paymentPlansPB)

    def UpdatePaymentPlan(self, request: UpdatePaymentPlanRequest, context) -> PaymentPlanPB:
        paymentPlanPB = request.payment_plan
        PaymentPlanDB.objects.get(id=paymentPlanPB.payment_plan_id).update(
            UserID=paymentPlanPB.user_id,
            PaymentTaskID=list(paymentPlanPB.payment_task_id), 
            Timeline=paymentPlanPB.timeline,
            PaymentFrequency=paymentPlanPB.payment_freq,
            AmountPerPayment=paymentPlanPB.amount_per_payment,
            PlanType=paymentPlanPB.plan_type,
            EndDate=paymentPlanPB.end_date.ToDatetime(),
            Active=paymentPlanPB.active,
            Status=paymentPlanPB.status,
            PaymentAction=payment_actions_pb_to_db(paymentPlanPB.payment_action)
        )
        return paymentPlanPB

    def DeletePaymentPlan(self, request: DeletePaymentPlanRequest, context) -> DeletePaymentPlanResponse:
        paymentPlanDB = PaymentPlanDB.objects.get(id=request.payment_plan_id)
        paymentPlanPB = paymentPlanDBToPB(paymentPlanDB)
        paymentPlanDB.delete()
        return DeletePaymentPlanResponse(status=DELETE_STATUS_PB.DELETE_STATUS_SUCCESS, payment_plan=paymentPlanPB)
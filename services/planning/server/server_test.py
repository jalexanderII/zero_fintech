from datetime import datetime
import os

from dotenv import load_dotenv

import sys
import json

from server import PlanningServicer
import sys, os
# make gen/Python importable by import Python.X
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, os.pardir, 'gen')))
from Python.planning.payment_plan_pb2 import DeletePaymentPlanRequest, GetPaymentPlanRequest, UpdatePaymentPlanRequest

# make ../database importable
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir)))
from database.database import initateMongoClient
from database.models import (PaymentPlan as PaymentPlanDB, PaymentFrequency as PaymentFrequencyDB, PlanType as PlanTypeDB,
    PaymentStatus as PaymentStatusDB, PaymentAction as PaymentActionDB, PaymentActionStatus as PaymentActionStatusDB)

from utils import paymentPlanDBToPB

def gen_server() -> PlanningServicer:
    # load .env file
    load_dotenv()

    mongoClient = initateMongoClient()
    planningCollection = mongoClient[os.getenv('PLANNING_COLLECTION')]

    server = PlanningServicer(planningCollection=planningCollection)

    return server

EXAMPLE_PAYMENT_PLANS = [
        PaymentPlanDB(UserID="afdf93f02c61d1be8eaf613",
            PaymentTaskID=['61dfa8296c734067e6726761'],
            Timeline=5.0,
            PaymentFrequency=PaymentFrequencyDB.PAYMENT_FREQUENCY_BIWEEKLY,
            AmountPerPayment=23.0,
            PlanType=PlanTypeDB.PLAN_TYPE_OPTIM_CREDIT_SCORE,
            EndDate=datetime.now(),
            Active=True,
            Status=PaymentStatusDB.PAYMENT_STATUS_IN_DEFAULT,
            PaymentAction=[
                PaymentActionDB(AccountID='ef159b621d3c2b15a6e53ec9',
                    Amount=23.0,
                    TransactionDate=datetime.now(),
                    PaymentActionStatus=PaymentActionStatusDB.PAYMENT_ACTION_STATUS_IN_DEFAULT)
        ]),
        PaymentPlanDB(PaymentPlanID="61e82eadbb7057f26763b443",
            UserID="61df93c0ac601d1be8e6af28",
            PaymentTaskID=['61dfa8296c734067e6726761', 'a2ffa82f6c734067e6726761'],
            Timeline=1.0,
            PaymentFrequency=PaymentFrequencyDB.PAYMENT_FREQUENCY_WEEKLY,
            AmountPerPayment=100.0,
            PlanType=PlanTypeDB.PLAN_TYPE_MIN_FEES,
            EndDate=datetime.now(),
            Active=True,
            Status=PaymentStatusDB.PAYMENT_STATUS_CURRENT,
            PaymentAction=[
                PaymentActionDB(AccountID='61df9b621d2c2b15a6e53ec9',
                    Amount=100.0,
                    TransactionDate=datetime.now(),
                    PaymentActionStatus=PaymentActionStatusDB.PAYMENT_ACTION_STATUS_PENDING)
        ]),
]

def test_list_payment_plans():
    server = gen_server()
    payment_plans = server.ListPaymentPlans(None, None).payment_plans
    assert len(payment_plans) > 0, f"Server did not return any PaymentPlan"

def test_create_no_building_and_delete_payment_plan():
    server = gen_server()
    # original number of PaymentPlan
    originalPaymentPlansLen = len(server.ListPaymentPlans(None, None).payment_plans)
    # adding one PaymentPlan
    paymentPlanPB = paymentPlanDBToPB(EXAMPLE_PAYMENT_PLANS[0])
    paymentPlanPB = server._createPaymentPlan(paymentPlanPB)
    updatedPaymentPlansLen = len(server.ListPaymentPlans(None, None).payment_plans)
    assert originalPaymentPlansLen == updatedPaymentPlansLen - 1, f"Before insertion number of PaymentPlans={originalPaymentPlansLen}; after insertion number of PaymentPlans={updatedPaymentPlansLen}"
    # deleting the added paymentPlan
    deleteRequest = DeletePaymentPlanRequest(payment_plan_id=paymentPlanPB.payment_plan_id)
    server.DeletePaymentPlan(request=deleteRequest, context=None)
    updatedPaymentPlansLen = len(server.ListPaymentPlans(None, None).payment_plans)
    assert originalPaymentPlansLen == updatedPaymentPlansLen, f"Before insertion number of PaymentPlans={originalPaymentPlansLen}; after insertion & deletion number of PaymentPlans={updatedPaymentPlansLen}"

def test_get_payment_plan():
    server = gen_server()
    payment_plan_id = '61e82eedbb7057f26763b444'
    paymentPlanGet = server.GetPaymentPlan(GetPaymentPlanRequest(payment_plan_id=payment_plan_id), context=None)
    assert payment_plan_id == paymentPlanGet.payment_plan_id, f"Plan inserted ID: {payment_plan_id}\nPlan retrieved ID: {paymentPlanGet.payment_plan_id}"

def test_update_payment_plan():
    server = gen_server()

    payment_plan_id = '61e82eadbb7057f26763b443'
    # check if the PaymentPlan has one PaymentTask
    paymentPlanPBOrig = server.GetPaymentPlan(GetPaymentPlanRequest(payment_plan_id=payment_plan_id), context=None)
    paymentTasks = paymentPlanPBOrig.payment_task_id
    assert len(paymentTasks) == 1, f"PaymentPlan has {len(paymentTasks)} PaymentTasks instead of 1"
    # update with the second element of EXAMPLE_PAYMENT_PLANS which has 2 PaymentTasks
    server.UpdatePaymentPlan(UpdatePaymentPlanRequest(
        payment_plan_id=payment_plan_id, payment_plan=paymentPlanDBToPB(EXAMPLE_PAYMENT_PLANS[1])
    ), context=None)
    # check if the PaymentPlan has two PaymentTasks (as the second entry of EXAMPLE_PAYMENT_PLANS has)
    paymentTasks = server.GetPaymentPlan(GetPaymentPlanRequest(payment_plan_id=payment_plan_id), context=None).payment_task_id
    assert len(paymentTasks) == 2, f"PaymentPlan has {len(paymentTasks)} PaymentTask(s) instead of 2"
    # update again to original PaymentPlan and check that we only have 1 PaymentTask
    server.UpdatePaymentPlan(UpdatePaymentPlanRequest(payment_plan_id=payment_plan_id, payment_plan=paymentPlanPBOrig),
                            context=None)
    paymentTasks = server.GetPaymentPlan(GetPaymentPlanRequest(payment_plan_id=payment_plan_id), context=None).payment_task_id
    assert len(paymentTasks) == 1, f"PaymentPlan has {len(paymentTasks)} PaymentTasks instead of 1"

def test_create_payment_plan():
    print("test_create_payment_plan wasn't implemented yet")
    pass

if __name__ == '__main__':
    test_list_payment_plans
    test_get_payment_plan()
    test_create_no_building_and_delete_payment_plan()
    test_update_payment_plan()
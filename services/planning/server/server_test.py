from datetime import datetime
import os

from dotenv import load_dotenv

import sys
import json

sys.path.append('/Users/joschkabraun/dev/zero_fintech/services/planning')

from server import PlanningServicer
from database.database import initateMongoClient
from database import models as db_models
from gen.planning import payment_plan_pb2

def gen_server() -> PlanningServicer:
    # load .env file
    load_dotenv()

    mongoClient = initateMongoClient()
    planningCollection = mongoClient[os.getenv('PLANNING_COLLECTION')]

    server = PlanningServicer(planningCollection=planningCollection)

    return server

EXAMPLE_PAYMENT_PLANS = [
    db_models.PaymentPlan(PaymentPlanID="61dfa3c6ac621d1be8e64613",
        UserID="61df93c0ac601d1be8e64613",
        PaymentTaskID=['61dfa8296c734067e6726761'],
        Timeline=1.0,
        PaymentFrequency=db_models.PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY,
        AmountPerPayment=100.0,
        PlanType=db_models.PlanType.PLAN_TYPE_MIN_FEES,
        EndDate=datetime.now(),
        Active=True,
        Status=db_models.PaymentStatus.PAYMENT_STATUS_CURRENT,
        PaymentAction=[
            db_models.PaymentAction(AccountID='61df9b621d2c2b15a6e53ec9',
                Amount=100.0,
                TransactionDate=datetime.now(),
                PaymentActionStatus=db_models.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING)
        ]),
        db_models.PaymentPlan(PaymentPlanID="a2efa1c6ac621d1bf8e64613",
        UserID="afdf93f02c61d1be8eaf613",
        PaymentTaskID=['61dfa8296c734067e6726761'],
        Timeline=1.0,
        PaymentFrequency=db_models.PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY,
        AmountPerPayment=100.0,
        PlanType=db_models.PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE,
        EndDate=datetime.now(),
        Active=True,
        Status=db_models.PaymentStatus.PAYMENT_STATUS_CURRENT,
        PaymentAction=[
            db_models.PaymentAction(AccountID='ef159b621d3c2b15a6e53ec9',
                Amount=100.0,
                TransactionDate=datetime.now(),
                PaymentActionStatus=db_models.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING)
        ]),
        db_models.PaymentPlan(PaymentPlanID="61dfa3c6ac621d1be8e64613",
        UserID="61df93c0ac601d1be8e6af28",
        PaymentTaskID=['61dfa8296c734067e6726761', 'a2ffa82f6c734067e6726761'],
        Timeline=1.0,
        PaymentFrequency=db_models.PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY,
        AmountPerPayment=100.0,
        PlanType=db_models.PlanType.PLAN_TYPE_MIN_FEES,
        EndDate=datetime.now(),
        Active=True,
        Status=db_models.PaymentStatus.PAYMENT_STATUS_CURRENT,
        PaymentAction=[
            db_models.PaymentAction(AccountID='61df9b621d2c2b15a6e53ec9',
                Amount=100.0,
                TransactionDate=datetime.now(),
                PaymentActionStatus=db_models.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING)
        ]),
]

def test_delete_payment_plan():
    server = gen_server()
    # original number of PaymentPlan
    originalPaymentPlansLen = len(server.ListPaymentPlans(None, None).payment_plans)
    # adding one PaymentPlan
    EXAMPLE_PAYMENT_PLANS[0].save()
    payment_plan_id = EXAMPLE_PAYMENT_PLANS[0].PaymentPlanID
    updatedPaymentPlansLen = len(server.ListPaymentPlans(None, None).payment_plans)
    assert originalPaymentPlansLen == updatedPaymentPlansLen - 1, f"Before insertion number of PaymentPlans={originalPaymentPlansLen}; after insertion number of PaymentPlans={updatedPaymentPlansLen}"
    # deleting the added paymentPlan
    deleteRequest = payment_plan_pb2.DeletePaymentPlanRequest(payment_plan_id=payment_plan_id)
    server.DeletePaymentPlan(request=deleteRequest, context=None)
    updatedPaymentPlansLen = len(server.ListPaymentPlans(None, None).payment_plans)
    assert originalPaymentPlansLen == updatedPaymentPlansLen, f"Before insertion number of PaymentPlans={originalPaymentPlansLen}; after insertion & deletion number of PaymentPlans={updatedPaymentPlansLen}"

def test_get_payment_plan():
    server = gen_server()
    # inserting a plan
    paymentPlanOrig = EXAMPLE_PAYMENT_PLANS[1]
    paymentPlanOrig.save()
    payment_plan_id = EXAMPLE_PAYMENT_PLANS[1].PaymentPlanID
    paymentPlanGet = server.GetPaymentPlan(payment_plan_pb2.GetPaymentPlanRequest(payment_plan_id=payment_plan_id), context=None)
    assert paymentPlanOrig.PaymentPlanID == paymentPlanGet.payment_plan_id, f"Plan inserted ID: {paymentPlanOrig.PaymentPlanID}\nPlan retrieved ID: {paymentPlanGet.payment_plan_id}"
    paymentPlanOrig.delete()

def test_update_payment_plan():
    server = gen_server()

    # # save first PaymentPlan
    paymentPlanOrig = EXAMPLE_PAYMENT_PLANS[0]
    paymentPlanOrig.save()
    assert len(server.GetPaymentPlan(payment_plan_pb2.GetPaymentPlanRequest(payment_plan_id=paymentPlanOrig.PaymentPlanID), context=None).payment_task_id) == 1
    server.UpdatePaymentPlan(payment_plan_pb2.UpdatePaymentPlanRequest(
        payment_plan_id=paymentPlanOrig.PaymentPlanID, payment_plan=server.paymentPlanDBToPB(EXAMPLE_PAYMENT_PLANS[2])
    ), context=None)
    assert len(server.GetPaymentPlan(payment_plan_pb2.GetPaymentPlanRequest(payment_plan_id=paymentPlanOrig.PaymentPlanID), context=None).payment_task_id) == 2, f"length = {len(server.GetPaymentPlan(payment_plan_pb2.GetPaymentPlanRequest(payment_plan_id=paymentPlanOrig.PaymentPlanID), context=None).payment_task_id)}"

if __name__ == '__main__':
    # test_delete_payment_plan()
    # test_get_payment_plan()
    test_update_payment_plan()
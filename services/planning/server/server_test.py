from datetime import datetime
import os

from dotenv import load_dotenv

import sys

# sys.path.append('/Users/joschkabraun/dev/zero_fintech/services/planning/database')
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
        EndDate=datetime.utcnow,
        Active=True,
        Status=db_models.PaymentStatus.PAYMENT_STATUS_CURRENT,
        PaymentAction=[
            db_models.PaymentAction(AccountID='61df9b621d2c2b15a6e53ec9',
                Amount=100.0,
                TransactionDate=datetime.utcnow,
                PaymentActionStatus=db_models.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING)
        ])
]

def test_delete_payment_plan():
    server = gen_server()

    # originalPaymentPlans = server.ListPaymentPlans(None, None)
    # originalPaymentPlansLen = len(originalPaymentPlans)
    # print(f"Number of payment plans before insertion: {originalPaymentPlansLen}")

    server.planningCollection.insert(EXAMPLE_PAYMENT_PLANS[0])
    updatedPaymentPlans = server.ListPaymentPlans(None, None)
    updatedPaymentPlansLen = len(updatedPaymentPlans)
    print(f"Number of payment plans after insertion: {updatedPaymentPlansLen}")

    deleteRequest = payment_plan_pb2.DeletePaymentPlanRequest(payment_plan_id="61dfa3c6ac621d1be8e64613")
    server.DeletePaymentPlan(request=deleteRequest, context=None)
    updatedPaymentPlans = server.ListPaymentPlans(None, None)
    updatedPaymentPlansLen = len(updatedPaymentPlans)
    print(f"Number of payment plans after deletion: {updatedPaymentPlansLen}")

if __name__ == '__main__':
    test_delete_payment_plan()
from gen.planning import payment_plan_pb2
from database import models as db_models
from database.database import initateMongoClient
from server import PlanningServicer
import pytest
from datetime import datetime
import os

from dotenv import load_dotenv

import sys
import json

from services.planning.server.server_test import EXAMPLE_PAYMENT_PLANS, gen_server

sys.path.append('/Users/joschkabraun/dev/zero_fintech/services/planning')


@pytest.fixture
def mock_payment_plan():
    return [
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
                              PaymentTaskID=[
                                  '61dfa8296c734067e6726761', 'a2ffa82f6c734067e6726761'],
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


def test_get_payment_plan(
    mock_payment_plan
):
    server = gen_server()
    # inserting a plan
    paymentPlanOrig = mock_payment_plan[1]
    paymentPlanOrig.save()
    payment_plan_id = mock_payment_plan[1].PaymentPlanID
    paymentPlanGet = server.GetPaymentPlan(payment_plan_pb2.GetPaymentPlanRequest(
        payment_plan_id=payment_plan_id), context=None)
    assert paymentPlanOrig.PaymentPlanID == paymentPlanGet.payment_plan_id, f"Plan inserted ID: {paymentPlanOrig.PaymentPlanID}\nPlan retrieved ID: {paymentPlanGet.payment_plan_id}"
    paymentPlanOrig.delete()

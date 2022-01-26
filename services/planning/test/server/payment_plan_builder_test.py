import datetime

import grpc
import pytest
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PlanType, PaymentFrequency, PaymentActionStatus
from gen.Python.common.payment_task_pb2 import PaymentTask, MetaData
from gen.Python.core.core_pb2_grpc import CoreStub
from services.planning.server.payment_plan_builder import PaymentPlanBuilder, payment_plan_builder
from services.planning.server.utils import shift_date_by_payment_frequency


@pytest.fixture
def gen_payment_plan_builder() -> PaymentPlanBuilder:
    return payment_plan_builder

user_id = '61df93c0ac601d1be8e64613'
accName2Id = {'Amex': '61df9b621d2c2b15a6e53ec9', 'Chase': '61df9af7f18b94fc44d09fb9'}

def test_create_payment_plan_min_fees_2_month_monthly(gen_payment_plan_builder):
    paymentTasks = [
        PaymentTask(user_id=user_id, account_id=accName2Id['Amex'], amount=500), # Amex
        PaymentTask(user_id=user_id, account_id=accName2Id['Chase'], amount=500), # Chase
    ]
    metaData = MetaData(preferred_plan_type=PlanType.PLAN_TYPE_MIN_FEES, preferred_timeline_in_months=2.0,
        preferred_payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)

    paymentPlans = gen_payment_plan_builder.createPaymentPlan(paymentTasks=paymentTasks, metaData=metaData)

    print(paymentPlans)

    assert len(paymentPlans) == 1

    paymentPlan = paymentPlans[0]
    assert paymentPlan.plan_type == PlanType.PLAN_TYPE_MIN_FEES
    # assert paymentPlan

    assert paymentPlan.payment_action[0].account_id == accName2Id['Chase']
    assert paymentPlan.payment_action[0].amount == 500
    assert paymentPlan.payment_action[0].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(datetime.datetime.now(), PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB = Timestamp()
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[0].transaction_date == transactionDatePB

    assert paymentPlan.payment_action[1].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[1].amount == 500
    assert paymentPlan.payment_action[1].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[1].transaction_date == transactionDatePB


import datetime

import grpc
import pytest
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PlanType, PaymentFrequency, PaymentActionStatus, PaymentStatus
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

    assert len(paymentPlans) == 1

    paymentPlan = paymentPlans[0]
    assert paymentPlan.plan_type == PlanType.PLAN_TYPE_MIN_FEES
    assert paymentPlan.user_id == user_id
    assert paymentPlan.timeline == 2.0
    assert paymentPlan.payment_freq == PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY
    assert paymentPlan.amount_per_payment == 500
    assert paymentPlan.active == True
    assert paymentPlan.status == PaymentStatus.PAYMENT_STATUS_CURRENT
    assert len(paymentPlan.payment_action) == 2

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

    assert paymentPlan.end_date == transactionDatePB

def test_create_payment_plan_credit_score_2_month_monthly(gen_payment_plan_builder):
    paymentTasks = [
        PaymentTask(user_id=user_id, account_id=accName2Id['Amex'], amount=500), # Amex
        PaymentTask(user_id=user_id, account_id=accName2Id['Chase'], amount=500), # Chase
    ]
    metaData = MetaData(preferred_plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE, preferred_timeline_in_months=2.0,
        preferred_payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)

    paymentPlans = gen_payment_plan_builder.createPaymentPlan(paymentTasks=paymentTasks, metaData=metaData)

    assert len(paymentPlans) == 1

    paymentPlan = paymentPlans[0]
    assert paymentPlan.plan_type == PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE
    assert paymentPlan.user_id == user_id
    assert paymentPlan.timeline == 2.0
    assert paymentPlan.payment_freq == PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY
    assert paymentPlan.amount_per_payment == 500
    assert paymentPlan.active == True
    assert paymentPlan.status == PaymentStatus.PAYMENT_STATUS_CURRENT
    assert len(paymentPlan.payment_action) == 2

    assert paymentPlan.payment_action[0].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[0].amount == 500
    assert paymentPlan.payment_action[0].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(datetime.datetime.now(), PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB = Timestamp()
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[0].transaction_date == transactionDatePB

    assert paymentPlan.payment_action[1].account_id == accName2Id['Chase']
    assert paymentPlan.payment_action[1].amount == 500
    assert paymentPlan.payment_action[1].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[1].transaction_date == transactionDatePB

    assert paymentPlan.end_date == transactionDatePB


def test_create_payment_plan_credit_score_no_other_prefs(gen_payment_plan_builder):
    paymentTasks = [
        PaymentTask(user_id=user_id, account_id=accName2Id['Amex'], amount=500), # Amex
        PaymentTask(user_id=user_id, account_id=accName2Id['Chase'], amount=500), # Chase
    ]
    metaData = MetaData(preferred_plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE, preferred_timeline_in_months=0.0,
        preferred_payment_freq=None)

    paymentPlans = gen_payment_plan_builder.createPaymentPlan(paymentTasks=paymentTasks, metaData=metaData)

    # print(paymentPlans)
    # print(len(paymentPlans))

    assert len(paymentPlans) == 3

    ### PaymentPlan 0
    paymentPlan = paymentPlans[0]
    assert paymentPlan.plan_type == PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE
    assert paymentPlan.user_id == user_id
    assert paymentPlan.timeline == 3.0
    assert paymentPlan.payment_freq == PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY
    assert (paymentPlan.amount_per_payment - 333.34) < 1e-2
    print(paymentPlan.amount_per_payment)
    assert paymentPlan.active == True
    assert paymentPlan.status == PaymentStatus.PAYMENT_STATUS_CURRENT
    assert len(paymentPlan.payment_action) == 4
    # PaymentAction 0
    assert paymentPlan.payment_action[0].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[0].amount == 333.34
    assert paymentPlan.payment_action[0].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(datetime.datetime.now(), PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB = Timestamp()
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[0].transaction_date == transactionDatePB
    # PaymentAction 1
    assert paymentPlan.payment_action[1].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[1].amount == 166.66
    assert paymentPlan.payment_action[1].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[1].transaction_date == transactionDatePB
    # PaymentAction 2
    assert paymentPlan.payment_action[1].account_id == accName2Id['Chase']
    assert paymentPlan.payment_action[1].amount == 166.68
    assert paymentPlan.payment_action[1].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    assert paymentPlan.payment_action[1].transaction_date == transactionDatePB
    # PaymentAction 3
    assert paymentPlan.payment_action[1].account_id == accName2Id['Chase']
    assert paymentPlan.payment_action[1].amount == 332.32
    assert paymentPlan.payment_action[1].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[1].transaction_date == transactionDatePB
    #
    assert paymentPlan.end_date == transactionDatePB
    ### PaymentPlan 1
    paymentPlan = paymentPlans[1]
    assert paymentPlan.plan_type == PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE
    assert paymentPlan.user_id == user_id
    assert paymentPlan.timeline == 6.0
    assert paymentPlan.payment_freq == PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY
    assert paymentPlan.amount_per_payment == 166.67
    assert paymentPlan.active == True
    assert paymentPlan.status == PaymentStatus.PAYMENT_STATUS_CURRENT
    assert len(paymentPlan.payment_action) == 4
    # PaymentAction 0
    assert paymentPlan.payment_action[0].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[0].amount == 166.67
    assert paymentPlan.payment_action[0].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(datetime.datetime.now(), PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB = Timestamp()
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[0].transaction_date == transactionDatePB
    # PaymentAction 1
    assert paymentPlan.payment_action[1].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[1].amount == 166.67
    assert paymentPlan.payment_action[1].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[1].transaction_date == transactionDatePB
    # PaymentAction 2
    assert paymentPlan.payment_action[2].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[2].amount == 166.66
    assert paymentPlan.payment_action[2].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[2].transaction_date == transactionDatePB
    # PaymentAction 3
    assert paymentPlan.payment_action[3].account_id == accName2Id['Amex']
    assert paymentPlan.payment_action[3].amount == 0.01
    assert paymentPlan.payment_action[3].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[3].transaction_date == transactionDatePB
    # PaymentAction 4
    assert paymentPlan.payment_action[4].account_id == accName2Id['Chase']
    assert paymentPlan.payment_action[4].amount == 166.67
    assert paymentPlan.payment_action[4].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    assert paymentPlan.payment_action[4].transaction_date == transactionDatePB
    # PaymentAction 5
    assert paymentPlan.payment_action[5].account_id == accName2Id['Chase']
    assert paymentPlan.payment_action[5].amount == 166.67
    assert paymentPlan.payment_action[5].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[5].transaction_date == transactionDatePB
    # PaymentAction 6
    assert paymentPlan.payment_action[6].account_id == accName2Id['Chase']
    assert paymentPlan.payment_action[6].amount == 166.65
    assert paymentPlan.payment_action[6].status == PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING
    transactionDateDt = shift_date_by_payment_frequency(transactionDateDt, PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY)
    transactionDatePB.FromDatetime(transactionDateDt)
    assert paymentPlan.payment_action[6].transaction_date == transactionDatePB
    #
    assert paymentPlan.end_date == transactionDatePB

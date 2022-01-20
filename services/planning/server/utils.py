import pandas as pd
from typing import List
from datetime import timedelta, datetime
from google.protobuf.timestamp_pb2 import Timestamp
import sys, os
# make gen/Python importable by import Python.X
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, os.pardir, 'gen')))
from Python.common.common_pb2 import (PlanType as PlanTypePB, PaymentFrequency as PaymentFrequencyPB)
from Python.planning.payment_plan_pb2 import (PaymentStatus as PaymentStatusPB, PaymentActionStatus as PaymentActionStatusPB, PaymentAction as PaymentActionPB, PaymentPlan as PaymentPlanPB)
# make ../database importable
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir)))
from database.models import (PaymentAction as PaymentActionDB, PaymentPlan as PaymentPlanDB)


def datetime_to_pb_timestamp(timestamp: datetime) -> Timestamp:
      pb_timestamp = Timestamp()
      pb_timestamp.FromDatetime(timestamp)
      return pb_timestamp

def payment_actions_db_to_pb(payment_actions: List[PaymentActionDB]) -> List[PaymentActionPB]:
    results = []
    for payment_action in payment_actions:
        new_payment_action = PaymentActionPB(
            account_id=payment_action.AccountID,
            amount=payment_action.Amount,
            transaction_date=datetime_to_pb_timestamp(payment_action.TransactionDate),
            status=PaymentActionStatusPB.Value(payment_action.PaymentActionStatus.name))
    results.append(new_payment_action)
    return results

def payment_actions_pb_to_db(payment_actions: List[PaymentActionPB]) -> List[PaymentActionDB]:
    results = []
    for payment_action in payment_actions:
        new_payment_action = PaymentActionDB(
            AccountID=payment_action.account_id,
            Amount=payment_action.amount,
            TransactionDate=payment_action.transaction_date.ToDatetime(),
            PaymentActionStatus=payment_action.status)
        results.append(new_payment_action)
    return results

def paymentPlanDBToPB(paymentPlanDB: PaymentPlanDB) -> PaymentPlanPB:
    """ Converts a MongoDB Document version of PaymentPlan to Protobuf version"""
    return PaymentPlanPB(payment_plan_id=paymentPlanDB.PaymentPlanID,
        user_id=paymentPlanDB.UserID,
        payment_task_id=paymentPlanDB.PaymentTaskID, 
        timeline=paymentPlanDB.Timeline,
        payment_freq=PaymentFrequencyPB.Value(paymentPlanDB.PaymentFrequency.name),
        amount_per_payment=paymentPlanDB.AmountPerPayment,
        plan_type=PlanTypePB.Value(paymentPlanDB.PlanType.name),
        end_date=datetime_to_pb_timestamp(paymentPlanDB.EndDate),
        active=paymentPlanDB.Active,
        status=PaymentStatusPB.Value(paymentPlanDB.Status.name),
        payment_action=payment_actions_db_to_pb(paymentPlanDB.PaymentAction))

def paymentPlanPBToDB(paymentPlanPB: PaymentPlanPB) -> PaymentPlanDB:
    """ Converts a Protobuf Document version of PaymentPlan to MongoDB version"""
    return PaymentPlanDB(
        PaymentPlanID=paymentPlanPB.payment_plan_id,
        UserID=paymentPlanPB.user_id,
        PaymentTaskID=paymentPlanPB.payment_task_id, 
        Timeline=paymentPlanPB.timeline,
        PaymentFrequency=paymentPlanPB.payment_freq,
        AmountPerPayment=paymentPlanPB.amount_per_payment,
        PlanType=paymentPlanPB.plan_type,
        EndDate=paymentPlanPB.end_date.ToDatetime(),
        Active=paymentPlanPB.active,
        Status=paymentPlanPB.status,
        PaymentAction=payment_actions_pb_to_db(paymentPlanPB.payment_action)
    )


def shift_date_by_payment_frequency(date: datetime, payment_freq: PaymentFrequencyPB) -> datetime:
    if payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_WEEKLY:
        return date + timedelta(days=7)
    elif payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_BIWEEKLY:
        return date + timedelta(days=14)
    date_format = '%m/%d/%Y'
    date_string = date.strftime(date_format)
    dtObj = datetime.strptime(date_string, date_format)
    if payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_MONTHLY:
        future_date = dtObj + pd.DateOffset(months=1)
    elif payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_QUARTERLY:
        future_date = dtObj + pd.DateOffset(months=3)
    return future_date.date()


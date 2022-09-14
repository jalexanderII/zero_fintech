from datetime import timedelta, datetime
from typing import List

import pandas as pd
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PaymentFrequency as PaymentFrequencyPB
from gen.Python.common.payment_plan_pb2 import PaymentPlan as PaymentPlanPB
from gen.Python.common.payment_plan_pb2 import PaymentAction as PaymentActionPB

from services.planning.database.models.common import PaymentAction as PaymentActionDB
from services.planning.database.models.common import PaymentPlanWName as PaymentPlanDB


def datetime_to_pb_timestamp(timestamp: datetime) -> Timestamp:
    pb_timestamp = Timestamp()
    pb_timestamp.FromDatetime(timestamp)
    return pb_timestamp


def payment_actions_db_to_pb(
    payment_actions: List[PaymentActionDB],
) -> List[PaymentActionPB]:
    results: List[PaymentActionPB] = []
    for payment_action in payment_actions:
        new_payment_action = PaymentActionPB(
            account_id=payment_action.account_id,
            amount=payment_action.amount,
            transaction_date=datetime_to_pb_timestamp(payment_action.transaction_date),
            status=payment_action.status,
        )
        results.append(new_payment_action)
    return results


def payment_actions_pb_to_db(
    payment_actions: List[PaymentActionPB],
) -> List[PaymentActionDB]:
    results: List[PaymentActionDB] = []
    for payment_action in payment_actions:
        new_payment_action = PaymentActionDB(
            account_id=payment_action.account_id,
            amount=payment_action.amount,
            transaction_date=payment_action.transaction_date.ToDatetime(),
            status=payment_action.status,
        )
        results.append(new_payment_action)
    return results


def payment_plan_DB_to_PB(paymentPlanDB: PaymentPlanDB) -> PaymentPlanPB:
    """Converts a MongoDB Document version of PaymentPlan to Protobuf version"""
    return PaymentPlanPB(
        payment_plan_id=paymentPlanDB.payment_plan_id,
        user_id=paymentPlanDB.user_id,
        payment_task_id=paymentPlanDB.payment_task_id,
        amount=paymentPlanDB.amount,
        timeline=paymentPlanDB.timeline,
        payment_freq=paymentPlanDB.payment_freq,
        amount_per_payment=paymentPlanDB.amount_per_payment,
        plan_type=paymentPlanDB.plan_type,
        end_date=datetime_to_pb_timestamp(paymentPlanDB.end_date),
        active=paymentPlanDB.active,
        status=paymentPlanDB.status,
        payment_action=payment_actions_db_to_pb(paymentPlanDB.payment_action),
    )


def payment_plan_PB_to_DB(paymentPlanPB: PaymentPlanPB) -> PaymentPlanDB:
    """Converts a Protobuf Document version of PaymentPlan to MongoDB version"""
    task_ids: List[str] = []
    for task_id in paymentPlanPB.payment_task_id:
        task_ids.append(str(task_id))
    return PaymentPlanDB(
        payment_plan_id=paymentPlanPB.payment_plan_id,
        user_id=paymentPlanPB.user_id,
        payment_task_id=task_ids,
        amount=paymentPlanPB.amount,
        timeline=paymentPlanPB.timeline,
        payment_freq=paymentPlanPB.payment_freq,
        amount_per_payment=paymentPlanPB.amount_per_payment,
        plan_type=paymentPlanPB.plan_type,
        end_date=paymentPlanPB.end_date.ToDatetime(),
        active=paymentPlanPB.active,
        status=paymentPlanPB.status,
        payment_action=payment_actions_pb_to_db(paymentPlanPB.payment_action),
    )


def shift_date_by_payment_frequency(
    date: datetime, payment_freq: PaymentFrequencyPB
) -> datetime:
    future_date: datetime = date
    date_format = "%m/%d/%Y"
    date_string = date.strftime(date_format)
    dtObj: datetime = datetime.strptime(date_string, date_format)

    if payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_WEEKLY:
        return date + timedelta(days=7)
    elif payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_BIWEEKLY:
        return date + timedelta(days=14)
    elif payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_MONTHLY:
        future_date = dtObj + pd.DateOffset(months=1)
    elif payment_freq == PaymentFrequencyPB.PAYMENT_FREQUENCY_QUARTERLY:
        future_date = dtObj + pd.DateOffset(months=3)
    return future_date


def payment_frequency_to_days(paymentFrequency: PaymentFrequencyPB) -> timedelta:
    """Converts PaymentFrequency protobuf to timedelta"""
    payment_freq_days = None
    if paymentFrequency == PaymentFrequencyPB.PAYMENT_FREQUENCY_WEEKLY:
        payment_freq_days = timedelta(days=7)
    elif paymentFrequency == PaymentFrequencyPB.PAYMENT_FREQUENCY_BIWEEKLY:
        payment_freq_days = timedelta(days=14)
    elif paymentFrequency == PaymentFrequencyPB.PAYMENT_FREQUENCY_MONTHLY:
        payment_freq_days = timedelta(days=30)
    elif paymentFrequency == PaymentFrequencyPB.PAYMENT_FREQUENCY_QUARTERLY:
        payment_freq_days = timedelta(days=90)
    return payment_freq_days

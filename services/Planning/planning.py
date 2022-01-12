import numpy as np
import pandas as pd
import grpc
from concurrent import futures
import logging
from datetime import timedelta, datetime

from google.protobuf.timestamp_pb2 import Timestamp

from gen.core import core_pb2_grpc, accounts_pb2, transactions_pb2
from gen.planning import planning_pb2_grpc, payment_plan_pb2, payment_task_pb2, common_pb2

def shift_date_by_payment_frequency(date: datetime, payment_freq: common_pb2.PaymentFrequency) -> timedelta:
    if payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY:
        return date + timedelta(days=7)
    elif payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_BIWEEKLY:
        return date + timedelta(days=14)
    date_format = '%m/%d/%Y'
    dtObj = datetime.strptime(date, date_format)
    if payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY:
        future_date = dtObj + pd.DateOffset(months=1)
    elif payment_freq == common_pb2.PaymentFrequency.PAYMENT_FREQUENCY_QUARTERLY:
        future_date = dtObj + pd.DateOffset(months=1)
    return future_date.date()

class PlanningServicer(planning_pb2_grpc.PlanningServicer):

    def CreatePaymentPlan(self, request, context: grpc.ServicerContext):
        user_id = None
        payment_task_ids = []
        transaction_ids = []
        account_ids = []
        pref_payment_freqs = []
        pref_plan_types = []
        pref_timelines = []
        for _payment_task in request.payment_tasks:
            if user_id is None:
                user_id = _payment_task.user_id
            payment_task_ids.append(_payment_task.payment_task_id)
            transaction_ids.append(_payment_task.transaction_id)
            account_ids.append(_payment_task.account_id)
            pref_payment_freqs.append(_payment_task.preferred_payment_freq)
            pref_plan_types.append(_payment_task.preferred_plan_type)
            pref_timelines.append(_payment_task.preferred_timeline)
        payment_task_ids = np.array(payment_task_ids)
        transaction_ids = np.array(transaction_ids)
        account_ids = np.array(account_ids)

        payment_plans = []
        for _pref_payment_freq, _pref_payment_type, _pref_timeline in set(zip(pref_payment_freqs, pref_plan_types, pref_timelines)):
            mask = pref_payment_freqs == _pref_payment_freq
            payment_plans.append(self._createPaymentPlan(payment_task_ids=payment_task_ids[mask], transaction_ids=transaction_ids[mask], user_id=user_id,
                account_ids=account_ids[mask], pref_payment_freq=_pref_payment_freq, pref_plan_type=_pref_payment_type, pref_timeline=_pref_timeline))

        return planning_pb2_grpc.CreatePaymentPlanResponse(payment_plans=payment_plans)
    
    def _createPaymentPlan(self, payment_task_ids, transaction_ids, account_ids, user_id, pref_payment_freq, pref_plan_type, pref_timeline) -> payment_plan_pb2.PaymentPlan:
        timestamp = Timestamp()
        timestamp.GetCurrentTime()
        start_date = timestamp.ToDatetime()

        channel = grpc.insecure_channel('localhost:50051')
        stub = core_pb2_grpc.CoreStub(channel)
        accounts = []
        account_aprs = []
        for _account_id in account_ids:
            _account = stub.GetAccount(accounts_pb2.GetAccountRequest(id=_account_id))
            accounts.append(_account)
            _apr = (_account.annual_percentage_rate.high_end + _account.annual_percentage_rate.low_end) / 2
            account_aprs.append(_apr)
        account_ids = np.array(account_ids)
        account_aprs = np.array(account_aprs)
        transaction_amounts = []
        for _transaction_id in transaction_ids:
            _transaction = stub.GetTransaction(transactions_pb2.GetTransactionRequest(id=_transaction_id))
            transaction_amounts.append(_transaction.amount)
        transaction_amounts = np.array(transaction_amounts)
        total_amount = sum(transaction_amounts)
        amount_per_payment = total_amount / (pref_payment_freq * pref_timeline)     # TODO: details

        if pref_plan_type == payment_task_pb2.PlanType.PLAN_TYPE_MIN_FEES or pref_plan_type == payment_task_pb2.PlanType.PLANTYPE_UNKNOWN:
            _array_sort = np.argsort(account_aprs)[::-1]
            account_ids = account_ids[_array_sort]
            account_aprs = account_aprs[_array_sort]
            transaction_amounts = transaction_amounts[_array_sort]
            transaction_ids = transaction_ids[_array_sort]

            account_ids = account_ids.tolist()
            transaction_amounts = transaction_amounts.tolist()

            payment_actions = []
            date = shift_date_by_payment_frequency(date=start_date, payment_freq=pref_payment_freq)
            timestamp.FromDatetime(date)
            pay_this_date = 0
            while len(transaction_amounts) > 0:
                _amount = transaction_amounts.pop(0)
                _account_id = account_ids.pop(0)
                _amount_this_date = min(amount_per_payment-pay_this_date, _amount)
                _amount_next_dates = _amount - _amount_this_date
                if _amount_this_date > 0:
                    payment_actions.append(payment_plan_pb2.PaymentAction(account_id=_account_id, amount=_amount, transaction_date=timestamp, status=payment_plan_pb2.PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
                    pay_this_date += _amount
                if _amount_next_dates > 0:
                    transaction_amounts.insert(0, _amount_next_dates)
                    account_ids.insert(0, _account_id)
                    pay_this_date = 0
                    date = shift_date_by_payment_frequency(date=date, payment_freq=pref_payment_freq)
                    timestamp.FromDatetime(date)
                
        elif pref_plan_type == payment_task_pb2.PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE:
            pass

        return payment_plan_pb2(payment_plan_id=1e-9, user_id=user_id, payment_task_id=payment_task_ids, amount_per_payment=amount_per_payment, plan_type=payment_task_pb2.PlanType.PLAN_TYPE_MIN_FEES, end_date=timestamp, active=True, status=payment_plan_pb2.PaymentStatus.PAYMENT_STATUS_CURRENT, payment_action=payment_actions)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    planning_pb2_grpc.add_PlanningServicer_to_server(
        PlanningServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    # logging.basicConfig()
    serve()
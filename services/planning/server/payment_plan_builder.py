from __future__ import annotations

import datetime
from datetime import timedelta
from itertools import product
from math import ceil
from typing import List, Tuple, Generator

import grpc
import pandas as pd
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PaymentFrequency, PlanType, PaymentActionStatus, PaymentStatus
from gen.Python.common.payment_task_pb2 import PaymentTask, MetaData
from gen.Python.common.payment_plan_pb2 import PaymentAction, PaymentPlan
from gen.Python.core.accounts_pb2 import GetAccountRequest
from gen.Python.core.core_pb2_grpc import CoreStub
from services.planning.server.utils import shift_date_by_payment_frequency

PAYMENT_FREQ_TO_TIMELINE = {PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY: 0.25,
                            PaymentFrequency.PAYMENT_FREQUENCY_BIWEEKLY: 0.5,
                            PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY: 1.0,
                            PaymentFrequency.PAYMENT_FREQUENCY_QUARTERLY: 3.0}


class PaymentPlanBuilder:
    def __init__(self, coreClient: CoreStub):
        self.coreClient = coreClient

    def createPaymentPlan(self, paymentTasks: List[PaymentTask], metaData: MetaData=None) -> List[PaymentPlan]:
        """ Creates a PaymentPlan protobuf from given list of PaymentTasks protos."""
        userId = None
        paymentTaskIds = []
        accountIds = []
        amounts = []
        for paymentTask in paymentTasks:
            if userId is None:
                userId = paymentTask.user_id
            paymentTaskIds.append(paymentTask.payment_task_id)
            accountIds.append(paymentTask.account_id)
            amounts.append(paymentTask.amount)

        paymentPlans = []
        for metaData in self._createMetaDataOptions(metaData=metaData, totalAmount=sum(amounts)):
            paymentPlan = self._createPaymentPlanMetaDataChosen(userId=userId, metaData=metaData,
                                paymentTaskIds=paymentTaskIds, accountIds=accountIds, amounts=amounts)
            paymentPlans.append(paymentPlan)

        return paymentPlans

    def _createMetaDataOptions(self, metaData: MetaData, totalAmount: float) -> List[MetaData]:
        """ Creates options for what kind of plans to create given type, timeline, frequency and total amount. """
        if metaData:
            planType = metaData.preferred_plan_type
            timelineInMonths = metaData.preferred_timeline_in_months
            paymentFreq = metaData.preferred_payment_freq
        else:
            planType = PlanType.PLAN_TYPE_UNKNOWN
            timelineInMonths = 0.0
            paymentFreq = PaymentFrequency.PAYMENT_FREQUENCY_UNKNOWN

        typeOptions = [planType] if planType != PlanType.PLAN_TYPE_UNKNOWN else [PlanType.PLAN_TYPE_MIN_FEES,
                                                                                 PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE]
        freqOptions = [paymentFreq]
        if paymentFreq == PaymentFrequency.PAYMENT_FREQUENCY_UNKNOWN:
            if timelineInMonths > 0:
                timelineOptions = [timelineInMonths]
            else:
                if totalAmount < 250.0:
                    timelineOptions = [1.0]
                else:
                    timelineOptions = [3.0, 6.0, 12.0]
            freqOptions = [PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY]
        elif paymentFreq in [PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY, PaymentFrequency.PAYMENT_FREQUENCY_BIWEEKLY]:
            if timelineInMonths > 0:
                timelineOptions = [timelineInMonths]
            else:
                timelineOptions = [1.0, 2.0]
        elif paymentFreq in [PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, PaymentFrequency.PAYMENT_FREQUENCY_QUARTERLY]:
            if timelineInMonths > 0:
                timelineOptions = [timelineInMonths]
            else:
                timelineOptions = [3.0, 6.0, 12.0]

        metaDataOptions = []
        for _type, _timeline, _freq in product(typeOptions, timelineOptions, freqOptions):
            metaDataOptions.append(MetaData(preferred_plan_type=_type, preferred_timeline_in_months=_timeline,
                                            preferred_payment_freq=_freq))
        return metaDataOptions

    def _createPaymentPlanMetaDataChosen(self, userId: str, metaData: MetaData, paymentTaskIds: List[str],
            accountIds: List[str], amounts: List[float]) -> PaymentPlan:
        """Creates a PaymentPlan for given choices of MetaData."""
        planType = metaData.preferred_plan_type
        timelineInMonths = metaData.preferred_timeline_in_months
        paymentFreq = metaData.preferred_payment_freq

        if planType == PlanType.PLAN_TYPE_UNKNOWN:
            raise ValueError("Using PLAN_TYPE_UNKNOWN not permitted")
        if timelineInMonths <= 0:
            raise ValueError("Need to specify timeline > 0")
        if paymentFreq == PaymentFrequency.PAYMENT_FREQUENCY_UNKNOWN:
            raise ValueError("Using PAYMENT_FREQUENCY_UNKNOWN not permitted")

        totalAmount = sum(amounts)
        amountPerPayment = totalAmount / (timelineInMonths / PAYMENT_FREQ_TO_TIMELINE[paymentFreq])
        amountPerPayment = round(ceil(amountPerPayment * 100) / 100, 2)

        if planType == PlanType.PLAN_TYPE_MIN_FEES:
            apr = []
            for accId in accountIds:
                acc = self.coreClient.GetAccount(GetAccountRequest(id=accId))
                apr.append((acc.annual_percentage_rate.high_end + acc.annual_percentage_rate.low_end) / 2)
            df = pd.DataFrame({'account_id': accountIds, 'apr': apr, 'amount': amounts,
                               'payment_task_id': paymentTaskIds}).sort_values('apr', ascending=False)
            accountIds = df['account_id'].values.tolist()
            amounts = df['amount'].values.tolist()
            paymentTaskIds = df['payment_task_id'].values.tolist()

            paymentActions = []
            startDate = datetime.datetime.now()
            dateDt = shift_date_by_payment_frequency(date=startDate, payment_freq=paymentFreq)
            datePB = Timestamp()
            datePB.FromDatetime(dateDt)
            payThisDate = 0
            while len(amounts) > 0:
                _amount = amounts.pop(0)
                _accId = accountIds.pop(0)
                _amountThisDate = min(amountPerPayment - payThisDate, _amount)
                _amountNextDates = _amount - _amountThisDate
                if _amountThisDate > 0:
                    paymentActions.append(PaymentAction(account_id=_accId, amount=_amountThisDate, transaction_date=datePB,
                                                        status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
                    payThisDate += _amount
                if _amountNextDates > 0:
                    amounts.insert(0, _amountNextDates)
                    accountIds.insert(0, _accId)
                    # move to next date
                    payThisDate = 0
                    dateDt = shift_date_by_payment_frequency(date=dateDt, payment_freq=paymentFreq)
                    datePB.FromDatetime(dateDt)
        elif planType == PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE:
            balance = []
            creditLimit = []
            for accId in accountIds:
                acc = self.coreClient.GetAccount(GetAccountRequest(id=accId))
                balance.append(acc.current_balance)
                creditLimit.append(acc.credit_limit)
            df = pd.DataFrame({'account_id': accountIds, 'amount': amounts, 'balance': balance, 'credit_limit': creditLimit})
            df['usage'] = df['balance'] / df['credit_limit']
            df = df.sort_values('usage', ascending=False)

            paymentActions = []
            startDate = datetime.datetime.now()
            dateDt = shift_date_by_payment_frequency(date=startDate, payment_freq=paymentFreq)
            datePB = Timestamp()
            datePB.FromDatetime(dateDt)
            payThisDate = 0
            while len(df) > 0:
                # pay off until hit limit for this payment date
                for _, row in df.iterrows():
                    _accId = row['account_id']
                    _amount = row['amount']
                    _balance = row['balance']
                    _amountThisDate = min(amountPerPayment - payThisDate, _amount)
                    if _amountThisDate == 0:
                        break
                    paymentActions.append(PaymentAction(account_id=_accId, amount=_amountThisDate, transaction_date=datePB,
                                      status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
                    payThisDate += _amountThisDate
                    df.loc[df['account_id'] == _accId, 'amount'] = _amount - _amountThisDate
                    df.loc[df['account_id'] == _accId, 'balance'] = _balance - _amountThisDate
                # drop any accounts which don't have any amounts left to pay off
                df = df.loc[df.amount > 0]
                # update credit card usage
                df['usage'] = df['balance'] / df['credit_limit']
                df = df.sort_values('usage', ascending=False)
                # move to next date
                datePB.FromDatetime(dateDt)
                payThisDate = 0
                dateDt = shift_date_by_payment_frequency(date=dateDt, payment_freq=paymentFreq)
                datePB.FromDatetime(dateDt)

        endDatePB = paymentActions[-1].transaction_date

        return PaymentPlan(user_id=userId, payment_task_id=paymentTaskIds, timeline=timelineInMonths, payment_freq=paymentFreq,
                           amount_per_payment=amountPerPayment, plan_type=planType, end_date=endDatePB, active=True,
                           status=PaymentStatus.PAYMENT_STATUS_CURRENT, payment_action=paymentActions)

payment_plan_builder = PaymentPlanBuilder(coreClient=CoreStub(channel=grpc.insecure_channel('localhost:9090')))

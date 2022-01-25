from __future__ import annotations

import datetime
from datetime import timedelta
from itertools import product
from typing import List, Tuple

import numpy as np
import pandas as pd
from bson.objectid import ObjectId
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.core.core_pb2_grpc import  CoreStub
from gen.Python.core.accounts_pb2 import GetAccountRequest
from gen.Python.common.common_pb2 import PaymentFrequency, PlanType
from gen.Python.common.payment_task_pb2 import PaymentTask, MetaData
from gen.Python.planning.payment_plan_pb2 import PAYMENT_STATUS_CURRENT, PaymentActionStatus, PaymentStatus
from gen.Python.planning.payment_plan_pb2 import PaymentAction, PaymentPlan
from services.planning.server.utils import paymentFrequencyToDays, shift_date_by_payment_frequency


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

        if metaData:
            planType = metaData.preferred_plan_type
            timelineInMonths = metaData.timeline_in_months
        else:
            planType = PlanType.PLAN_TYPE_UNKNOWN
            timelineInMonths = 0.0
            paymentFreq = PaymentFrequency.PAYMENT_FREQUENCY_UNKNOWN

        typeTimelineFrequencyOptions = self._createTypeTimelineFrequencyOptions(planType=planType,
            timelineInMonths=timelineInMonths, paymentFreq=paymentFreq, totalAmount=sum(amounts))

        paymentPlans = []
        for prefPlanType, prefTimelineInMonths, prefPaymentFreq in typeTimelineFrequencyOptions:
            paymentPlan = self._createPaymentPlanSameTypeFreq(userId=userId, planType=prefPlanType,
                timelineInMonths=prefTimelineInMonths, paymentFreq=prefPaymentFreq, paymentTaskIds=paymentTaskIds,
                accountIds=accountIds, amounts=amounts)
            paymentPlans.append(paymentPlan)

        return paymentPlans

    def _createTypeTimelineFrequencyOptions(self, planType: PlanType, timelineInMonths: float, paymentFreq: PaymentFrequency,
            totalAmount: float) -> List[Tuple[PlanType, float, PaymentFrequency]]:
        """ Creates options for what kind of plans to create given type, timeline, frequency and total amount. """
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
        return list(product(typeOptions, timelineOptions, freqOptions))

    def _createPaymentPlanSameTypeFreq(self, userId: str, planType: PlanType, timelineInMonths: float,
            paymentFreq: PaymentFrequency, paymentTaskIds: List[str], accountIds: List[str], amounts: List[float],
            ) -> PaymentPlan:
        """Creates a PaymentPlan for given choices of PaymentFrequency and PlanType for the Accounts and Amounts. """
        assert planType != PlanType.PLAN_TYPE_UNKNOWN, "PlanType needs to be chosen"

        totalAmount = sum(amounts)
        amountPerPayment = totalAmount / (timedelta(days=30) * timelineInMonths / paymentFrequencyToDays(paymentFreq))

        if planType == PlanType.PLAN_TYPE_MIN_FEES:
            apr = []
            for accId in accountIds:
                acc = self.coreClient.GetAccount(GetAccountRequest(id=accId))
                apr.append((acc.annual_percentage_rate.high_end + acc.annual_percentage_rate.low_end) / 2)
            df = pd.DataFrame({'account_id': accountIds, 'apr': apr, 'amount': amounts,
                               'payment_task_id': paymentTaskIds}).sort_values('apr', ascending=False)
            accountIds = df['account_id'].values
            amounts = df['amount'].values
            paymentTaskIds = df['payment_task_id'].values

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
                elif _amountNextDates > 0:
                    amounts.insert(0, _amountNextDates)
                    accountIds.insert(0, _accId)
                    # move to next date
                    payThisDate = 0
                    dateDt = shift_date_by_payment_frequency(date=dateDt, payment_freq=paymentFreq)
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
                    paymentActions.append(PaymentAction(account_id=_accId, amount=_amountThisDate, transaction_date=datePB,
                                      status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING))
                    payThisDate += _amountThisDate
                    df.loc[df['account_id'] == _accId, 'amount'] = _amount - _amountThisDate
                    df.loc[df['account_id'] == _accId, 'balance'] = _balance - _amountThisDate
                # drop any accounts which don't have any amounts left to pay off
                df = df.loc[df.amount > 0]
                # update credit card usage
                df['usage'] = df['balance'] / df['credit_limit']
                # move to next date
                datePB.FromDatetime(dateDt)
                payThisDate = 0
                dateDt = shift_date_by_payment_frequency(date=dateDt, payment_freq=paymentFreq)
                datePB.FromDatetime(dateDt)

        endDatePB = paymentActions[-1].transaction_date
        # calculate timeline from how many dates were chosen
        if paymentFreq == PaymentFrequency.PAYMENT_FREQUENCY_WEEKLY:
            paymentFreqAsTimelineInMonths = 0.25
        elif paymentFreq == PaymentFrequency.PAYMENT_FREQUENCY_BIWEEKLY:
            paymentFreqAsTimelineInMonths = 0.5
        elif paymentFreq == PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY:
            paymentFreqAsTimelineInMonths = 1.0
        elif paymentFreq == PaymentFrequency.PAYMENT_FREQUENCY_QUARTERLY:
            paymentFreqAsTimelineInMonths = 3.0
        timeline = paymentFreqAsTimelineInMonths * len(set([pa.transaction_date for pa in paymentActions]))

        return PaymentPlan(user_id=userId, payment_task_id=paymentTaskIds, timeline=timeline, payment_freq=paymentFreq,
                           amount_per_payment=amountPerPayment, plan_type=planType, end_date=endDatePB, active=True,
                           status=PaymentStatus.PAYMENT_STATUS_CURRENT, payment_action=paymentActions)

    @staticmethod
    def _mock_payment_plan(paymentTasks: List[PaymentTask]) -> List[PaymentPlan]:
        """A mock payment plan for end-to-end testing"""
        ids: List[str] = []
        plans: List[PaymentPlan] = []
        actions: List[PaymentAction] = []
        total = 0
        pb_timestamp = Timestamp()
        for task in paymentTasks:
            ids.append(task.payment_task_id)
            total += task.amount
            a = PaymentAction(
                account_id=task.account_id,
                amount=task.amount,
                transaction_date=pb_timestamp.GetCurrentTime(),
                status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
            )
            actions.append(a)
        bson_id = ObjectId()
        plan = PaymentPlan(
            payment_plan_id=str(bson_id),
            user_id=paymentTasks[0].user_id,
            payment_task_id=ids,
            timeline=12,
            payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
            amount_per_payment=total / 12,
            plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE,
            end_date=pb_timestamp.GetCurrentTime(),
            active=True,
            status=PAYMENT_STATUS_CURRENT,
            payment_action=actions,
        )
        plans.append(plan)
        return plans


# payment_plan_builder = PaymentPlanBuilder()

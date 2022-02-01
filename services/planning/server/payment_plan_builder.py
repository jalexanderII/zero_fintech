from __future__ import annotations

import datetime
import logging
import os
import sys
from itertools import product
from math import ceil
from typing import List, Dict

import grpc
import pandas as pd
from attr import define
from dotenv import load_dotenv
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PaymentFrequency, PlanType, PaymentActionStatus, PaymentStatus, \
    PAYMENT_FREQUENCY_UNKNOWN, PAYMENT_FREQUENCY_WEEKLY, PAYMENT_FREQUENCY_QUARTERLY, PAYMENT_FREQUENCY_MONTHLY, \
    PAYMENT_FREQUENCY_BIWEEKLY, PLAN_TYPE_OPTIM_CREDIT_SCORE, PLAN_TYPE_MIN_FEES, PLAN_TYPE_UNKNOWN, \
    PAYMENT_ACTION_STATUS_PENDING, PAYMENT_STATUS_CURRENT
from gen.Python.common.payment_plan_pb2 import PaymentAction, PaymentPlan
from gen.Python.common.payment_task_pb2 import PaymentTask, MetaData
from gen.Python.core.accounts_pb2 import GetAccountRequest, Account, AnnualPercentageRates
from gen.Python.core.core_pb2_grpc import CoreStub
from services.planning.server.utils import shift_date_by_payment_frequency

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger("PaymentPlanBuilder")

load_dotenv()
CORE_CLIENT = CoreStub(grpc.insecure_channel(f'localhost:{os.getenv("CORE_SERVER_PORT")}'))


@define
class PaymentPlanBuilder:
    coreClient: CoreStub = CORE_CLIENT

    def create(self, paymentTasks: List[PaymentTask], metaData: MetaData = None) -> List[PaymentPlan]:
        """ Creates a PaymentPlan protobuf from given list of PaymentTasks protos."""
        userId: str = paymentTasks[-1].user_id
        paymentTaskIds: List[str] = []
        accountIds: List[str] = []
        amounts: List[float] = []
        for payment_task in paymentTasks:
            paymentTaskIds.append(payment_task.payment_task_id)
            accountIds.append(payment_task.account_id)
            amounts.append(payment_task.amount)

        return [
            self._create_from_meta_data(
                userId=userId, metaData=metaData,
                paymentTaskIds=paymentTaskIds, accountIds=accountIds,
                amounts=amounts
            )
            for metaData in self._get_meta_data_options(metaData=metaData, totalAmount=sum(amounts) > 250.0)
        ]

    @staticmethod
    def _get_meta_data_options(metaData: MetaData, totalAmount: bool) -> List[MetaData]:
        """ Creates options for what kind of plans to create given type, timeline, frequency and total amount. """
        timelineInMonths = 0.0
        paymentFreq = PAYMENT_FREQUENCY_UNKNOWN
        typeOptions = [PlanType.PLAN_TYPE_MIN_FEES, PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE]
        freqOptions: List[PaymentFrequency] = [PAYMENT_FREQUENCY_MONTHLY]
        if metaData:
            planType = metaData.preferred_plan_type
            timelineInMonths = metaData.preferred_timeline_in_months
            paymentFreq = metaData.preferred_payment_freq
            typeOptions = [planType] if planType != PLAN_TYPE_UNKNOWN else typeOptions
            freqOptions = [paymentFreq] if paymentFreq != PAYMENT_FREQUENCY_UNKNOWN else freqOptions

        payment_freq_to_timeline_options: Dict[PaymentFrequency, List[float]] = {
            PAYMENT_FREQUENCY_UNKNOWN: [3.0, 6.0, 12.0] if totalAmount else [1.0],
            PAYMENT_FREQUENCY_WEEKLY: [1.0, 2.0],
            PAYMENT_FREQUENCY_BIWEEKLY: [1.0, 2.0],
            PAYMENT_FREQUENCY_MONTHLY: [3.0, 6.0, 12.0],
            PAYMENT_FREQUENCY_QUARTERLY: [3.0, 6.0, 12.0],
        }
        timelineOptions = [timelineInMonths] if timelineInMonths > 0 else payment_freq_to_timeline_options[paymentFreq]

        return [
            MetaData(preferred_plan_type=_type, preferred_timeline_in_months=_timeline,
                     preferred_payment_freq=_freq)
            for _type, _timeline, _freq in product(typeOptions, timelineOptions, freqOptions)
        ]

    def _create_from_meta_data(self, userId: str, metaData: MetaData, paymentTaskIds: List[str],
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

        datePB = Timestamp()
        start_date = datetime.datetime.now()

        payment_actions: List[PaymentAction] = []

        balance: List[float] = []
        creditLimit: List[float] = []
        apr: List[float] = []
        accounts: List[Account] = [
            self.coreClient.GetAccount(GetAccountRequest(id=acc_id))
            for acc_id in accountIds
        ]
        logger.info(accounts)
        for acc in accounts:
            apr.append(avg_apr(acc.annual_percentage_rate))
            balance.append(acc.current_balance)
            creditLimit.append(acc.credit_limit)

        df = pd.DataFrame(
            {
                "account_id": accountIds,
                "apr": apr,
                "amount": amounts,
                "payment_task_id": paymentTaskIds,
                "balance": balance,
                "credit_limit": creditLimit,
            }
        )

        PAYMENT_FREQ_TO_TIMELINE: Dict[PaymentFrequency, float] = {
            PAYMENT_FREQUENCY_WEEKLY: 0.25,
            PAYMENT_FREQUENCY_BIWEEKLY: 0.5,
            PAYMENT_FREQUENCY_MONTHLY: 1.0,
            PAYMENT_FREQUENCY_QUARTERLY: 3.0,
        }

        num_payments = timelineInMonths / PAYMENT_FREQ_TO_TIMELINE[paymentFreq]
        amount_per_payment = sum(amounts) / num_payments
        amount_per_payment = round(ceil(amount_per_payment * 100) / 100, 2)

        if planType == PLAN_TYPE_MIN_FEES:
            df.sort_values("apr", ascending=False)
            accountIds = df["account_id"].values.tolist()
            amounts = df["amount"].values.tolist()
            paymentTaskIds = df["payment_task_id"].values.tolist()

            shifted_date = shift_date_by_payment_frequency(start_date, paymentFreq)
            datePB.FromDatetime(shifted_date)

            pay_on_date = 0
            while len(amounts) > 0:
                _amount = amounts.pop(0)
                _accId = accountIds.pop(0)
                _amountThisDate = min(amount_per_payment - pay_on_date, _amount)
                _amountNextDates = _amount - _amountThisDate
                if _amountThisDate > 0:
                    payment_actions.append(
                        PaymentAction(
                            account_id=_accId,
                            amount=_amountThisDate,
                            transaction_date=datePB,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                        )
                    )
                    pay_on_date += _amount
                if _amountNextDates > 0:
                    amounts.insert(0, _amountNextDates)
                    accountIds.insert(0, _accId)
                    # move to next date
                    pay_on_date = 0
                    shifted_date = shift_date_by_payment_frequency(
                        shifted_date, paymentFreq
                    )
                    datePB.FromDatetime(shifted_date)
        elif planType == PLAN_TYPE_OPTIM_CREDIT_SCORE:
            df["usage"] = df["balance"] / df["credit_limit"]
            df = df.sort_values("usage", ascending=False)

            shifted_date = shift_date_by_payment_frequency(start_date, paymentFreq)
            datePB.FromDatetime(shifted_date)
            pay_on_date = 0
            while len(df) > 0:
                # pay off until hit limit for this payment date
                for _, row in df.iterrows():
                    _accId = row["account_id"]
                    _amount = row["amount"]
                    _balance = row["balance"]
                    _amountThisDate = min(amount_per_payment - pay_on_date, _amount)
                    if _amountThisDate == 0:
                        break
                    payment_actions.append(
                        PaymentAction(
                            account_id=_accId,
                            amount=_amountThisDate,
                            transaction_date=datePB,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                        )
                    )
                    pay_on_date += _amountThisDate
                    remaining_amount = _amount - _amountThisDate
                    remaining_balance = _balance - _amountThisDate
                    df.loc[df["account_id"] == _accId, "amount"] = remaining_amount
                    df.loc[df["account_id"] == _accId, "balance"] = remaining_balance

                # drop any accounts which don't have any amounts left to pay off
                df = df.loc[df.amount > 0]
                # update credit card usage
                df["usage"] = df["balance"] / df["credit_limit"]
                df = df.sort_values("usage", ascending=False)
                # move to next date
                pay_on_date = 0
                shifted_date = shift_date_by_payment_frequency(
                    shifted_date, paymentFreq
                )
                datePB.FromDatetime(shifted_date)

        return PaymentPlan(
            user_id=userId,
            payment_task_id=paymentTaskIds,
            timeline=timelineInMonths,
            payment_freq=paymentFreq,
            amount_per_payment=amount_per_payment,
            plan_type=planType,
            end_date=payment_actions[-1].transaction_date,
            active=True,
            status=PAYMENT_STATUS_CURRENT,
            payment_action=payment_actions,
        )


def avg_apr(apr: AnnualPercentageRates) -> float:
    return (apr.high_end + apr.low_end) / 2


payment_plan_builder = PaymentPlanBuilder()

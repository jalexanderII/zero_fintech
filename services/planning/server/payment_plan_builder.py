from __future__ import annotations

import datetime
import logging
import os
import sys
from itertools import product
from math import ceil
from typing import List, Dict, Optional

import grpc
import pandas as pd
from attr import define
from dotenv import load_dotenv
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PAYMENT_ACTION_STATUS_PENDING
from gen.Python.common.common_pb2 import PAYMENT_STATUS_CURRENT
from gen.Python.common.common_pb2 import (
    PaymentFrequency,
    PAYMENT_FREQUENCY_UNKNOWN,
    PAYMENT_FREQUENCY_WEEKLY,
    PAYMENT_FREQUENCY_BIWEEKLY,
    PAYMENT_FREQUENCY_MONTHLY,
    PAYMENT_FREQUENCY_QUARTERLY,
)
from gen.Python.common.common_pb2 import (
    PlanType,
    PLAN_TYPE_UNKNOWN,
    PLAN_TYPE_MIN_FEES,
    PLAN_TYPE_OPTIM_CREDIT_SCORE,
)
from gen.Python.common.payment_plan_pb2 import PaymentAction, PaymentPlan
from gen.Python.common.payment_task_pb2 import PaymentTask, MetaData
from gen.Python.core.accounts_pb2 import AnnualPercentageRates
from gen.Python.core.accounts_pb2 import GetAccountRequest, Account
from gen.Python.core.core_pb2_grpc import CoreStub
from services.planning.server.utils import shift_date_by_payment_frequency

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger("PaymentPlanBuilder")

load_dotenv()

CORE_CLIENT = CoreStub(
    grpc.insecure_channel(f'localhost:{os.getenv("CORE_SERVER_PORT")}')
)

AMOUNT_THRESHOLD = 250.0
DEFAULT_APR = 10.0


@define
class PaymentPlanBuilder:
    core_client: CoreStub = CORE_CLIENT

    def create(
        self, payment_tasks: List[PaymentTask], meta_data: Optional[MetaData] = None
    ) -> List[PaymentPlan]:
        """Creates a PaymentPlan given a list of PaymentTasks"""
        user_id: str = payment_tasks[0].user_id
        payment_task_ids: List[str] = []
        account_ids: List[str] = []
        amounts: List[float] = []
        for payment_task in payment_tasks:
            payment_task_ids.append(payment_task.payment_task_id)
            account_ids.append(payment_task.account_id)
            amounts.append(payment_task.amount)

        return [
            self._create_from_meta_data(
                user_id=user_id,
                plan_type=_meta_data.preferred_plan_type,
                timeline_months=_meta_data.preferred_timeline_in_months,
                payment_freq=_meta_data.preferred_payment_freq,
                payment_task_ids=payment_task_ids,
                account_ids=account_ids,
                amounts=amounts,
            )
            for _meta_data in self._get_meta_data_options(
                meta_data, sum(amounts) > AMOUNT_THRESHOLD
            )
        ]

    @staticmethod
    def _get_meta_data_options(
        meta_data: Optional[MetaData], gt_threshold: bool
    ) -> List[MetaData]:
        """Creates options for what kind of plans to create given type, timeline, frequency and total amount."""
        timeline_months: float = 0.0
        payment_freq: PaymentFrequency = PAYMENT_FREQUENCY_UNKNOWN
        plan_type_options: List[PlanType] = [
            PLAN_TYPE_MIN_FEES,
            PLAN_TYPE_OPTIM_CREDIT_SCORE,
        ]
        freq_options: List[PaymentFrequency] = [PAYMENT_FREQUENCY_MONTHLY]

        if meta_data:
            plan_type = meta_data.preferred_plan_type
            timeline_months = meta_data.preferred_timeline_in_months
            payment_freq = meta_data.preferred_payment_freq
            plan_type_options = (
                [plan_type] if plan_type != PLAN_TYPE_UNKNOWN else plan_type_options
            )
            freq_options = (
                [payment_freq]
                if payment_freq != PAYMENT_FREQUENCY_UNKNOWN
                else freq_options
            )

        payment_freq_to_timeline_options: Dict[PaymentFrequency, List[float]] = {
            PAYMENT_FREQUENCY_UNKNOWN: [3.0, 6.0, 12.0] if gt_threshold else [1.0],
            PAYMENT_FREQUENCY_WEEKLY: [1.0, 2.0],
            PAYMENT_FREQUENCY_BIWEEKLY: [1.0, 2.0],
            PAYMENT_FREQUENCY_MONTHLY: [3.0, 6.0, 12.0],
            PAYMENT_FREQUENCY_QUARTERLY: [3.0, 6.0, 12.0],
        }

        timeline_options: List[float] = (
            [timeline_months]
            if timeline_months > 0
            else payment_freq_to_timeline_options[payment_freq]
        )

        if PLAN_TYPE_UNKNOWN in plan_type_options:
            raise ValueError(
                f"Chosen PlanTypes contains PLAN_TYPE_UNKNOWN: {plan_type_options}"
            )
        elif len(list(filter(lambda x: x <= 0.0, timeline_options))) > 0:
            raise ValueError(
                f"Chosen timeline options contain non-positive values: {timeline_options}"
            )
        elif PAYMENT_FREQUENCY_UNKNOWN in freq_options:
            raise ValueError(
                f"Chosen PaymentFrequency options contains PAYMENT_FREQUENCY_UNKNOWN: {freq_options}"
            )

        return [
            MetaData(
                preferred_plan_type=pref_type,
                preferred_timeline_in_months=timeline,
                preferred_payment_freq=freq,
            )
            for pref_type, timeline, freq in product(
                plan_type_options, timeline_options, freq_options
            )
        ]

    def _create_from_meta_data(
        self,
        user_id: str,
        plan_type: PlanType,
        timeline_months: float,
        payment_freq: PaymentFrequency,
        payment_task_ids: List[str],
        account_ids: List[str],
        amounts: List[float],
    ) -> PaymentPlan:
        """Creates a PaymentPlan for given choices of MetaData."""
        start_date = datetime.datetime.now()

        balance: List[float] = []
        credit_limit: List[float] = []
        apr: List[float] = []
        accounts: List[Account] = self._fetch_accounts(account_ids)

        for acc in accounts:
            apr.append(get_purchase_apr(acc.annual_percentage_rate))
            balance.append(acc.current_balance)
            credit_limit.append(acc.credit_limit)

        df = pd.DataFrame(
            {
                "account_id": account_ids,
                "apr": apr,
                "amount": amounts,
                "payment_task_id": payment_task_ids,
                "balance": balance,
                "credit_limit": credit_limit,
            }
        )

        PAYMENT_FREQ_TO_TIMELINE: Dict[PaymentFrequency, float] = {
            PAYMENT_FREQUENCY_WEEKLY: 0.25,
            PAYMENT_FREQUENCY_BIWEEKLY: 0.5,
            PAYMENT_FREQUENCY_MONTHLY: 1.0,
            PAYMENT_FREQUENCY_QUARTERLY: 3.0,
        }

        num_payments = timeline_months / PAYMENT_FREQ_TO_TIMELINE[payment_freq]
        amount_per_payment = sum(amounts) / num_payments
        amount_per_payment = round(ceil(amount_per_payment * 100) / 100, 2)

        payment_actions: List[PaymentAction] = []
        if plan_type == PLAN_TYPE_MIN_FEES:
            payment_actions = self._create_payment_actions_min_fees(
                payment_freq=payment_freq,
                df=df,
                start_date=start_date,
                amount_per_payment=amount_per_payment,
            )
        elif plan_type == PLAN_TYPE_OPTIM_CREDIT_SCORE:
            payment_actions = self._create_payment_actions_optim_credit_score(
                payment_freq=payment_freq,
                df=df,
                start_date=start_date,
                amount_per_payment=amount_per_payment,
            )

        return PaymentPlan(
            user_id=user_id,
            payment_task_id=payment_task_ids,
            amount=sum(amounts),
            timeline=timeline_months,
            payment_freq=payment_freq,
            amount_per_payment=amount_per_payment,
            plan_type=plan_type,
            end_date=payment_actions[-1].transaction_date,
            active=True,
            status=PAYMENT_STATUS_CURRENT,
            payment_action=payment_actions,
        )

    def _fetch_accounts(self, account_ids: List[str]) -> List[Account]:
        return [
            self.core_client.GetAccount(GetAccountRequest(id=acc_id))
            for acc_id in account_ids
        ]

    @staticmethod
    def _create_payment_actions_min_fees(
        payment_freq: PaymentFrequency,
        df: pd.DataFrame,
        start_date: datetime.datetime,
        amount_per_payment: float,
    ) -> List[PaymentAction]:
        """Creates a list of PaymentActions to minimize fees for given freq, start date, amount per payment and
        DataFrame on the accounts."""
        payment_actions: List[PaymentAction] = []
        datePB = Timestamp()

        df.sort_values("apr", ascending=False)
        account_ids = df["account_id"].values.tolist()
        amounts = df["amount"].values.tolist()

        shifted_date = shift_date_by_payment_frequency(start_date, payment_freq)
        datePB.FromDatetime(shifted_date)

        pay_on_date = 0
        while len(amounts) > 0:
            _amount = amounts.pop(0)
            _accId = account_ids.pop(0)
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
                account_ids.insert(0, _accId)
                # move to next date
                pay_on_date = 0
                shifted_date = shift_date_by_payment_frequency(
                    shifted_date, payment_freq
                )
                datePB.FromDatetime(shifted_date)

        return payment_actions

    @staticmethod
    def _create_payment_actions_optim_credit_score(
        payment_freq: PaymentFrequency,
        df: pd.DataFrame,
        start_date: datetime.datetime,
        amount_per_payment: float,
    ) -> List[PaymentAction]:
        """Creates a list of PaymentActions to optimize credit score for given freq, start date, amount per payment and
        DataFrame on the accounts."""
        payment_actions: List[PaymentAction] = []
        datePB = Timestamp()

        df["usage"] = df["balance"] / df["credit_limit"]
        df = df.sort_values("usage", ascending=False)

        shifted_date = shift_date_by_payment_frequency(start_date, payment_freq)
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
            shifted_date = shift_date_by_payment_frequency(shifted_date, payment_freq)
            datePB.FromDatetime(shifted_date)

        return payment_actions


def get_purchase_apr(aprs: List[AnnualPercentageRates]) -> float:
    for apr in aprs:
        if apr.apr_type == "purchase_apr":
            return apr.apr_percentage
    return DEFAULT_APR


payment_plan_builder = PaymentPlanBuilder()

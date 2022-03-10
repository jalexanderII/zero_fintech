import datetime
from typing import List

import attr
import pandas as pd

from gen.Python.common.common_pb2 import PlanType, PaymentFrequency
from gen.Python.common.payment_plan_pb2 import PaymentPlan, PaymentAction
from gen.Python.common.payment_task_pb2 import MetaData


@attr.s(auto_attribs=True, kw_only=True)
class ParamSet:
    """
    Use as base class for sets for parameters for @pytest.mark.parametrize
    when you find yourself passing too many parameters positionally.
    """

    id: str

    @property
    def __name__(self) -> str:  # indicate the id to pytest
        return self.id


@attr.s(auto_attribs=True, kw_only=True)
class MetaDataToPaymentPlanParams(ParamSet):
    user_id: str
    plan_type: PlanType
    timeline_months: float
    payment_freq: PaymentFrequency
    payment_task_ids: List[str]
    account_ids: List[str]
    amounts: List[float]
    expected: PaymentPlan


@attr.s(auto_attribs=True, kw_only=True)
class CreatePaymentActionsParams(ParamSet):
    payment_freq: PaymentFrequency
    df: pd.DataFrame
    start_date: datetime.datetime
    amount_per_payment: float
    expected: List[PaymentAction]
    min_fee: bool


@attr.s(auto_attribs=True, kw_only=True)
class GetMetaDataParams(ParamSet):
    meta_data: MetaData
    gt_threshold: bool
    expected: List[MetaData]

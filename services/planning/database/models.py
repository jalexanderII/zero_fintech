# from datetime import datetime
# from enum import IntEnum
# from typing import List, Optional
#
# from pydantic import BaseModel, Field
#
#
# class PlanType(IntEnum):
#     PLAN_TYPE_UNKNOWN = 0
#     PLAN_TYPE_OPTIM_CREDIT_SCORE = 1
#     PLAN_TYPE_MIN_FEES = 2
#
#
# class PaymentStatus(IntEnum):
#     PAYMENT_STATUS_UNKNOWN = 0
#     PAYMENT_STATUS_CURRENT = 1
#     PAYMENT_STATUS_COMPLETED = 2
#     PAYMENT_STATUS_IN_DEFAULT = 4
#
#
# class PaymentActionStatus(IntEnum):
#     PAYMENT_ACTION_STATUS_UNKNOWN = 0
#     PAYMENT_ACTION_STATUS_PENDING = 1
#     PAYMENT_ACTION_STATUS_COMPLETED = 2
#     PAYMENT_ACTION_STATUS_IN_DEFAULT = 3
#
#
# class PaymentFrequency(IntEnum):
#     PAYMENT_FREQUENCY_UNKNOWN = 0
#     PAYMENT_FREQUENCY_WEEKLY = 1
#     PAYMENT_FREQUENCY_BIWEEKLY = 2
#     PAYMENT_FREQUENCY_MONTHLY = 3
#     PAYMENT_FREQUENCY_QUARTERLY = 4
#
#
# class PaymentAction(BaseModel):
#     AccountID: str = Field(...)
#     Amount: float = Field(...)
#     TransactionDate: datetime = None
#     PaymentActionStatus: 'PaymentActionStatus' = None
#
#
# class PaymentPlan(BaseModel):
#     PaymentPlanID: Optional[str] = Field(...)
#     UserID: str = Field(...)
#     PaymentTaskID: List[str] = None
#     Timeline: float = Field(...)
#     PaymentFrequency: 'PaymentFrequency' = None
#     AmountPerPayment: float = Field(...)
#     PlanType: 'PlanType' = None
#     EndDate: datetime = None
#     Active: bool = True
#     Status: 'PaymentStatus' = None
#     PaymentAction: List['PaymentAction'] = None

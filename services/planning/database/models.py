from enum import Enum
from mongoengine import Document, ListField, StringField, FloatField, EmbeddedDocument, EnumField, BooleanField, ObjectIdField, DateTimeField, EmbeddedDocumentField

class PlanType(Enum):
    PLAN_TYPE_UNKNOWN = 0
    PLAN_TYPE_OPTIM_CREDIT_SCORE = 1
    PLAN_TYPE_MIN_FEES = 2

class PaymentStatus(Enum):
    PAYMENT_STATUS_UNKNOWN=0
    PAYMENT_STATUS_CURRENT=1
    PAYMENT_STATUS_COMPLETED=2
    PAYMENT_STATUS_IN_DEFAULT=4

class PaymentActionStatus(Enum):
    PAYMENT_ACTION_STATUS_UNKNOWN = 0
    PAYMENT_ACTION_STATUS_PENDING = 1
    PAYMENT_ACTION_STATUS_COMPLETED = 2
    PAYMENT_ACTION_STATUS_IN_DEFAULT = 3

class PaymentFrequency(Enum):
    PAYMENT_FREQUENCY_UNKNOWN = 0
    PAYMENT_FREQUENCY_WEEKLY = 1
    PAYMENT_FREQUENCY_BIWEEKLY = 2
    PAYMENT_FREQUENCY_MONTHLY = 3
    PAYMENT_FREQUENCY_QUARTERLY = 4

class PaymentAction(EmbeddedDocument):
    AccountID = StringField()
    Amount = FloatField()
    TransactionDate = DateTimeField()
    PaymentActionStatus = EnumField(PaymentActionStatus)

class PaymentPlan(Document):
    PaymentPlanID = StringField()
    UserID = StringField()
    PaymentTaskID = ListField(StringField())
    Timeline = FloatField()
    PaymentFrequency = EnumField(PaymentFrequency)
    AmountPerPayment = FloatField()
    PlanType = EnumField(PlanType)
    EndDate = DateTimeField()
    Active = BooleanField()
    Status = EnumField(PaymentStatus)
    PaymentAction = ListField(EmbeddedDocumentField(PaymentAction))


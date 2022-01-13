from mongoengine import Document, ListField, StringField, FloatField, EmbeddedDocument, EnumField, BooleanField, ObjectIdField, DateTimeField

class PlanType:
    PAYMENT_FREQUENCY_UNKNOWN = 0
    PAYMENT_FREQUENCY_WEEKLY = 1
    PAYMENT_FREQUENCY_BIWEEKLY = 2
    PAYMENT_FREQUENCY_MONTHLY = 3
    PAYMENT_FREQUENCY_QUARTERLY = 4

class PaymentStatus:
    PAYMENT_STATUS_UNKNOWN=0
    PAYMENT_STATUS_CURRENT=1
    PAYMENT_STATUS_COMPLETED=2
    PAYMENT_STATUS_IN_DEFAULT=4

class PaymentActionStatus:
    PAYMENT_ACTION_STATUS_UNKNOWN = 0
    PAYMENT_ACTION_STATUS_PENDING = 1
    PAYMENT_ACTION_STATUS_COMPLETED = 2
    PAYMENT_ACTION_STATUS_IN_DEFAULT = 3

class PaymentAction(Document):
    AccountID = StringField(required=True)  # ObjectIdField
    Amount = FloatField()
    TransactionDate = DateTimeField()
    PaymentActionStatus = EnumField(PaymentActionStatus)

class PaymentPlan(Document):
    PaymentPlanID = StringField(required=True, unique=True) # ObjectIdField
    UserID = StringField(required=True) # ObjectIdField
    PaymentTaskID = ListField(StringField(required=True))   # ListField(ObjectIdField)
    AmountPerPayment = FloatField()
    PlanType = EnumField(PlanType)
    EndDate = DateTimeField()
    Active = BooleanField()
    Status = EnumField(PaymentStatus)
    PaymentAction = ListField(EmbeddedDocument(PaymentAction))


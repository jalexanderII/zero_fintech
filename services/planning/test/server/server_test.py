from datetime import datetime

import pytest
from bson.objectid import ObjectId
from dotenv import load_dotenv
from google.protobuf.timestamp_pb2 import Timestamp
from pymongo.collection import Collection

from gen.Python.common.common_pb2 import PAYMENT_STATUS_CURRENT, PAYMENT_ACTION_STATUS_PENDING
from gen.Python.common.common_pb2 import PAYMENT_FREQUENCY_WEEKLY, PLAN_TYPE_MIN_FEES, DELETE_STATUS_SUCCESS
from gen.Python.common.payment_plan_pb2 import PaymentAction as PaymentActionPB, PaymentPlan as PaymentPlanPB, \
    DeletePaymentPlanRequest, GetPaymentPlanRequest, ListPaymentPlanRequest, ListPaymentPlanResponse, UpdatePaymentPlanRequest
from services.planning.database.database import initiate_mongo_client
from services.planning.database.models.common import PaymentFrequency as PaymentFrequencyDB
from services.planning.database.models.common import PlanType as PlanTypeDB
from services.planning.database.models.common import PaymentAction as PaymentActionDB
from services.planning.database.models.common import PaymentActionStatus as PaymentActionStatusDB
from services.planning.database.models.common import PaymentPlan as PaymentPlanDB
from services.planning.database.models.common import PaymentStatus as PaymentStatusDB
from services.planning.server.server import PlanningServicer


@pytest.fixture
def gen_server() -> PlanningServicer:
    load_dotenv()
    planningCollection: Collection = initiate_mongo_client()
    return PlanningServicer(planningCollection)


def test_create_payment_plan(gen_server):
    tt = Timestamp()
    pp = PaymentPlanPB(
        payment_plan_id=str(ObjectId()),
        user_id="61df93c0ac601d1be8e6af28",
        payment_task_id=['61dfa8296c734067e6726761', 'a2ffa82f6c734067e6726761'],
        timeline=1.0,
        payment_freq=PAYMENT_FREQUENCY_WEEKLY,
        amount_per_payment=100.0,
        plan_type=PLAN_TYPE_MIN_FEES,
        end_date=tt.GetCurrentTime(),
        active=True,
        status=PAYMENT_STATUS_CURRENT,
        payment_action=[
            PaymentActionPB(
                account_id='61df9b621d2c2b15a6e53ec9',
                amount=100.0,
                transaction_date=tt.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING
            ),
        ]
    )
    inserted = gen_server._createPaymentPlan(pp)
    assert inserted is not None, f"Failed to create a new payment plan"


def test_get_payment_plan(gen_server):
    payment_plan_id = '61e8e60a186cad9b7e6db48f'
    paymentPlanGet = gen_server.GetPaymentPlan(GetPaymentPlanRequest(payment_plan_id=payment_plan_id), context=None)
    assert payment_plan_id == paymentPlanGet.payment_plan_id,\
        f"Plan inserted ID: {payment_plan_id}\nPlan retrieved ID: {paymentPlanGet.payment_plan_id}"


def test_list_payment_plans(gen_server):
    payment_plans = gen_server.ListPaymentPlans(ListPaymentPlanRequest(), None).payment_plans
    assert len(payment_plans) > 0, f"Server did not return any PaymentPlan"


def test_update_payment_plan(gen_server):
    payment_plan_id = '61e8e60a186cad9b7e6db48f'
    tt = Timestamp()
    pp = PaymentPlanPB(
        payment_plan_id=payment_plan_id,
        user_id="61df93c0ac601d1be8e6af28",
        payment_task_id=['61dfa8296c734067e6726761', 'a2ffa82f6c734067e6726761', 'a2ffa82f6c734067e6726761'],
        timeline=3.0,
        payment_freq=PAYMENT_FREQUENCY_WEEKLY,
        amount_per_payment=150.0,
        plan_type=PLAN_TYPE_MIN_FEES,
        end_date=tt.GetCurrentTime(),
        active=True,
        status=PAYMENT_STATUS_CURRENT,
        payment_action=[
            PaymentActionPB(
                account_id='61df9b621d2c2b15a6e53ec9',
                amount=150.0,
                transaction_date=tt.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING
            ),
        ]
    )
    updated_payment_plan = gen_server.UpdatePaymentPlan(
        UpdatePaymentPlanRequest(payment_plan_id=payment_plan_id, payment_plan=pp),
        None
    )
    assert updated_payment_plan.payment_plan_id == pp.payment_plan_id, "Not equal"
    assert updated_payment_plan.timeline == pp.timeline, "Not updated"
    assert updated_payment_plan.end_date == pp.end_date, "Not updated"
    assert updated_payment_plan.amount_per_payment == pp.amount_per_payment, "Not updated"
    assert len(updated_payment_plan.payment_task_id) == len(pp.payment_task_id), "Not updated"


def test_delete_payment_plan(gen_server):
    pp = PaymentPlanDB(
        payment_plan_id=str(ObjectId()),
        user_id="to_delete",
        payment_task_id=['to_delete'],
        timeline=0,
        payment_freq=PaymentFrequencyDB.PAYMENT_FREQUENCY_UNKNOWN,
        amount_per_payment=0,
        plan_type=PlanTypeDB.PLAN_TYPE_UNKNOWN,
        end_date=datetime.now(),
        active=False,
        status=PaymentStatusDB.PAYMENT_STATUS_CANCELLED,
        payment_action=[PaymentActionDB(
            account_id='to_delete',
            amount=0,
            transaction_date=datetime.now(),
            status=PaymentActionStatusDB.PAYMENT_ACTION_STATUS_UNKNOWN
        )]
    )
    originalPaymentPlansLen = len(gen_server.ListPaymentPlans(ListPaymentPlanRequest(), None).payment_plans)
    new_payment_plan = gen_server.planning_collection.insert_one(pp.to_dict())
    new_id = new_payment_plan.inserted_id
    assert new_id is not None, f"Failed to create a new payment plan"
    updatedPaymentPlansLen = len(gen_server.ListPaymentPlans(ListPaymentPlanRequest(), None).payment_plans)
    assert updatedPaymentPlansLen == originalPaymentPlansLen + 1, f"Failed add a new payment plan"
    deleteResponse = gen_server.DeletePaymentPlan(
        request=DeletePaymentPlanRequest(payment_plan_id=str(new_id)),
        context=None,
    )
    assert deleteResponse.status == DELETE_STATUS_SUCCESS, f"Failed status is {deleteResponse.status}"

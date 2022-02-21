from datetime import datetime

import pytest
from bson.objectid import ObjectId
from dotenv import load_dotenv
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import DELETE_STATUS_SUCCESS
from gen.Python.common.common_pb2 import PAYMENT_ACTION_STATUS_PENDING
from gen.Python.common.common_pb2 import PAYMENT_FREQUENCY_WEEKLY
from gen.Python.common.common_pb2 import PAYMENT_STATUS_CURRENT
from gen.Python.common.common_pb2 import PLAN_TYPE_MIN_FEES
from gen.Python.common.payment_plan_pb2 import (
    DeletePaymentPlanRequest,
    GetPaymentPlanRequest,
    ListPaymentPlanRequest,
    UpdatePaymentPlanRequest,
)
from gen.Python.common.payment_plan_pb2 import PaymentAction as PaymentActionPB
from gen.Python.common.payment_plan_pb2 import PaymentPlan as PaymentPlanPB
from services.planning.database.db import initiate_mongo_test_client
from services.planning.database.models.common import PaymentAction as PaymentActionDB
from services.planning.database.models.common import (
    PaymentActionStatus as PaymentActionStatusDB,
)
from services.planning.database.models.common import (
    PaymentFrequency as PaymentFrequencyDB,
)
from services.planning.database.models.common import PaymentPlan as PaymentPlanDB
from services.planning.database.models.common import PaymentStatus as PaymentStatusDB
from services.planning.database.models.common import PlanType as PlanTypeDB
from services.planning.server.planning_servicer import PlanningService

load_dotenv()

tt = Timestamp()


@pytest.fixture
def gen_server() -> PlanningService:
    return PlanningService(planning_collection=initiate_mongo_test_client())


def test_save_payment_plan(gen_server):
    pp = PaymentPlanPB(
        payment_plan_id=str(ObjectId()),
        user_id="6212a0101fca9390a37a32d2",
        payment_task_id=["61dfa8296c734067e6726761", "a2ffa82f6c734067e6726761"],
        timeline=4.0,
        payment_freq=PAYMENT_FREQUENCY_WEEKLY,
        amount_per_payment=150.0,
        plan_type=PLAN_TYPE_MIN_FEES,
        end_date=tt.GetCurrentTime(),
        active=True,
        status=PAYMENT_STATUS_CURRENT,
        payment_action=[
            PaymentActionPB(
                account_id="6212a29794c88ffb3de9d764",
                amount=150.0,
                transaction_date=tt.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING,
            ),
        ],
    )
    inserted = gen_server.SavePaymentPlan(pp)
    assert inserted is not None, f"Failed to create a new payment plan"


def test_get_payment_plan(gen_server):
    payment_plan_id = "6213fd01541ef06de2168ecc"
    paymentPlanGet = gen_server.GetPaymentPlan(
        GetPaymentPlanRequest(payment_plan_id=payment_plan_id),
    )
    assert (
        payment_plan_id == paymentPlanGet.payment_plan_id
    ), f"Plan inserted ID: {payment_plan_id}\nPlan retrieved ID: {paymentPlanGet.payment_plan_id}"


def test_list_payment_plans(gen_server):
    payment_plans = gen_server.ListPaymentPlans(ListPaymentPlanRequest()).payment_plans
    assert len(payment_plans) > 0, f"Server did not return any PaymentPlan"


def test_update_payment_plan(gen_server):
    payment_plan_id = "61e8e60a186cad9b7e6db48f"
    pp = PaymentPlanPB(
        payment_plan_id=payment_plan_id,
        user_id="61df93c0ac601d1be8e6af28",
        payment_task_id=[
            "61dfa8296c734067e6726761",
            "a2ffa82f6c734067e6726761",
            "a2ffa82f6c734067e6726761",
        ],
        timeline=3.0,
        payment_freq=PAYMENT_FREQUENCY_WEEKLY,
        amount_per_payment=150.0,
        plan_type=PLAN_TYPE_MIN_FEES,
        end_date=tt.GetCurrentTime(),
        active=True,
        status=PAYMENT_STATUS_CURRENT,
        payment_action=[
            PaymentActionPB(
                account_id="61df9b621d2c2b15a6e53ec9",
                amount=150.0,
                transaction_date=tt.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING,
            ),
        ],
    )
    updated_payment_plan = gen_server.UpdatePaymentPlan(
        UpdatePaymentPlanRequest(payment_plan_id=payment_plan_id, payment_plan=pp)
    )
    assert updated_payment_plan.payment_plan_id == pp.payment_plan_id, "Not equal"
    assert updated_payment_plan.timeline == pp.timeline, "Not updated"
    assert updated_payment_plan.end_date == pp.end_date, "Not updated"
    assert (
        updated_payment_plan.amount_per_payment == pp.amount_per_payment
    ), "Not updated"
    assert len(updated_payment_plan.payment_task_id) == len(
        pp.payment_task_id
    ), "Not updated"


def test_delete_payment_plan(gen_server):
    pp = PaymentPlanDB(
        payment_plan_id=str(ObjectId()),
        user_id="to_delete",
        payment_task_id=["to_delete"],
        timeline=0,
        payment_freq=PaymentFrequencyDB.PAYMENT_FREQUENCY_UNKNOWN,
        amount_per_payment=0,
        plan_type=PlanTypeDB.PLAN_TYPE_UNKNOWN,
        end_date=datetime.now(),
        active=False,
        status=PaymentStatusDB.PAYMENT_STATUS_CANCELLED,
        payment_action=[
            PaymentActionDB(
                account_id="to_delete",
                amount=0,
                transaction_date=datetime.now(),
                status=PaymentActionStatusDB.PAYMENT_ACTION_STATUS_UNKNOWN,
            )
        ],
    )
    originalPaymentPlansLen = gen_server.planning_collection.count_documents({})
    new_payment_plan = gen_server.planning_collection.insert_one(pp.to_dict())
    new_id = new_payment_plan.inserted_id
    assert new_id is not None, f"Failed to create a new payment plan"
    updatedPaymentPlansLen = gen_server.planning_collection.count_documents({})
    assert (
        updatedPaymentPlansLen == originalPaymentPlansLen + 1
    ), f"Failed add a new payment plan"
    deleteResponse = gen_server.DeletePaymentPlan(
        request=DeletePaymentPlanRequest(payment_plan_id=str(new_id)),
    )
    assert (
        deleteResponse.status == DELETE_STATUS_SUCCESS
    ), f"Failed status is {deleteResponse.status}"

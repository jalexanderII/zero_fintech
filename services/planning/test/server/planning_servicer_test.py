from datetime import datetime, timedelta
from unittest.mock import MagicMock

import pytest
from bson.objectid import ObjectId
from dotenv import load_dotenv
from google.protobuf.timestamp_pb2 import Timestamp
from pymongo.collection import Collection
from pytest_mock import MockerFixture

from gen.Python.common.common_pb2 import DELETE_STATUS_SUCCESS, PAYMENT_ACTION_STATUS_PENDING, PLAN_TYPE_MIN_FEES, \
    PAYMENT_ACTION_STATUS_COMPLETED, PAYMENT_FREQUENCY_WEEKLY, PAYMENT_STATUS_CURRENT
from gen.Python.common.payment_plan_pb2 import (
    DeletePaymentPlanRequest,
    GetPaymentPlanRequest,
    ListPaymentPlanRequest,
    UpdatePaymentPlanRequest,
)
from gen.Python.common.payment_plan_pb2 import PaymentAction as PaymentActionPB
from gen.Python.common.payment_plan_pb2 import PaymentPlan as PaymentPlanPB
from gen.Python.planning.planning_pb2 import GetUserOverviewRequest
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
from services.planning.test.helpers.shared_objects import MOCK_USER_ID, MOCK_CHASE_ACC, MOCK_AMEX_ACC

load_dotenv()

tt = Timestamp()


@pytest.fixture
def mock_planning_server() -> PlanningService:
    return PlanningService(planning_collection=initiate_mongo_test_client())


@pytest.fixture
def patch__fetch_accounts(mocker: MockerFixture) -> None:
    mocker.patch.object(PlanningService, "_fetch_accounts").return_value = [MOCK_CHASE_ACC, MOCK_AMEX_ACC]


@pytest.fixture
def patch_planning_collection_find(mocker: MockerFixture) -> None:
    last_month, this_month, next_month = datetime.now()+timedelta(days=-31), datetime.now(),\
                                         datetime.now()+timedelta(days=31)
    mocker.patch.object(Collection, 'find').return_value = [
        {'_id': ObjectId('622874c7e9f2ee76c713c40e'), 'paymentPlanId': '622874c7e9f2ee76c713c40d', 'userId': 'test',
         'paymentTaskId': ['61dfa8296c734067e6726761', 'a2ffa82f6c734067e6726761'], 'amount': 300.0, 'timeline': 1.0,
         'paymentFreq': 'PAYMENT_FREQUENCY_WEEKLY', 'amountPerPayment': 300.0, 'planType': 'PLAN_TYPE_MIN_FEES',
         'endDate': str(this_month), 'active': True, 'status': 'PAYMENT_STATUS_CURRENT', 'paymentAction': [
            {'accountId': '1', 'amount': 150.0, 'transactionDate': str(this_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '2', 'amount': 150.0, 'transactionDate': str(this_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'}]},
        {'_id': ObjectId('622876a215c75f609b0d65ee'), 'paymentPlanId': '622876a215c75f609b0d65ed', 'userId': 'test',
         'paymentTaskId': ['01', '02'], 'amount': 450.0, 'timeline': 3.0, 'paymentFreq': 'PAYMENT_FREQUENCY_WEEKLY',
         'amountPerPayment': 150.0, 'planType': 'PLAN_TYPE_MIN_FEES', 'endDate': str(next_month),
         'active': True, 'status': 'PAYMENT_STATUS_CURRENT', 'paymentAction': [
            {'accountId': '1', 'amount': 150.0, 'transactionDate': str(last_month),
             'status': 'PAYMENT_ACTION_STATUS_COMPLETED'},
            {'accountId': '1', 'amount': 150.0, 'transactionDate': str(this_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '2', 'amount': 150.0, 'transactionDate': str(next_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'}]},
        {'_id': ObjectId('6228a94425de4b0368862441'), 'paymentPlanId': '6228a94425de4b0368862440', 'userId': 'test',
         'paymentTaskId': ['61dfa8296c734067e6726761', 'a2ffa82f6c734067e6726761'], 'amount': 300.0, 'timeline': 1.0,
         'paymentFreq': 'PAYMENT_FREQUENCY_WEEKLY', 'amountPerPayment': 300.0, 'planType': 'PLAN_TYPE_MIN_FEES',
         'endDate': str(this_month), 'active': True, 'status': 'PAYMENT_STATUS_CURRENT', 'paymentAction': [
            {'accountId': '1', 'amount': 150.0, 'transactionDate': str(this_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '2', 'amount': 150.0, 'transactionDate': str(this_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'}]},
        {'_id': ObjectId('623882d0da531f4cb175600b'), 'paymentPlanId': '623882d0da531f4cb1756009', 'userId': 'test',
         'paymentTaskId': ['01', '02'], 'amount': 450.0, 'timeline': 3.0, 'paymentFreq': 'PAYMENT_FREQUENCY_WEEKLY',
         'amountPerPayment': 150.0, 'planType': 'PLAN_TYPE_MIN_FEES', 'endDate': str(next_month),
         'active': True, 'status': 'PAYMENT_STATUS_CURRENT', 'paymentAction': [
            {'accountId': '1', 'amount': 150.0, 'transactionDate': str(last_month),
             'status': 'PAYMENT_ACTION_STATUS_COMPLETED'},
            {'accountId': '1', 'amount': 150.0, 'transactionDate': str(this_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '2', 'amount': 150.0, 'transactionDate': str(next_month),
             'status': 'PAYMENT_ACTION_STATUS_PENDING'}]}
    ]


def test_save_payment_plan(mock_planning_server: PlanningService):
    pp = PaymentPlanPB(
        payment_plan_id=str(ObjectId()),
        user_id=MOCK_USER_ID,
        payment_task_id=["61dfa8296c734067e6726761", "a2ffa82f6c734067e6726761"],
        amount=300.0,
        timeline=1.0,
        payment_freq=PAYMENT_FREQUENCY_WEEKLY,
        amount_per_payment=300.0,
        plan_type=PLAN_TYPE_MIN_FEES,
        end_date=tt.GetCurrentTime(),
        active=True,
        status=PAYMENT_STATUS_CURRENT,
        payment_action=[
            PaymentActionPB(
                account_id=MOCK_CHASE_ACC.account_id,
                amount=150.0,
                transaction_date=tt.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING,
            ),
            PaymentActionPB(
                account_id=MOCK_AMEX_ACC.account_id,
                amount=150.0,
                transaction_date=tt.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING,
            ),
        ],
    )
    inserted = mock_planning_server.SavePaymentPlan(pp)
    assert inserted is not None, f"Failed to create a new payment plan"


def test_get_payment_plan(mock_planning_server: PlanningService):
    payment_plan_id = "622874c7e9f2ee76c713c40e"
    paymentPlanGet = mock_planning_server.GetPaymentPlan(
        GetPaymentPlanRequest(payment_plan_id=payment_plan_id),
    )
    assert (
        payment_plan_id == paymentPlanGet.payment_plan_id
    ), f"Plan inserted ID: {payment_plan_id}\nPlan retrieved ID: {paymentPlanGet.payment_plan_id}"


def test_list_payment_plans(mock_planning_server: PlanningService):
    payment_plans = mock_planning_server.ListPaymentPlans(ListPaymentPlanRequest()).payment_plans
    assert len(payment_plans) > 0, f"Server did not return any PaymentPlan"


def test_update_payment_plan(mock_planning_server: PlanningService):
    payment_plan_id = "6213fd01541ef06de2168ecc"
    pp = PaymentPlanPB(
        payment_plan_id=payment_plan_id,
        user_id="61df93c0ac601d1be8e6af28",
        payment_task_id=[
            "61dfa8296c734067e6726761",
            "a2ffa82f6c734067e6726761",
            "a2ffa82f6c734067e6726761",
        ],
        amount=150.0,
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
    updated_payment_plan = mock_planning_server.UpdatePaymentPlan(
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


def test_delete_payment_plan(mock_planning_server: PlanningService):
    pp = PaymentPlanDB(
        payment_plan_id=str(ObjectId()),
        user_id="to_delete",
        payment_task_id=["to_delete"],
        amount=0.0,
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
    originalPaymentPlansLen = mock_planning_server.planning_collection.count_documents({})
    new_payment_plan = mock_planning_server.planning_collection.insert_one(pp.to_dict())
    new_id = new_payment_plan.inserted_id
    assert new_id is not None, f"Failed to create a new payment plan"
    updatedPaymentPlansLen = mock_planning_server.planning_collection.count_documents({})
    assert (
        updatedPaymentPlansLen == originalPaymentPlansLen + 1
    ), f"Failed add a new payment plan"
    deleteResponse = mock_planning_server.DeletePaymentPlan(
        request=DeletePaymentPlanRequest(payment_plan_id=str(new_id)),
    )
    assert (
        deleteResponse.status == DELETE_STATUS_SUCCESS
    ), f"Failed status is {deleteResponse.status}"


def test_get_amount_paid_percentage(patch__fetch_accounts: None, patch_planning_collection_find: None,
                                    mock_planning_server: PlanningService):
    amount_paid_percentage_response = mock_planning_server.GetAmountPaidPercentage(
        request=GetUserOverviewRequest(user_id=MOCK_USER_ID)
    )

    assert amount_paid_percentage_response.percentage_paid == 0.2


def test_get_percentage_covered_by_plans(patch__fetch_accounts: None, patch_planning_collection_find: None,
                                         mock_planning_server: PlanningService):
    percentage_covered_response = mock_planning_server.GetPercentageCoveredByPlans(
        request=GetUserOverviewRequest(user_id=MOCK_USER_ID)
    )

    assert percentage_covered_response.overall_covered == 2/3
    assert percentage_covered_response.account_to_percent_covered["1"] == 0.75
    assert percentage_covered_response.account_to_percent_covered["2"] == 0.6


def test_get_waterfall_overview(patch__fetch_accounts: None, patch_planning_collection_find: None,
                                mock_planning_server: PlanningService):
    waterfall_overview_response = mock_planning_server.GetWaterfallOverview(
        request=GetUserOverviewRequest(user_id=MOCK_USER_ID)
    ).monthly_waterfall

    assert waterfall_overview_response[0].account_to_amounts["1"] == 600
    assert waterfall_overview_response[0].account_to_amounts["2"] == 300
    # adding 31 days isn't precisely a month such that either next or the month thereafter will a PaymentAction
    assert waterfall_overview_response[1].account_to_amounts["2"] == 300 or \
           waterfall_overview_response[2].account_to_amounts["2"] == 300



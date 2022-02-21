import datetime
from unittest.mock import MagicMock

import pandas as pd
import pytest
from google.protobuf.timestamp_pb2 import Timestamp
from pytest_mock import MockerFixture

from gen.Python.common.common_pb2 import (
    PAYMENT_FREQUENCY_MONTHLY,
    PAYMENT_FREQUENCY_BIWEEKLY,
    PAYMENT_ACTION_STATUS_PENDING,
    PAYMENT_FREQUENCY_QUARTERLY,
)
from gen.Python.common.common_pb2 import (
    PAYMENT_STATUS_CURRENT,
    PLAN_TYPE_OPTIM_CREDIT_SCORE,
)
from gen.Python.common.common_pb2 import PLAN_TYPE_MIN_FEES
from gen.Python.common.payment_plan_pb2 import PaymentAction, PaymentPlan
from gen.Python.core.accounts_pb2 import Account, AnnualPercentageRates
from services.planning.server.payment_plan_builder import PaymentPlanBuilder
from services.planning.test.helpers.paramset import (
    MetaDataToPaymentPlanParams,
    CreatePaymentActionsParams,
)

START_DATE = datetime.datetime(2022, 2, 21, 17, 26, 12)

MOCK_CHASE_ACC = Account(
    account_id="1",
    user_id="test",
    name="Chase",
    available_balance=2500,
    current_balance=500,
    credit_limit=3000,
    annual_percentage_rate=[
        AnnualPercentageRates(apr_percentage=22, apr_type="purchase_apr")
    ],
)

MOCK_AMEX_ACC = Account(
    account_id="2",
    user_id="test",
    name="Amex",
    available_balance=4000,
    current_balance=1000,
    credit_limit=5000,
    annual_percentage_rate=[
        AnnualPercentageRates(apr_percentage=42, apr_type="purchase_apr")
    ],
)

MOCK_DF = pd.DataFrame(
    {
        "account_id": [MOCK_CHASE_ACC.account_id, MOCK_AMEX_ACC.account_id],
        "apr": [
            MOCK_CHASE_ACC.annual_percentage_rate[0].apr_percentage,
            MOCK_AMEX_ACC.annual_percentage_rate[0].apr_percentage,
        ],
        "amount": [MOCK_CHASE_ACC.current_balance, MOCK_AMEX_ACC.current_balance],
        "payment_task_id": ["01", "02"],
        "balance": [MOCK_CHASE_ACC.current_balance, MOCK_AMEX_ACC.current_balance],
        "credit_limit": [MOCK_CHASE_ACC.credit_limit, MOCK_AMEX_ACC.credit_limit],
    }
)

MOCK_MIN_FEE_MONTHLY_ACTIONS = [
    PaymentAction(
        account_id="1",
        amount=500,
        transaction_date=Timestamp(seconds=1647820800),
        status=PAYMENT_ACTION_STATUS_PENDING,
    ),
    PaymentAction(
        account_id="2",
        amount=500,
        transaction_date=Timestamp(seconds=1650499200),
        status=PAYMENT_ACTION_STATUS_PENDING,
    ),
    PaymentAction(
        account_id="2",
        amount=500,
        transaction_date=Timestamp(seconds=1653091200),
        status=PAYMENT_ACTION_STATUS_PENDING,
    ),
]
MOCK_OPTIM_MONTHLY_ACTIONS = [
    PaymentAction(
        account_id="2",
        amount=500,
        transaction_date=Timestamp(seconds=1647820800),
        status=PAYMENT_ACTION_STATUS_PENDING,
    ),
    PaymentAction(
        account_id="1",
        amount=500,
        transaction_date=Timestamp(seconds=1650499200),
        status=PAYMENT_ACTION_STATUS_PENDING,
    ),
    PaymentAction(
        account_id="2",
        amount=500,
        transaction_date=Timestamp(seconds=1653091200),
        status=PAYMENT_ACTION_STATUS_PENDING,
    ),
]


@pytest.fixture
def mock_payment_plan_builder() -> PaymentPlanBuilder:
    ppb = PaymentPlanBuilder()
    return ppb


@pytest.fixture
def patch__fetch_accounts(mocker: MockerFixture):
    return mocker.patch.object(PaymentPlanBuilder, "_fetch_accounts")


@pytest.mark.parametrize(
    "p",
    [
        MetaDataToPaymentPlanParams(
            id="Test Payment Plan for Monthly min fee",
            user_id="test",
            plan_type=PLAN_TYPE_MIN_FEES,
            timeline_months=3,
            payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            payment_task_ids=["01", "02"],
            account_ids=[MOCK_CHASE_ACC.account_id, MOCK_AMEX_ACC.account_id],
            amounts=[MOCK_CHASE_ACC.current_balance, MOCK_AMEX_ACC.current_balance],
            expected=PaymentPlan(
                user_id="test",
                payment_task_id=["01", "02"],
                timeline=3,
                payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                amount_per_payment=500,
                plan_type=PLAN_TYPE_MIN_FEES,
                end_date=Timestamp(seconds=1653091200),
                active=True,
                status=PAYMENT_STATUS_CURRENT,
                payment_action=MOCK_MIN_FEE_MONTHLY_ACTIONS,
            ),
        ),
        MetaDataToPaymentPlanParams(
            id="Test Payment Plan for Monthly optim credit",
            user_id="test",
            plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
            timeline_months=3,
            payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            payment_task_ids=["01", "02"],
            account_ids=[MOCK_CHASE_ACC.account_id, MOCK_AMEX_ACC.account_id],
            amounts=[MOCK_CHASE_ACC.current_balance, MOCK_AMEX_ACC.current_balance],
            expected=PaymentPlan(
                user_id="test",
                payment_task_id=["01", "02"],
                timeline=3,
                payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                amount_per_payment=500,
                plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                end_date=Timestamp(seconds=1653091200),
                active=True,
                status=PAYMENT_STATUS_CURRENT,
                payment_action=MOCK_OPTIM_MONTHLY_ACTIONS,
            ),
        ),
    ],
)
def test__create_from_meta_data(
    p: MetaDataToPaymentPlanParams,
    patch__fetch_accounts: MagicMock,
    mock_payment_plan_builder: PaymentPlanBuilder,
):
    ppb = mock_payment_plan_builder
    patch__fetch_accounts.return_value = [MOCK_CHASE_ACC, MOCK_AMEX_ACC]
    args = (
        p.user_id,
        p.plan_type,
        p.timeline_months,
        p.payment_freq,
        p.payment_task_ids,
        p.account_ids,
        p.amounts,
    )
    actual = ppb._create_from_meta_data(*args)
    assert actual == p.expected


@pytest.mark.parametrize(
    "p",
    [
        CreatePaymentActionsParams(
            id="Test Monthly Plan min fee",
            payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            df=MOCK_DF,
            start_date=START_DATE,
            amount_per_payment=500.0,
            expected=MOCK_MIN_FEE_MONTHLY_ACTIONS,
            min_fee=True,
        ),
        CreatePaymentActionsParams(
            id="Test BiWeekly Plan min fee",
            payment_freq=PAYMENT_FREQUENCY_BIWEEKLY,
            df=MOCK_DF,
            start_date=START_DATE,
            amount_per_payment=250.0,
            expected=[
                PaymentAction(
                    account_id="1",
                    amount=250,
                    transaction_date=Timestamp(seconds=1646673972),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
                PaymentAction(
                    account_id="1",
                    amount=250,
                    transaction_date=Timestamp(seconds=1647883572),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
                PaymentAction(
                    account_id="2",
                    amount=250,
                    transaction_date=Timestamp(seconds=1649093172),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
                PaymentAction(
                    account_id="2",
                    amount=250,
                    transaction_date=Timestamp(seconds=1650302772),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
                PaymentAction(
                    account_id="2",
                    amount=250,
                    transaction_date=Timestamp(seconds=1651512372),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
                PaymentAction(
                    account_id="2",
                    amount=250,
                    transaction_date=Timestamp(seconds=1652721972),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
            ],
            min_fee=True,
        ),
        CreatePaymentActionsParams(
            id="Test Monthly Plan credit",
            payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            df=MOCK_DF,
            start_date=START_DATE,
            amount_per_payment=500.0,
            expected=MOCK_OPTIM_MONTHLY_ACTIONS,
            min_fee=False,
        ),
        CreatePaymentActionsParams(
            id="Test Quarterly Plan credit",
            payment_freq=PAYMENT_FREQUENCY_QUARTERLY,
            df=MOCK_DF,
            start_date=START_DATE,
            amount_per_payment=750,
            expected=[
                PaymentAction(
                    account_id="2",
                    amount=750,
                    transaction_date=Timestamp(seconds=1653091200),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
                PaymentAction(
                    account_id="1",
                    amount=500,
                    transaction_date=Timestamp(seconds=1661040000),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
                PaymentAction(
                    account_id="2",
                    amount=250,
                    transaction_date=Timestamp(seconds=1661040000),
                    status=PAYMENT_ACTION_STATUS_PENDING,
                ),
            ],
            min_fee=False,
        ),
    ],
)
def test__create_payment_actions_min_fees(
    p: CreatePaymentActionsParams, mock_payment_plan_builder: PaymentPlanBuilder
):
    ppb = mock_payment_plan_builder
    if p.min_fee:
        actual = ppb._create_payment_actions_min_fees(
            p.payment_freq, p.df, p.start_date, p.amount_per_payment
        )
    else:
        actual = ppb._create_payment_actions_optim_credit_score(
            p.payment_freq, p.df, p.start_date, p.amount_per_payment
        )
    assert actual == p.expected


# TODO(Joschka)
def test__get_meta_data_options():
    pass

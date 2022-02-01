import datetime
from typing import List

from google.protobuf.timestamp_pb2 import Timestamp
from pytest_cases import parametrize_with_cases, case, fixture as cases_fixture

from gen.Python.common.common_pb2 import PLAN_TYPE_MIN_FEES, PAYMENT_FREQUENCY_MONTHLY, \
    PAYMENT_FREQUENCY_UNKNOWN, PLAN_TYPE_OPTIM_CREDIT_SCORE, PLAN_TYPE_UNKNOWN, PAYMENT_ACTION_STATUS_PENDING, \
    PAYMENT_FREQUENCY_BIWEEKLY, PAYMENT_FREQUENCY_QUARTERLY, PAYMENT_STATUS_CURRENT
from gen.Python.common.common_pb2 import PaymentFrequency
from gen.Python.common.common_pb2 import PaymentStatus
from gen.Python.common.payment_plan_pb2 import PaymentPlan, PaymentAction
from gen.Python.common.payment_task_pb2 import PaymentTask, MetaData
from services.planning.server.payment_plan_builder import PaymentPlanBuilder
from services.planning.server.utils import datetime_to_pb_timestamp
from services.planning.server.utils import shift_date_by_payment_frequency


@cases_fixture
def gen_payment_plan_builder() -> PaymentPlanBuilder:
    return PaymentPlanBuilder()


def shift_now_by_payment_frequency_multiple_times(
        paymentFreq: PaymentFrequency, howOften: int
) -> Timestamp:
    """Helper function to shift current date/now by PaymentFrequency multiple times."""
    date = datetime.datetime.now()
    for _ in range(howOften):
        date = shift_date_by_payment_frequency(date=date, payment_freq=paymentFreq)
    return datetime_to_pb_timestamp(date)


class Cases:
    user_id = "61df93c0ac601d1be8e64613"
    accName2Id = {
        "Amex": "61df9b621d2c2b15a6e53ec9",
        "Chase": "61df9af7f18b94fc44d09fb9",
        "Barclay": "61df9f3397fa9f3b7a9b67a8",
    }

    # APR: 0.10, 0.1975, 0.34
    # BAL: 10000, 9000, 345.23
    # CLI: 15000,  25000, , 5000
    # %:    0.67, 0.36, , 0.07

    @case(tags="end_to_end")
    def case_min_fees_2_month_monthly(self):
        return (  # first tuple
            [
                PaymentTask(
                    user_id=self.user_id, account_id=self.accName2Id["Amex"], amount=500
                ),  # Amex
                PaymentTask(
                    user_id=self.user_id,
                    account_id=self.accName2Id["Chase"],
                    amount=500,
                ),  # Chase
            ],
            MetaData(
                preferred_plan_type=PLAN_TYPE_MIN_FEES,
                preferred_timeline_in_months=2.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            ),
            [
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=["", ""],
                    timeline=2.0,
                    payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=500,
                    plan_type=PLAN_TYPE_MIN_FEES,
                    end_date=shift_now_by_payment_frequency_multiple_times(
                        PAYMENT_FREQUENCY_MONTHLY, 2
                    ),
                    active=True,
                    status=PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=500,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 2
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=500,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 1
                            ),
                        ),
                    ],
                )
            ],
        )

    @case(tags="end_to_end")
    def case_credit_score_2_month_monthly(self):
        return (  # first tuple
            [
                PaymentTask(
                    user_id=self.user_id, account_id=self.accName2Id["Amex"], amount=500
                ),  # Amex
                PaymentTask(
                    user_id=self.user_id,
                    account_id=self.accName2Id["Chase"],
                    amount=500,
                ),  # Chase
            ],
            MetaData(
                preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                preferred_timeline_in_months=2.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            ),
            [
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=["", ""],
                    timeline=2.0,
                    payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=500,
                    plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(
                        PAYMENT_FREQUENCY_MONTHLY, 2
                    ),
                    active=True,
                    status=PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=500,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 1
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=500,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 2
                            ),
                        ),
                    ],
                )
            ],
        )

    @case(tags="end_to_end")
    def case_credit_score_no_other_meta_data(self):
        return (  # first tuple
            [
                PaymentTask(
                    user_id=self.user_id, account_id=self.accName2Id["Amex"], amount=500
                ),  # Amex
                PaymentTask(
                    user_id=self.user_id,
                    account_id=self.accName2Id["Chase"],
                    amount=500,
                ),  # Chase
            ],
            MetaData(
                preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                preferred_timeline_in_months=0.0,
                preferred_payment_freq=None,
            ),
            [
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=["", ""],
                    timeline=3.0,
                    payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=333.34,
                    plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(
                        PAYMENT_FREQUENCY_MONTHLY, 3
                    ),
                    active=True,
                    status=PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=333.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 1
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=166.66,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 2
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=166.68,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 2
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=333.32,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 3
                            ),
                        ),
                    ],
                ),
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=["", ""],
                    timeline=6.0,
                    payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=166.67,
                    plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(
                        PAYMENT_FREQUENCY_MONTHLY, 6
                    ),
                    active=True,
                    status=PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=166.67,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 1
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=166.67,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 2
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=166.66,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 3
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=0.01,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 3
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=166.67,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 4
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=166.67,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 5
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=166.65,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 6
                            ),
                        ),
                    ],
                ),
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=["", ""],
                    timeline=12.0,
                    payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=83.34,
                    plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(
                        PAYMENT_FREQUENCY_MONTHLY, 12
                    ),
                    active=True,
                    status=PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 1
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 2
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 3
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 4
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 5
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Amex"],
                            amount=83.30,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 6
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=0.04,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 6
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 7
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 8
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 9
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 10
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=83.34,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 11
                            ),
                        ),
                        PaymentAction(
                            account_id=self.accName2Id["Chase"],
                            amount=83.26,
                            status=PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PAYMENT_FREQUENCY_MONTHLY, 12
                            ),
                        ),
                    ],
                ),
            ],
        )

    @case(tags="create_single_payment_plan")
    def case_min_fees_2_months_biweekly(self):
        return (
            self.user_id,
            MetaData(
                preferred_plan_type=PLAN_TYPE_MIN_FEES,
                preferred_timeline_in_months=2.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_BIWEEKLY,
            ),
            ["", ""],
            [self.accName2Id["Chase"], self.accName2Id["Amex"]],
            [300, 200],
            PaymentPlan(
                user_id=self.user_id,
                payment_task_id=["", ""],
                timeline=2.0,
                payment_freq=PAYMENT_FREQUENCY_BIWEEKLY,
                amount_per_payment=125,
                plan_type=PLAN_TYPE_MIN_FEES,
                end_date=shift_now_by_payment_frequency_multiple_times(
                    PAYMENT_FREQUENCY_BIWEEKLY, 4
                ),
                active=True,
                status=PAYMENT_STATUS_CURRENT,
                payment_action=[
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=125,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_BIWEEKLY, 1
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=125,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_BIWEEKLY, 2
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=50,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_BIWEEKLY, 3
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Amex"],
                        amount=75,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_BIWEEKLY, 3
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Amex"],
                        amount=125,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_BIWEEKLY, 4
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                ],
            ),
        )

    @case(tags="create_single_payment_plan")
    def case_credit_score_5_months_monthly(self):
        return (
            self.user_id,
            MetaData(
                preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                preferred_timeline_in_months=5.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            ),
            ["", "", ""],
            [
                self.accName2Id["Chase"],
                self.accName2Id["Amex"],
                self.accName2Id["Barclay"],
            ],
            [180, 170, 170],
            PaymentPlan(
                user_id=self.user_id,
                payment_task_id=["", "", ""],
                timeline=5.0,
                payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                amount_per_payment=104.0,
                plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                end_date=shift_now_by_payment_frequency_multiple_times(
                    PAYMENT_FREQUENCY_MONTHLY, 5
                ),
                active=True,
                status=PAYMENT_STATUS_CURRENT,
                payment_action=[
                    PaymentAction(
                        account_id=self.accName2Id["Amex"],
                        amount=104,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 1
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Amex"],
                        amount=66,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 2
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=38,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 2
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=104,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 3
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=38,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 4
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Barclay"],
                        amount=66,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 4
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Barclay"],
                        amount=104,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 5
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                ],
            ),
        )

    @case(tags="create_single_payment_plan")
    def case_credit_score_3_months_monthly_pay_off_whole_acc(self):
        return (
            self.user_id,
            MetaData(
                preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                preferred_timeline_in_months=3.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            ),
            ["", ""],
            [self.accName2Id["Chase"], self.accName2Id["Amex"]],
            [9000, 10000],
            PaymentPlan(
                user_id=self.user_id,
                payment_task_id=["", ""],
                timeline=3.0,
                payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                amount_per_payment=6333.34,
                plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                end_date=shift_now_by_payment_frequency_multiple_times(
                    PAYMENT_FREQUENCY_MONTHLY, 3
                ),
                active=True,
                status=PAYMENT_STATUS_CURRENT,
                payment_action=[
                    PaymentAction(
                        account_id=self.accName2Id["Amex"],
                        amount=6333.34,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 1
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=6333.34,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 2
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Amex"],
                        amount=3666.66,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 3
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                    PaymentAction(
                        account_id=self.accName2Id["Chase"],
                        amount=2666.66,
                        transaction_date=shift_now_by_payment_frequency_multiple_times(
                            PAYMENT_FREQUENCY_MONTHLY, 3
                        ),
                        status=PAYMENT_ACTION_STATUS_PENDING,
                    ),
                ],
            ),
        )

    @case(tags="meta_data_options")
    def case_no_meta_data_amount_lt_250(self):
        return (
            MetaData(
                preferred_plan_type=PLAN_TYPE_UNKNOWN,
                preferred_timeline_in_months=0.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_UNKNOWN,
            ),
            240,
            [
                MetaData(
                    preferred_plan_type=PLAN_TYPE_MIN_FEES,
                    preferred_timeline_in_months=1.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                ),
                MetaData(
                    preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    preferred_timeline_in_months=1.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                ),
            ],
        )

    @case(tags="meta_data_options")
    def case_min_fees_no_other_meta_data_amount_gt_250(self):
        return (
            MetaData(
                preferred_plan_type=PLAN_TYPE_MIN_FEES,
                preferred_timeline_in_months=0.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_UNKNOWN,
            ),
            500,
            [
                MetaData(
                    preferred_plan_type=PLAN_TYPE_MIN_FEES,
                    preferred_timeline_in_months=3.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                ),
                MetaData(
                    preferred_plan_type=PLAN_TYPE_MIN_FEES,
                    preferred_timeline_in_months=6.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                ),
                MetaData(
                    preferred_plan_type=PLAN_TYPE_MIN_FEES,
                    preferred_timeline_in_months=12.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_MONTHLY,
                ),
            ],
        )

    @case(tags="meta_data_options")
    def case_credit_score_biweekly_payments_no_timeline_amount_gt_250(self):
        return (
            MetaData(
                preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                preferred_timeline_in_months=0.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_BIWEEKLY,
            ),
            500,
            [
                MetaData(
                    preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    preferred_timeline_in_months=1.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_BIWEEKLY,
                ),
                MetaData(
                    preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    preferred_timeline_in_months=2.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_BIWEEKLY,
                ),
            ],
        )

    @case(tags="meta_data_options")
    def case_quarterly_6_month_no_type(self):
        return (
            MetaData(
                preferred_plan_type=PLAN_TYPE_UNKNOWN,
                preferred_timeline_in_months=6.0,
                preferred_payment_freq=PAYMENT_FREQUENCY_QUARTERLY,
            ),
            500,
            [
                MetaData(
                    preferred_plan_type=PLAN_TYPE_MIN_FEES,
                    preferred_timeline_in_months=6.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_QUARTERLY,
                ),
                MetaData(
                    preferred_plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    preferred_timeline_in_months=6.0,
                    preferred_payment_freq=PAYMENT_FREQUENCY_QUARTERLY,
                ),
            ],
        )


@parametrize_with_cases(
    "paymentTasks, metaData, paymentPlans", cases=Cases, has_tag="end_to_end"
)
def test_create_payment_plan_end_to_end(
        paymentTasks, metaData, paymentPlans, gen_payment_plan_builder
):
    paymentPlansCreated = gen_payment_plan_builder.create(
        paymentTasks, metaData
    )

    assert len(paymentPlans) == len(paymentPlansCreated), (
        f"Tests asks for {len(paymentPlans)} but"
        f"{len(paymentPlansCreated)} were created"
    )

    for i, (paymentPlan, paymentPlanCreated) in enumerate(
            zip(paymentPlans, paymentPlansCreated)
    ):
        assert paymentPlan.user_id == paymentPlanCreated.user_id, f"Plan {i} differs"
        assert sorted(paymentPlan.payment_task_id) == sorted(
            paymentPlanCreated.payment_task_id
        ), f"Plan {i} differs"
        assert paymentPlan.timeline == paymentPlanCreated.timeline, f"Plan {i} differs"
        assert (
                paymentPlan.payment_freq == paymentPlanCreated.payment_freq
        ), f"Plan {i} differs"
        assert (
                paymentPlan.plan_type == paymentPlanCreated.plan_type
        ), f"Plan {i} differs"
        assert (
                abs(paymentPlan.amount_per_payment - paymentPlanCreated.amount_per_payment)
                < 1e-3
        ), f"Plan {i} differs"
        assert (
                abs(paymentPlan.end_date.seconds - paymentPlanCreated.end_date.seconds) < 2
        ), f"Plan {i} differs"
        assert paymentPlan.active == paymentPlanCreated.active, f"Plan {i} differs"
        assert paymentPlan.status == paymentPlanCreated.status, f"Plan {i} differs"
        assert len(paymentPlan.payment_action) == len(
            paymentPlanCreated.payment_action
        ), f"Plan {i} differs"
        for ii, (paymentAction, paymentActionCreated) in enumerate(
                zip(paymentPlan.payment_action, paymentPlanCreated.payment_action)
        ):
            assert (
                    paymentAction.account_id == paymentActionCreated.account_id
            ), f"Action {ii} of plan {i} differs"
            assert (
                    abs(paymentAction.amount - paymentActionCreated.amount) < 1e-3
            ), f"Action {ii} of plan {i} differs"
            assert (
                    paymentAction.status == paymentAction.status
            ), f"Action {ii} of plan {i} differs"


@parametrize_with_cases(
    "userId, metaData, paymentTaskIds, accountIds, amounts, paymentPlan",
    cases=Cases,
    has_tag="create_single_payment_plan",
)
def test_create_single_payment_plan(
        userId: str,
        metaData: MetaData,
        paymentTaskIds: List[str],
        accountIds: List[str],
        amounts: List[float],
        paymentPlan: PaymentPlan,
        gen_payment_plan_builder,
):
    paymentPlanCreated = gen_payment_plan_builder._create_from_meta_data(
        user_id=userId,
        plan_type=metaData.preferred_plan_type,
        timeline_months=metaData.preferred_timeline_in_months,
        payment_freq=metaData.preferred_payment_freq,
        payment_task_ids=paymentTaskIds,
        account_ids=accountIds,
        amounts=amounts,
    )

    assert paymentPlan.user_id == paymentPlanCreated.user_id
    assert sorted(paymentPlan.payment_task_id) == sorted(
        paymentPlanCreated.payment_task_id
    )
    assert paymentPlan.timeline == paymentPlanCreated.timeline
    assert paymentPlan.payment_freq == paymentPlanCreated.payment_freq
    assert paymentPlan.plan_type == paymentPlanCreated.plan_type
    print(paymentPlan.amount_per_payment - paymentPlanCreated.amount_per_payment)
    assert (
            abs(paymentPlan.amount_per_payment - paymentPlanCreated.amount_per_payment)
            < 1e-3
    )
    assert abs(paymentPlan.end_date.seconds - paymentPlanCreated.end_date.seconds) < 2
    assert paymentPlan.active == paymentPlanCreated.active
    assert paymentPlan.status == paymentPlanCreated.status
    assert len(paymentPlan.payment_action) == len(paymentPlanCreated.payment_action)
    for ii, (paymentAction, paymentActionCreated) in enumerate(
            zip(paymentPlan.payment_action, paymentPlanCreated.payment_action)
    ):
        assert (
                paymentAction.account_id == paymentActionCreated.account_id
        ), f"Action {ii} differs"
        assert (
                abs(paymentAction.amount - paymentActionCreated.amount) < 1e-3
        ), f"Action {ii} differs"
        assert (
                abs(
                    paymentAction.transaction_date.seconds
                    - paymentActionCreated.transaction_date.seconds
                )
                < 2
        ), f"Action {ii} differs"
        assert paymentAction.status == paymentAction.status, f"Action {ii} differs"


@parametrize_with_cases(
    "metaData, totalAmount, metaDataOptions", cases=Cases, has_tag="meta_data_options"
)
def test_create_meta_data_options(
        metaData: MetaData,
        totalAmount: float,
        metaDataOptions: List[MetaData],
        gen_payment_plan_builder,
):
    metaDataOptionsCreated = gen_payment_plan_builder._get_meta_data_options(
        metaData, totalAmount > 250.0
    )

    assert len(metaDataOptions) == len(metaDataOptionsCreated), (
        f"Test asks for {len(metaDataOptions)} but "
        f"{len(metaDataOptionsCreated)} were created"
    )
    for i, (metaData, metaDataCreated) in enumerate(
            zip(metaDataOptions, metaDataOptionsCreated)
    ):
        assert (
                metaData.preferred_plan_type == metaDataCreated.preferred_plan_type
        ), f"MetaData {i} differs"
        assert (
                metaData.preferred_timeline_in_months
                == metaDataCreated.preferred_timeline_in_months
        ), f"MetaData {i} differs"
        assert (
                metaData.preferred_payment_freq == metaDataCreated.preferred_payment_freq
        ), f"MetaData {i} differs"

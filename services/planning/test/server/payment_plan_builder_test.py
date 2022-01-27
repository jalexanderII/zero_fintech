import datetime

from pytest_cases import parametrize_with_cases, fixture as cases_fixture
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PlanType, PaymentFrequency, PaymentActionStatus, PaymentStatus
from gen.Python.common.payment_plan_pb2 import PaymentPlan, PaymentAction
from gen.Python.common.payment_task_pb2 import PaymentTask, MetaData
from services.planning.server.payment_plan_builder import PaymentPlanBuilder, payment_plan_builder
from services.planning.server.utils import shift_date_by_payment_frequency

@cases_fixture
def gen_payment_plan_builder() -> PaymentPlanBuilder:
    return payment_plan_builder


def datetime2timestamp(date: datetime) -> Timestamp:
    """ Helper function for pytest.mark.paramterize"""
    timestamp = Timestamp()
    timestamp.FromDatetime(date)
    return timestamp


def shift_now_by_payment_frequency_multiple_times(paymentFreq: PaymentFrequency, howOften: int) -> Timestamp:
    """ Helper function to shift current date/now by PaymentFrequency multiple times. """
    date = datetime.datetime.now()
    for _ in range(howOften):
        date = shift_date_by_payment_frequency(date=date, payment_freq=paymentFreq)
    return datetime2timestamp(date)


class Cases:
    def __init__(self):
        self.user_id = '61df93c0ac601d1be8e64613'
        self.accName2Id = {'Amex': '61df9b621d2c2b15a6e53ec9', 'Chase': '61df9af7f18b94fc44d09fb9'}

    def case_min_fees_2_month_monthly(self):
        return (  # first tuple
            [
                PaymentTask(user_id=self.user_id, account_id=self.accName2Id['Amex'], amount=500),  # Amex
                PaymentTask(user_id=self.user_id, account_id=self.accName2Id['Chase'], amount=500),  # Chase
            ],
            MetaData(preferred_plan_type=PlanType.PLAN_TYPE_MIN_FEES, preferred_timeline_in_months=2.0,
                     preferred_payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY),
            [
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=['', ''],
                    timeline=2.0,
                    payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=500,
                    plan_type=PlanType.PLAN_TYPE_MIN_FEES,
                    end_date=shift_now_by_payment_frequency_multiple_times(PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                                                                           2),
                    active=True,
                    status=PaymentStatus.PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=500,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 1)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=500,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 2))
                    ])
            ]
        )

    def case_credit_score_2_month_monthly(self):
        return (  # first tuple
            [
                PaymentTask(user_id=self.user_id, account_id=self.accName2Id['Amex'], amount=500),  # Amex
                PaymentTask(user_id=self.user_id, account_id=self.accName2Id['Chase'], amount=500),  # Chase
            ],
            MetaData(preferred_plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE, preferred_timeline_in_months=2.0,
                     preferred_payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY),
            [
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=['', ''],
                    timeline=2.0,
                    payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=500,
                    plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                                                                           2),
                    active=True,
                    status=PaymentStatus.PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=500,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 1)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=500,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 2))
                    ])
            ]
        )

    def case_credit_score_no_other_meta_data(self):
        return (  # first tuple
            [
                PaymentTask(user_id=self.user_id, account_id=self.accName2Id['Amex'], amount=500),  # Amex
                PaymentTask(user_id=self.user_id, account_id=self.accName2Id['Chase'], amount=500),  # Chase
            ],
            MetaData(preferred_plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE, preferred_timeline_in_months=0.0,
                     preferred_payment_freq=None),
            [
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=['', ''],
                    timeline=3.0,
                    payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=333.34,
                    plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                                                                           3),
                    active=True,
                    status=PaymentStatus.PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=333.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 1)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=166.66,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 2)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=166.68,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 2)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=333.32,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 3))
                    ]),
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=['', ''],
                    timeline=6.0,
                    payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=166.67,
                    plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                                                                           6),
                    active=True,
                    status=PaymentStatus.PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=166.67,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 1)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=166.67,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 2)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=166.66,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 3)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=0.01,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 3)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=166.67,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 4)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=166.67,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 5)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=166.65,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 6))
                    ]),
                PaymentPlan(
                    user_id=self.user_id,
                    payment_task_id=['', ''],
                    timeline=12.0,
                    payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                    amount_per_payment=83.34,
                    plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE,
                    end_date=shift_now_by_payment_frequency_multiple_times(PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
                                                                           12),
                    active=True,
                    status=PaymentStatus.PAYMENT_STATUS_CURRENT,
                    payment_action=[
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 1)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 2)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 3)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 4)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 5)),
                        PaymentAction(
                            account_id=self.accName2Id['Amex'],
                            amount=83.30,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 6)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=0.04,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 6)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 7)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 8)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 9)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 10)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=83.34,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 11)),
                        PaymentAction(
                            account_id=self.accName2Id['Chase'],
                            amount=83.26,
                            status=PaymentActionStatus.PAYMENT_ACTION_STATUS_PENDING,
                            transaction_date=shift_now_by_payment_frequency_multiple_times(
                                PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY, 12)),
                    ])
            ]
        )


@parametrize_with_cases("paymentTasks, metaData, paymentPlans", cases=[Cases.case_min_fees_2_month_monthly,
    Cases.case_credit_score_2_month_monthly, Cases.case_credit_score_no_other_meta_data])
def test_create_payment_plan(paymentTasks, metaData, paymentPlans, gen_payment_plan_builder):
    paymentPlansCreated = gen_payment_plan_builder.createPaymentPlan(paymentTasks=paymentTasks, metaData=metaData)

    assert len(paymentPlans) == len(paymentPlansCreated), f"Tests asks for {len(paymentPlans)} but" \
                                                          f"{len(paymentPlansCreated)} were created"

    for i, (paymentPlan, paymentPlanCreated) in enumerate(zip(paymentPlans, paymentPlansCreated)):
        assert paymentPlan.user_id == paymentPlanCreated.user_id, f"Plan {i} differs"
        assert sorted(paymentPlan.payment_task_id) == sorted(paymentPlanCreated.payment_task_id), f"Plan {i} differs"
        assert paymentPlan.timeline == paymentPlanCreated.timeline, f"Plan {i} differs"
        assert paymentPlan.payment_freq == paymentPlanCreated.payment_freq, f"Plan {i} differs"
        assert paymentPlan.plan_type == paymentPlanCreated.plan_type, f"Plan {i} differs"
        assert (paymentPlan.amount_per_payment - paymentPlanCreated.amount_per_payment) < 1e-3, f"Plan {i} differs"
        assert paymentPlan.end_date == paymentPlanCreated.end_date, f"Plan {i} differs"
        assert paymentPlan.active == paymentPlanCreated.active, f"Plan {i} differs"
        assert paymentPlan.status == paymentPlanCreated.status, f"Plan {i} differs"
        assert len(paymentPlan.payment_action) == len(paymentPlanCreated.payment_action), f"Plan {i} differs"
        for ii, (paymentAction, paymentActionCreated) in enumerate(zip(paymentPlan.payment_action, paymentPlanCreated.payment_action)):
            assert paymentAction.account_id == paymentActionCreated.account_id, f"Action {ii} of plan {i} differs"
            assert (paymentAction.amount - paymentActionCreated.amount) < 1e-3, f"Action {ii} of plan {i} differs"
            assert paymentAction.transaction_date == paymentActionCreated.transaction_date, f"Action {ii} of plan {i} differs"
            assert paymentAction.status == paymentAction.status, f"Action {ii} of plan {i} differs"

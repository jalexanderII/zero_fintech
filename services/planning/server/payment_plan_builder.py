from typing import List

from bson.objectid import ObjectId
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PAYMENT_FREQUENCY_MONTHLY
from gen.Python.common.common_pb2 import PLAN_TYPE_OPTIM_CREDIT_SCORE
from gen.Python.common.payment_task_pb2 import PaymentTask
from gen.Python.planning.payment_plan_pb2 import PAYMENT_ACTION_STATUS_PENDING, PAYMENT_STATUS_CURRENT
from gen.Python.planning.payment_plan_pb2 import PaymentAction, PaymentPlan


class PaymentPlanBuilder:

    def createPaymentPlan(self, paymentTasks: List[PaymentTask]) -> List[PaymentPlan]:
        """ Creates a PaymentPlan protobuf from given list of PaymentTasks protos."""
        return self._mock_payment_plan(paymentTasks)

    @staticmethod
    def _mock_payment_plan(paymentTasks: List[PaymentTask]) -> List[PaymentPlan]:
        """A mock payment plan for end-to-end testing"""
        ids: List[str] = []
        plans: List[PaymentPlan] = []
        actions: List[PaymentAction] = []
        total = 0
        pb_timestamp = Timestamp()
        for task in paymentTasks:
            ids.append(task.payment_task_id)
            total += task.amount
            a = PaymentAction(
                account_id=task.account_id,
                amount=task.amount,
                transaction_date=pb_timestamp.GetCurrentTime(),
                status=PAYMENT_ACTION_STATUS_PENDING,
            )
            actions.append(a)
        bson_id = ObjectId()
        plan = PaymentPlan(
            payment_plan_id=str(bson_id),
            user_id=paymentTasks[0].user_id,
            payment_task_id=ids,
            timeline=12,
            payment_freq=PAYMENT_FREQUENCY_MONTHLY,
            amount_per_payment=total / 12,
            plan_type=PLAN_TYPE_OPTIM_CREDIT_SCORE,
            end_date=pb_timestamp.GetCurrentTime(),
            active=True,
            status=PAYMENT_STATUS_CURRENT,
            payment_action=actions,
        )
        plans.append(plan)
        return plans


payment_plan_builder = PaymentPlanBuilder()

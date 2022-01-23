from typing import List

import numpy as np
from bson.objectid import ObjectId
from google.protobuf.timestamp_pb2 import Timestamp

from gen.Python.common.common_pb2 import PaymentFrequency, PlanType
from gen.Python.common.payment_task_pb2 import PaymentTask
from gen.Python.core.core_pb2_grpc import  CoreStub
from gen.Python.planning.payment_plan_pb2 import PAYMENT_ACTION_STATUS_PENDING, PAYMENT_STATUS_CURRENT
from gen.Python.planning.payment_plan_pb2 import PaymentAction, PaymentPlan


class PaymentPlanBuilder:

    def __init__(self, coreClient: CoreStub):
        self.coreClient = coreClient

    def createPaymentPlan(self, paymentTasks: List[PaymentTask]) -> List[PaymentPlan]:
        """ Creates a PaymentPlan protobuf from given list of PaymentTasks protos."""

        """ Assumptions:
         - don't pay off minimum
         - we group those PaymentTasks together which shall be paid with
           the same PlanType and PaymentFrequency
        """

        # acc2amount = paymentTasks.account_to_amount
        # prefPlanType = paymentTasks.preferred_payment_type
        # prefPaymentFreq = paymentTasks.preferred_payment_frequency

        user_id = None
        paymentTaskIds = []
        accountIds = []
        amounts = []
        prefPaymentFreqs = []
        prefPlanTypes = []

        for paymentTask in paymentTasks:
            if user_id is None:
                user_id = paymentTask.user_id
            paymentTaskIds.append(paymentTask.payment_task_id)
            accountIds.append(paymentTask.account_id)
            amounts.append(paymentTask.amount)
            prefPaymentFreqs.append(paymentTask.preferred_payment_freq)
            prefPlanTypes.append(paymentTask.preferred_plan_type)
        paymentTaskIds = np.array(paymentTaskIds)
        accountIds = np.array(accountIds)
        amounts = np.array(amounts)

        payment_plans = []
        for prefPaymentFreq, prefPlanType in set(zip(prefPaymentFreqs, prefPlanTypes)):
            mask = (prefPaymentFreqs == prefPaymentFreq) and (prefPlanTypes == prefPlanType)
            # for any of them with UNKNOWN PaymentFrequency or PlanType, choose a way either adding them to existing
            # other partial matches in those tuples or simply select a few choices (at least possible in PlanType
            payment_plans.extend(self._createPaymentPlanSameTypeFreq(paymentTaskIds[mask], accountIds[mask],
                                amounts[mask], prefPaymentFreq, prefPlanType))
        return self._mock_payment_plan(paymentTasks)

    def _createPaymentPlanSameTypeFreq(self, paymentTaskIds: np.ndarray, accountIds: np.ndarray, amounts: np.ndarray,
            paymentFreq: PaymentFrequency, planType: PlanType) -> PaymentPlan:
        """Creates a PaymentPlan for given choices of PaymentFrequency and PlanType for the Accounts and Amounts. """
        pass

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
            payment_freq=PaymentFrequency.PAYMENT_FREQUENCY_MONTHLY,
            amount_per_payment=total / 12,
            plan_type=PlanType.PLAN_TYPE_OPTIM_CREDIT_SCORE,
            end_date=pb_timestamp.GetCurrentTime(),
            active=True,
            status=PAYMENT_STATUS_CURRENT,
            payment_action=actions,
        )
        plans.append(plan)
        return plans


payment_plan_builder = PaymentPlanBuilder()

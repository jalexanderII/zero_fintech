import numpy as np
import gen.planning.planning_pb2_grpc as planning_pb2_grpc
import gen.planning.payment_plan_pb2_grpc as payment_plan_pb2_grpc
# import gen.planning.payment_task_pb2_grpc as payment_task_pb2_grpc

class Planning(planning_pb2_grpc.PlanningServicer):

    def CreatePaymentPlan(self, request: planning_pb2_grpc.CreatePaymentPlanRequest, context) -> planning_pb2_grpc.CreatePaymentPlanResponse:
        user_id = None
        payment_task_ids = []
        transaction_ids = []
        account_ids = []
        pref_payment_freqs = []
        pref_plan_types = []
        pref_timelines = []
        for _payment_task in request.payment_tasks:
            if user_id is None:
                user_id = _payment_task.user_id
            payment_task_ids.append(_payment_task.payment_task_id)
            transaction_ids.append(_payment_task.transaction_id)
            account_ids.append(_payment_task.account_id)
            pref_payment_freqs.append(_payment_task.preferred_payment_freq)
            pref_plan_types.append(_payment_task.preferred_plan_type)
            pref_timelines.append(_payment_task.preferred_timeline)
        payment_task_ids = np.array(payment_task_ids)
        transaction_ids = np.array(transaction_ids)
        account_ids = np.array(account_ids)

        payment_plans = []
        for _pref_payment_freq, _pref_payment_type, _pref_timeline in set(zip(pref_payment_freqs, pref_plan_types, pref_timelines)):
            mask = pref_payment_freqs == _pref_payment_freq
            payment_plans.append(self._createPaymentPlan(payment_task_ids=payment_task_ids[mask], transaction_ids=transaction_ids[mask], user_id=user_id,
                account_ids=account_ids[mask], pref_payment_freq=_pref_payment_freq, pref_plan_type=_pref_payment_type, pref_timeline=_pref_timeline))

        return planning_pb2_grpc.CreatePaymentPlanResponse(payment_plans=payment_plans)
    
    def _createPaymentPlan(self, payment_task_ids, transaction_ids, account_ids, user_id, pref_payment_freq, pref_plan_type, pref_timeline) -> payment_plan_pb2_grpc.PaymentPlan:
        pass
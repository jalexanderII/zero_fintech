import sys, os
import grpc
from pymongo.database import Database
# make gen/Python importable by import Python.X
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir, os.pardir, os.pardir, 'gen')))
from Python.common.common_pb2 import DELETE_STATUS_SUCCESS
from Python.planning.planning_pb2_grpc import PlanningServicer as PlanningServicerPB
from Python.planning.planning_pb2 import CreatePaymentPlanRequest, CreatePaymentPlanResponse
from Python.planning.payment_plan_pb2 import GetPaymentPlanRequest, ListPaymentPlanRequest, DeletePaymentPlanRequest
from Python.planning.payment_plan_pb2 import DeletePaymentPlanResponse, ListPaymentPlanResponse, UpdatePaymentPlanRequest
from Python.planning.payment_plan_pb2 import PaymentPlan as PaymentPlanPB
# make ../database importable
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), os.pardir)))
from database.models import PaymentPlan as PaymentPlanDB
from .utils import paymentPlanDBToPB, paymentPlanPBToDB, payment_actions_pb_to_db
from payment_plan_builder import PaymentPlanBuilder
from dotenv import load_dotenv


class PlanningServicer(PlanningServicerPB):

    def __init__(self, planningCollection: Database) -> None:
        self.planningCollection = planningCollection
        self.ppb = PaymentPlanBuilder()

    def CreatePaymentPlan(self, request: CreatePaymentPlanRequest, context: grpc.ServicerContext) -> CreatePaymentPlanResponse:
        """ Creates PaymentPlan(s) for given request containing a list of PaymentTasks"""
        paymentPlanListPB = self.ppb.createPaymentPlan(request.payment_tasks)
        # for paymentPlanPB in paymentPlanListPB:
        #      self._createPaymentPlan(paymentPlanPB)
        return CreatePaymentPlanResponse(payment_plans=paymentPlanListPB)

    # TODO(JB): this needs fixing returning error on save
    def _createPaymentPlan(self, paymentPlanPB: PaymentPlanPB) -> None:
        """ Saves a PaymentPlan into the database without creating it."""
        paymentPlanDB = paymentPlanPBToDB(paymentPlanPB)
        paymentPlanDB.save()
    
    def GetPaymentPlan(self, request: GetPaymentPlanRequest, context) -> PaymentPlanPB:
        paymentPlanDB = PaymentPlanDB.objects.get(id=request.payment_plan_id)
        return paymentPlanDBToPB(paymentPlanDB)

    def ListPaymentPlans(self, request: ListPaymentPlanRequest, context) -> ListPaymentPlanResponse:
        paymentPlansDB = PaymentPlanDB.objects
        paymentPlansPB = []
        for pp in paymentPlansDB:
            paymentPlansPB.append(paymentPlanDBToPB(pp))
        return ListPaymentPlanResponse(payment_plans=paymentPlansPB)

    def UpdatePaymentPlan(self, request: UpdatePaymentPlanRequest, context) -> PaymentPlanPB:
        paymentPlanPB = request.payment_plan
        PaymentPlanDB.objects.get(id=paymentPlanPB.payment_plan_id).update(
            UserID=paymentPlanPB.user_id,
            PaymentTaskID=list(paymentPlanPB.payment_task_id), 
            Timeline=paymentPlanPB.timeline,
            PaymentFrequency=paymentPlanPB.payment_freq,
            AmountPerPayment=paymentPlanPB.amount_per_payment,
            PlanType=paymentPlanPB.plan_type,
            EndDate=paymentPlanPB.end_date.ToDatetime(),
            Active=paymentPlanPB.active,
            Status=paymentPlanPB.status,
            PaymentAction=payment_actions_pb_to_db(paymentPlanPB.payment_action)
        )
        return paymentPlanPB

    def DeletePaymentPlan(self, request: DeletePaymentPlanRequest, context) -> DeletePaymentPlanResponse:
        paymentPlanDB = PaymentPlanDB.objects.get(id=request.payment_plan_id)
        paymentPlanPB = paymentPlanDBToPB(paymentPlanDB)
        paymentPlanDB.delete()
        return DeletePaymentPlanResponse(status=DELETE_STATUS_SUCCESS, payment_plan=paymentPlanPB)
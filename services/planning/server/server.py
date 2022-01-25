import logging
import sys

from bson.objectid import ObjectId

from database.models.common import PaymentPlan as PaymentPlanDB
from pymongo.collection import Collection

from gen.Python.common.common_pb2 import DELETE_STATUS_SUCCESS, DELETE_STATUS_FAILED
from gen.Python.planning.payment_plan_pb2 import DeletePaymentPlanRequest
from gen.Python.planning.payment_plan_pb2 import DeletePaymentPlanResponse, ListPaymentPlanResponse
from gen.Python.planning.payment_plan_pb2 import GetPaymentPlanRequest, ListPaymentPlanRequest
from gen.Python.common.payment_plan_pb2 import PaymentPlan as PaymentPlanPB
from gen.Python.common.payment_task_pb2 import PaymentPlanResponse
from gen.Python.planning.payment_plan_pb2 import UpdatePaymentPlanRequest
from gen.Python.planning.planning_pb2 import CreatePaymentPlanRequest
from gen.Python.planning.planning_pb2_grpc import PlanningServicer as PlanningServicerPB
from .payment_plan_builder import PaymentPlanBuilder, payment_plan_builder
from .utils import paymentPlanDBToPB, paymentPlanPBToDB

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)]
)
logger = logging.getLogger('PlanningServicer')


class PlanningServicer(PlanningServicerPB):

    def __init__(self, planning_collection: Collection) -> None:
        self.planning_collection = planning_collection
        self.payment_plan_builder: PaymentPlanBuilder = payment_plan_builder

    def CreatePaymentPlan(self, request: CreatePaymentPlanRequest, context) -> PaymentPlanResponse:
        """ Creates PaymentPlan(s) for given request containing a list of PaymentTasks"""
        logger.info('CreatePaymentPlan called')
        paymentPlanListPB = self.payment_plan_builder.createPaymentPlan(request.payment_tasks)
        for paymentPlanPB in paymentPlanListPB:
            self._createPaymentPlan(paymentPlanPB)
        return PaymentPlanResponse(payment_plans=paymentPlanListPB)

    def _createPaymentPlan(self, payment_plan_PB: PaymentPlanPB) -> str:
        """ Saves a PaymentPlan into the database without creating it."""
        paymentPlanDB = paymentPlanPBToDB(payment_plan_PB).to_dict()
        new_payment_plan = self.planning_collection.insert_one(paymentPlanDB)
        return new_payment_plan.inserted_id

    def GetPaymentPlan(self, request: GetPaymentPlanRequest, context) -> PaymentPlanPB:
        logger.info('GetPaymentPlan called')
        paymentPlanDB = self.planning_collection.find_one({"_id": ObjectId(request.payment_plan_id)})
        pp_id = paymentPlanDB["_id"]
        paymentPlanDB = PaymentPlanDB().from_dict(paymentPlanDB)
        paymentPlanDB.payment_plan_id = str(pp_id)
        return paymentPlanDBToPB(paymentPlanDB)

    def ListPaymentPlans(self, request: ListPaymentPlanRequest, context) -> ListPaymentPlanResponse:
        logger.info('ListPaymentPlans called')
        payment_plans = self.planning_collection.find()
        paymentPlansPB = []
        for payment_plan in payment_plans:
            pp_id = payment_plan["_id"]
            paymentPlanDB = PaymentPlanDB().from_dict(payment_plan)
            paymentPlanDB.payment_plan_id = str(pp_id)
            paymentPlansPB.append(paymentPlanDBToPB(paymentPlanDB))
        return ListPaymentPlanResponse(payment_plans=paymentPlansPB)

    def UpdatePaymentPlan(self, request: UpdatePaymentPlanRequest, context) -> PaymentPlanPB:
        logger.info('UpdatePaymentPlan called')
        paymentPlanDB = paymentPlanPBToDB(request.payment_plan)
        payment_plan = {k: v for k, v in paymentPlanDB.to_dict().items() if v is not None}
        _ = self.planning_collection.update_one({"_id": ObjectId(request.payment_plan_id)}, {"$set": payment_plan})
        updated_payment_plan = self.planning_collection.find_one({"_id": ObjectId(request.payment_plan_id)})
        paymentPlanDB = PaymentPlanDB().from_dict(updated_payment_plan)
        paymentPlanDB.payment_plan_id = request.payment_plan_id
        return paymentPlanDBToPB(paymentPlanDB)

    def DeletePaymentPlan(self, request: DeletePaymentPlanRequest, context) -> DeletePaymentPlanResponse:
        logger.info('DeletePaymentPlan called')
        # fetch
        paymentPlanDB = self.planning_collection.find_one({"_id": ObjectId(request.payment_plan_id)})
        paymentPlanDB = PaymentPlanDB().from_dict(paymentPlanDB)
        paymentPlanDB.payment_plan_id = request.payment_plan_id
        paymentPlanPB = paymentPlanDBToPB(paymentPlanDB)
        # delete
        delete_result = self.planning_collection.delete_one({"_id": ObjectId(request.payment_plan_id)})
        if delete_result.deleted_count == 1:
            return DeletePaymentPlanResponse(status=DELETE_STATUS_SUCCESS, payment_plan=paymentPlanPB)
        return DeletePaymentPlanResponse(status=DELETE_STATUS_FAILED, payment_plan=paymentPlanPB)

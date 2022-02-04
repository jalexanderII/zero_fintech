import logging
import sys
from typing import List

from attr import define, field
from bson.objectid import ObjectId
from database.models.common import PaymentPlan as PaymentPlanDB
from pymongo.collection import Collection

from gen.Python.common.common_pb2 import DELETE_STATUS_SUCCESS, DELETE_STATUS_FAILED
from gen.Python.common.payment_plan_pb2 import DeletePaymentPlanRequest
from gen.Python.common.payment_plan_pb2 import DeletePaymentPlanResponse
from gen.Python.common.payment_plan_pb2 import GetPaymentPlanRequest
from gen.Python.common.payment_plan_pb2 import ListPaymentPlanRequest
from gen.Python.common.payment_plan_pb2 import ListPaymentPlanResponse
from gen.Python.common.payment_plan_pb2 import PaymentPlan as PaymentPlanPB
from gen.Python.common.payment_plan_pb2 import UpdatePaymentPlanRequest
from gen.Python.common.payment_plan_pb2 import PaymentPlanResponse
from gen.Python.planning.planning_pb2 import CreatePaymentPlanRequest
from gen.Python.planning.planning_pb2_grpc import PlanningServicer
from services.planning.server.payment_plan_builder import PaymentPlanBuilder
from services.planning.server.payment_plan_builder import payment_plan_builder
from services.planning.server.utils import payment_plan_PB_to_DB, payment_plan_DB_to_PB

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[logging.StreamHandler(sys.stdout)],
)
logger = logging.getLogger("PlanningServicer")


@define
class PlanningService(PlanningServicer):
    planning_collection: Collection
    _payment_plan_builder: PaymentPlanBuilder = field(
        init=False, default=payment_plan_builder
    )

    def CreatePaymentPlan(
        self, request: CreatePaymentPlanRequest, ctx=None
    ) -> PaymentPlanResponse:
        """Calls PaymentPlanBuilder to generate a list of Payments plans given a list of PaymentTasks"""
        logger.info("CreatePaymentPlan called")
        payment_tasks, meta_data = request.payment_tasks, request.meta_data

        payment_plans_pb: List[PaymentPlanPB] = self._payment_plan_builder.create(
            payment_tasks, meta_data
        )
        resp = PaymentPlanResponse(payment_plans=payment_plans_pb)

        map(lambda plan: self.SavePaymentPlan(plan), payment_plans_pb)
        return resp

    def SavePaymentPlan(self, payment_plan_pb: PaymentPlanPB) -> str:
        """Adds a given PaymentPlan to the database"""
        payment_plan_db = payment_plan_PB_to_DB(payment_plan_pb).to_dict()
        return self.planning_collection.insert_one(payment_plan_db).inserted_id

    def GetPaymentPlan(self, request: GetPaymentPlanRequest, ctx=None) -> PaymentPlanPB:
        logger.info("GetPaymentPlan called")
        payment_plan_response = self.planning_collection.find_one(
            {"_id": ObjectId(request.payment_plan_id)}
        )
        _id = payment_plan_response["_id"]
        payment_plan_db = PaymentPlanDB().from_dict(payment_plan_response)
        payment_plan_db.payment_plan_id = str(_id)
        return payment_plan_DB_to_PB(payment_plan_db)

    def ListPaymentPlans(
        self, request: ListPaymentPlanRequest, ctx=None
    ) -> ListPaymentPlanResponse:
        logger.info("ListPaymentPlans called")
        payment_plans_pb: List[PaymentPlanPB] = []
        payment_plans = self.planning_collection.find()
        for payment_plan in payment_plans:
            payment_plan_db = PaymentPlanDB().from_dict(payment_plan)
            payment_plan_db.payment_plan_id = str(payment_plan["_id"])
            payment_plans_pb.append(payment_plan_DB_to_PB(payment_plan_db))

        return ListPaymentPlanResponse(payment_plans=payment_plans_pb)

    def UpdatePaymentPlan(
        self, request: UpdatePaymentPlanRequest, ctx=None
    ) -> PaymentPlanPB:
        logger.info("UpdatePaymentPlan called")
        payment_plan = {
            k: v
            for k, v in payment_plan_PB_to_DB(request.payment_plan).to_dict().items()
            if v is not None
        }

        _ = self.planning_collection.update_one(
            {"_id": ObjectId(request.payment_plan_id)}, {"$set": payment_plan}
        )
        updated_payment_plan = self.planning_collection.find_one(
            {"_id": ObjectId(request.payment_plan_id)}
        )
        payment_plan_db = PaymentPlanDB().from_dict(updated_payment_plan)
        payment_plan_db.payment_plan_id = request.payment_plan_id

        return payment_plan_DB_to_PB(payment_plan_db)

    def DeletePaymentPlan(
        self, request: DeletePaymentPlanRequest, ctx=None
    ) -> DeletePaymentPlanResponse:
        logger.info("DeletePaymentPlan called")

        payment_plan_db = self.planning_collection.find_one(
            {"_id": ObjectId(request.payment_plan_id)}
        )
        payment_plan_db = PaymentPlanDB().from_dict(payment_plan_db)
        payment_plan_db.payment_plan_id = request.payment_plan_id

        resp = self.planning_collection.delete_one(
            {"_id": ObjectId(request.payment_plan_id)}
        )

        return DeletePaymentPlanResponse(
            status=DELETE_STATUS_SUCCESS
            if resp.deleted_count == 1
            else DELETE_STATUS_FAILED,
            payment_plan=payment_plan_DB_to_PB(payment_plan_db),
        )

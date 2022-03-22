import datetime
import logging
import os
import sys
from collections import defaultdict
from typing import List

import grpc
from attr import define, field
from bson.objectid import ObjectId
from pymongo.collection import Collection
from pymongo.results import InsertOneResult

from gen.Python.common.common_pb2 import PaymentActionStatus, PAYMENT_ACTION_STATUS_PENDING
from gen.Python.common.common_pb2 import DELETE_STATUS_SUCCESS, DELETE_STATUS_FAILED
from gen.Python.common.payment_plan_pb2 import DeletePaymentPlanRequest, ListUserPaymentPlansRequest
from gen.Python.common.payment_plan_pb2 import DeletePaymentPlanResponse
from gen.Python.common.payment_plan_pb2 import GetPaymentPlanRequest
from gen.Python.common.payment_plan_pb2 import ListPaymentPlanRequest
from gen.Python.common.payment_plan_pb2 import ListPaymentPlanResponse
from gen.Python.common.payment_plan_pb2 import PaymentPlan as PaymentPlanPB
from gen.Python.common.payment_plan_pb2 import PaymentPlanResponse
from gen.Python.common.payment_plan_pb2 import UpdatePaymentPlanRequest
from gen.Python.core.accounts_pb2 import Account, ListUserAccountsRequest, ListAccountResponse
from gen.Python.core.core_pb2_grpc import CoreStub
from gen.Python.planning.planning_pb2 import CreatePaymentPlanRequest, GetUserOverviewRequest
from gen.Python.planning.planning_pb2 import GetAmountPaidPercentageResponse, GetPercentageCoveredByPlansResponse
from gen.Python.planning.planning_pb2 import WaterfallOverviewResponse, WaterfallMonth
from gen.Python.planning.planning_pb2_grpc import PlanningServicer
from services.planning.database.models.common import PaymentPlan as PaymentPlanDB
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
    core_client: CoreStub = CoreStub(grpc.insecure_channel(f'localhost:{os.getenv("CORE_SERVER_PORT")}'))

    def CreatePaymentPlan(
            self, request: CreatePaymentPlanRequest, ctx=None
    ) -> PaymentPlanResponse:
        """Calls PaymentPlanBuilder to generate a list of Payments plans given a list of PaymentTasks"""
        logger.info("CreatePaymentPlan called")
        payment_tasks, meta_data = request.payment_tasks, request.meta_data

        payment_plans_pb: List[PaymentPlanPB] = self._payment_plan_builder.create(
            payment_tasks, meta_data
        )
        if request.save_plan:
            for payment_plan in payment_plans_pb:
                new_id = self.SavePaymentPlan(payment_plan)
                logger.info(f"New plan created with id {new_id}")

        return PaymentPlanResponse(payment_plans=payment_plans_pb)

    def SavePaymentPlan(self, payment_plan_pb: PaymentPlanPB) -> str:
        """Adds a given PaymentPlan to the database"""
        logger.info("SavePaymentPlan called")
        payment_plan_db = payment_plan_PB_to_DB(payment_plan_pb).to_dict()
        resp: InsertOneResult = self.planning_collection.insert_one(payment_plan_db)
        return str(resp.inserted_id)

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

    def ListUserPaymentPlans(
            self, request: ListUserPaymentPlansRequest, ctx=None
    ) -> ListPaymentPlanResponse:
        logger.info("ListUserPaymentPlans called")
        payment_plans_pb: List[PaymentPlanPB] = []
        payment_plans = self.planning_collection.find({"userId": request.user_id})
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

    def GetWaterfallOverview(self, request: GetUserOverviewRequest, ctx=None) -> WaterfallOverviewResponse:
        logger.info("GetWaterfallOverview called")

        payment_plans_cursor = self.planning_collection.find({'userId': request.user_id, 'active': True})
        # for this month and 11 month afterwards have a dictionary mapping account_id to amount to be paid
        waterfall = [defaultdict(float) for _ in range(12)]
        now = datetime.datetime.now()
        for _pp in payment_plans_cursor:
            pp = PaymentPlanDB().from_dict(_pp)
            for pa in pp.payment_action:
                _waterfall_month = pa.transaction_date.month - now.month  # if in same month, before or after
                if 0 <= _waterfall_month <= 11:
                    waterfall[_waterfall_month][pa.account_id] += pa.amount  # update account with amount
        monthly_waterfall = [WaterfallMonth(account_to_amounts=_waterfall_month) for _waterfall_month in waterfall]
        return WaterfallOverviewResponse(monthly_waterfall=monthly_waterfall)

    def GetAmountPaidPercentage(self, request: GetUserOverviewRequest, ctx=None) -> GetAmountPaidPercentageResponse:
        logger.info("GetAmountPaidPercentage called")

        # retrieve all active PaymentPlans for user
        # iterate over all PaymentActions and depending on its status add it to paid amount or only total_amount
        payment_plans_cursor = self.planning_collection.find({'userId': request.user_id, 'active': True})
        amount_paid, total_amount = 0, 0
        for _pp in payment_plans_cursor:
            pp = PaymentPlanDB().from_dict(_pp)
            for pa in pp.payment_action:
                _amount = pa.amount
                if pa.status == PaymentActionStatus.PAYMENT_ACTION_STATUS_COMPLETED:
                    amount_paid += _amount
                total_amount += _amount
        prcnt_paid = amount_paid / total_amount if total_amount > 0 else 1
        return GetAmountPaidPercentageResponse(percentage_paid=prcnt_paid)

    def GetPercentageCoveredByPlans(
            self, request: GetUserOverviewRequest, ctx=None
    ) -> GetPercentageCoveredByPlansResponse:
        logger.info("GetPercentageCoveredByPlans")

        user_id = request.user_id
        # retrieve all accounts
        accounts = self._fetch_accounts(user_id=user_id)
        # get total balance
        acc2balance = {}
        for acc in accounts:
            balance = acc.current_balance
            if balance > 0:
                acc2balance[acc.account_id] = balance
        # retrieve all active PaymentPlans for user and all accounts
        # and see how much every account is covered
        payment_plans_cursor = self.planning_collection.find({'userId': user_id, 'active': True})
        acc2coverage = defaultdict(float)
        for _pp in payment_plans_cursor:
            pp = PaymentPlanDB().from_dict(_pp)
            for pa in pp.payment_action:
                if pa.status == PAYMENT_ACTION_STATUS_PENDING:
                    acc2coverage[pa.account_id] += pa.amount
        # see coverage in percentage
        total_balance = sum(acc2balance.values())
        total_coverage_prcnt = sum(acc2coverage.values()) / total_balance if total_balance > 0 else 1
        acc2coverage_prcnt = {}
        for acc_id in acc2balance.keys():
            balance = acc2balance[acc_id]
            acc2coverage_prcnt[acc_id] = acc2coverage[acc_id] / balance if balance > 0 else 1

        return GetPercentageCoveredByPlansResponse(
            overall_covered=total_coverage_prcnt, account_to_percent_covered=acc2coverage_prcnt)

    def _fetch_accounts(self, user_id) -> List[Account]:
        account_resp: ListAccountResponse = self.core_client.ListUserAccounts(ListUserAccountsRequest(user_id=user_id))
        return account_resp.accounts

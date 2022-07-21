import os
from collections import defaultdict
from typing import List

import grpc
import numpy as np
from fastapi import APIRouter, Depends
from motor.motor_asyncio import AsyncIOMotorCollection

from gen.Python.common.common_pb2 import (
    PAYMENT_ACTION_STATUS_PENDING,
    PAYMENT_STATUS_CURRENT,
    PAYMENT_STATUS_IN_DEFAULT,
    PAYMENT_STATUS_COMPLETED,
    PAYMENT_ACTION_STATUS_COMPLETED,
    PAYMENT_ACTION_STATUS_IN_DEFAULT,
    PAYMENT_STATUS_CANCELLED,
    PaymentStatus,
    PlanType,
    PaymentFrequency,
)
from gen.Python.core.accounts_pb2 import Account, GetAccountRequest
from gen.Python.core.core_pb2_grpc import CoreStub
from gen.Python.core.users_pb2 import GetUserRequest
from services.planning.app.depends import get_async_mongo_collection, get_core_client
from services.planning.app.models.dashboard_models import (
    PercentCoveredByPlansResponseSchema,
    StatisticalOverviewCategoricalDataResponseSchema,
    StatisticalOverviewContinuousDataResponseSchema,
    OverviewAllPreferencesResponseSchema,
)
from services.planning.database.models.common import PaymentPlan as PaymentPlanDB

router = APIRouter()


@router.get("/overview_amount_per_payment")
async def get_overview_amount_per_payment(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> StatisticalOverviewContinuousDataResponseSchema:
    amount_per_payment = []
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        amount_per_payment.append(pp.amount_per_payment)
    amount_per_payment = np.array(amount_per_payment)
    return StatisticalOverviewContinuousDataResponseSchema(
        mean=np.mean(amount_per_payment),
        std=np.std(amount_per_payment),
        median=np.median(amount_per_payment),
        p25=np.percentile(amount_per_payment, q=25),
        p75=np.percentile(amount_per_payment, q=75),
    )


@router.get("/overview_payment_plan_amounts")
async def get_overview_payment_plan_amounts(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> StatisticalOverviewContinuousDataResponseSchema:
    amounts = []
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        amounts.append(pp.amount)
    amounts = np.array(amounts)
    return StatisticalOverviewContinuousDataResponseSchema(
        mean=np.mean(amounts),
        std=np.std(amounts),
        median=np.median(amounts),
        p25=np.percentile(amounts, q=25),
        p75=np.percentile(amounts, q=75),
    )


@router.get("/overview_percent_covered_by_plans")
async def get_overview_percent_covered_by_plans(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
    core_client: CoreStub = Depends(get_core_client),
) -> PercentCoveredByPlansResponseSchema:
    # calculate amount covered by plans in total and by user
    total_plan_amounts, user_plan_amounts = 0, defaultdict(float)
    async for _pp in planning_collection.find({"active": True}):
        pp = PaymentPlanDB().from_dict(_pp)
        amount = 0
        for pa in pp.payment_action:
            if pa.status == PAYMENT_ACTION_STATUS_PENDING:
                amount += pa.amount
        total_plan_amounts += amount
        user_plan_amounts[pp.user_id] += amount
    # calculate total amount of all users and the percentage of coverage for every user
    total_amount, percent_covered_users = 0, {}
    for user_id in user_plan_amounts.keys():
        accounts = _fetch_accounts(core_client, user_id)
        amount = 0
        for account in accounts:
            amount += account.current_balance
        total_amount += amount
        percent_covered_users[user_id] = (
            user_plan_amounts[user_id] / amount if amount > 0 else 1
        )
    percent_covered_users = np.array(list(percent_covered_users.values()))
    return PercentCoveredByPlansResponseSchema(
        percent_amount_by_plans_vs_overall=total_plan_amounts / total_amount,
        mean=np.mean(percent_covered_users),
        std=np.std(percent_covered_users),
        median=np.median(percent_covered_users),
        p25=np.percentile(percent_covered_users, q=25),
        p75=np.percentile(percent_covered_users, q=75),
    )


@router.get("/overview_current_vs_defaulted_plans")
async def get_overview_current_defaulted_plans(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> StatisticalOverviewCategoricalDataResponseSchema:
    category_to_count = defaultdict(int)
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        if pp.status in [PAYMENT_STATUS_CURRENT, PAYMENT_STATUS_IN_DEFAULT]:
            category_to_count[PaymentStatus.Name(pp.status)] += 1
    return StatisticalOverviewCategoricalDataResponseSchema(
        category_to_count=category_to_count
    )


@router.get("/overview_amounts_completed_pending_defaulted")
async def get_overview_amounts_completed_pending_in_default(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> StatisticalOverviewCategoricalDataResponseSchema:
    category_to_count = defaultdict(float)
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        if pp.status == PAYMENT_STATUS_COMPLETED:
            category_to_count["completed"] += pp.amount
        elif pp.status in [PAYMENT_STATUS_CURRENT, PAYMENT_STATUS_CANCELLED]:
            for pa in pp.payment_action:
                if pa.status == PAYMENT_ACTION_STATUS_PENDING:
                    category_to_count["pending"] += pa.amount
                elif pa.status == PAYMENT_ACTION_STATUS_COMPLETED:
                    category_to_count["completed"] += pa.amount
        elif pp.status == PAYMENT_STATUS_IN_DEFAULT:
            # TODO: validate if below assumption is correct:
            # if one PaymentAction defaulted does that mean that all following ones, should be also considered
            # defaulted?
            for pa in pp.payment_action:
                if pa.status in [
                    PAYMENT_ACTION_STATUS_PENDING,
                    PAYMENT_ACTION_STATUS_IN_DEFAULT,
                ]:
                    category_to_count["in_default"] += pa.amount
                elif pa.status == PAYMENT_ACTION_STATUS_COMPLETED:
                    category_to_count["completed"] += pa.amount
    return StatisticalOverviewCategoricalDataResponseSchema(
        category_to_count=category_to_count
    )


@router.get("/overview_timeline_preferences")
async def get_overview_timeline_preferences(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> StatisticalOverviewContinuousDataResponseSchema:
    timeline = []
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        timeline.append(pp.timeline)
    timeline = np.array(timeline)
    return StatisticalOverviewContinuousDataResponseSchema(
        mean=np.mean(timeline),
        std=np.std(timeline),
        median=np.median(timeline),
        p25=np.percentile(timeline, q=25),
        p75=np.percentile(timeline, q=75),
    )


@router.get("/overview_plan_type_preferences")
async def get_overview_plan_type_preferences(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> StatisticalOverviewCategoricalDataResponseSchema:
    plan_type_counts = defaultdict(int)
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        plan_type_counts[PlanType.Name(pp.plan_type)] += 1
    return StatisticalOverviewCategoricalDataResponseSchema(
        category_to_count=plan_type_counts
    )


@router.get("/overview_payment_frequency_preferences")
async def get_overview_payment_frequency_preferences(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> StatisticalOverviewCategoricalDataResponseSchema:
    payment_freq_counts = defaultdict(int)
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        payment_freq_counts[PaymentFrequency.Name(pp.payment_freq)] += 1
    return StatisticalOverviewCategoricalDataResponseSchema(
        category_to_count=payment_freq_counts
    )


@router.get("/overview_all_preferences")
async def get_overview_all_preferences(
    planning_collection: AsyncIOMotorCollection = Depends(get_async_mongo_collection),
) -> OverviewAllPreferencesResponseSchema:
    timeline, plan_type_counts, payment_freq_counts = (
        [],
        defaultdict(int),
        defaultdict(int),
    )
    async for _pp in planning_collection.find({}):
        pp = PaymentPlanDB().from_dict(_pp)
        timeline.append(pp.timeline)
        plan_type_counts[PlanType.Name(pp.plan_type)] += 1
        payment_freq_counts[PaymentFrequency.Name(pp.payment_freq)] += 1
    timeline = np.array(timeline)
    overview_timeline_preferences = StatisticalOverviewContinuousDataResponseSchema(
        mean=np.mean(timeline),
        std=np.std(timeline),
        median=np.median(timeline),
        p25=np.percentile(timeline, q=25),
        p75=np.percentile(timeline, q=75),
    )
    overview_plan_type_option_preferences = (
        StatisticalOverviewCategoricalDataResponseSchema(
            category_to_count=plan_type_counts
        )
    )
    overview_payment_frequency_preferences = (
        StatisticalOverviewCategoricalDataResponseSchema(
            category_to_count=payment_freq_counts
        )
    )
    return OverviewAllPreferencesResponseSchema(
        overview_timeline_preferences=overview_timeline_preferences,
        overview_plan_type_option_preferences=overview_plan_type_option_preferences,
        overview_payment_frequency_preferences=overview_payment_frequency_preferences,
    )


def _fetch_accounts(core_client: CoreStub, user_id: str) -> List[Account]:
    user = core_client.GetUser(GetUserRequest(id=user_id))
    account_ids: List[str] = list(user.account_id_to_token.keys())

    return [
        core_client.GetAccount(GetAccountRequest(id=acc_id)) for acc_id in account_ids
    ]

from typing import Any, Dict, Union

from pydantic import BaseModel


class StatisticalOverviewContinuousDataResponseSchema(BaseModel):
    mean: float
    std: float
    median: float
    p25: float
    p75: float


class StatisticalOverviewCategoricalDataResponseSchema(BaseModel):
    category_to_count: Dict[Any, Union[float, int]]


class PercentCoveredByPlansResponseSchema(
    StatisticalOverviewContinuousDataResponseSchema
):
    percent_amount_by_plans_vs_overall: float


class OverviewAllPreferencesResponseSchema(BaseModel):
    overview_timeline_preferences: StatisticalOverviewContinuousDataResponseSchema
    overview_plan_type_option_preferences: StatisticalOverviewCategoricalDataResponseSchema
    overview_payment_frequency_preferences: StatisticalOverviewCategoricalDataResponseSchema

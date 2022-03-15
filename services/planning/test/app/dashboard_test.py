import os
from typing import Generator, Any
from unittest.mock import MagicMock

import pytest
from dotenv import load_dotenv
from motor.motor_asyncio import AsyncIOMotorCollection, AsyncIOMotorClient
from pytest_mock import MockerFixture
from starlette.testclient import TestClient

from services.planning.app.api import dashboard
from services.planning.app.depends import get_async_mongo_collection
from services.planning.app.main import create_application
from services.planning.test.helpers.shared_objects import MOCK_CHASE_ACC, MOCK_AMEX_ACC


async def get_async_mongo_test_collection() -> Generator[AsyncIOMotorCollection, Any, None]:
    try:
        load_dotenv()
        client = AsyncIOMotorClient(os.getenv("MONGOURI"))
        db = client[os.getenv("PLANNING_DB_NAME")]
        yield db[f'{os.getenv("PLANNING_COLLECTION")}_TEST']
    finally:
        client.close()


@pytest.fixture
def test_client() -> Generator[TestClient, Any, None]:
    application = create_application()
    application.dependency_overrides[get_async_mongo_collection] = get_async_mongo_test_collection
    with TestClient(application) as test_client:
        yield test_client


@pytest.fixture
def patch__fetch_accounts(mocker: MockerFixture) -> MagicMock:
    return mocker.patch.object(dashboard, '_fetch_accounts')


def test_overview_amount_per_payment(test_client: TestClient):
    response = test_client.get('/dashboard/overview_amount_per_payment')

    response_body = response.json()
    assert response_body['mean'] > 0
    assert response_body['std'] > 0
    assert response_body['median'] > 0
    assert response_body['p25'] > 0
    assert response_body['p75'] > 0


def test_overview_payment_plan_amounts(test_client: TestClient):
    response = test_client.get('/dashboard/overview_payment_plan_amounts')

    response_body = response.json()
    assert response_body['mean'] > 0
    assert response_body['std'] > 0
    assert response_body['median'] > 0
    assert response_body['p25'] > 0
    assert response_body['p75'] > 0


def test_overview_percent_covered_by_plans(test_client: TestClient, patch__fetch_accounts: MagicMock):
    patch__fetch_accounts.return_value = [MOCK_CHASE_ACC, MOCK_AMEX_ACC]

    response = test_client.get('/dashboard/overview_percent_covered_by_plans')

    response_body = response.json()
    assert response_body['mean'] > 0
    assert response_body['median'] > 0
    assert response_body['p25'] > 0
    assert response_body['p75'] > 0


def test_overview_current_vs_defaulted_plans(test_client: TestClient):
    response = test_client.get('/dashboard/overview_current_vs_defaulted_plans')

    category_to_count = response.json()['category_to_count']
    assert len(category_to_count.keys()) > 0
    assert len(list(filter(lambda x: x > 0, category_to_count.values()))) == len(category_to_count.values())


def test_overview_amounts_completed_pending_defaulted(test_client: TestClient):
    response = test_client.get('/dashboard/overview_amounts_completed_pending_defaulted')

    category_to_count = response.json()['category_to_count']
    assert len(category_to_count.keys()) > 0
    assert len(list(filter(lambda x: x > 0, category_to_count.values()))) == len(category_to_count.values())


def test_overview_timeline_preferences(test_client: TestClient):
    response = test_client.get('/dashboard/overview_timeline_preferences')

    response_body = response.json()
    assert response_body['mean'] > 0
    assert response_body['std'] > 0
    assert response_body['median'] > 0
    assert response_body['p25'] > 0
    assert response_body['p75'] > 0


def test_overview_plan_type_preferences(test_client: TestClient):
    response = test_client.get('/dashboard/overview_plan_type_preferences')

    category_to_count = response.json()['category_to_count']
    assert len(category_to_count.keys()) > 0
    assert len(list(filter(lambda x: x > 0, category_to_count.values()))) == len(category_to_count.values())


def test_overview_payment_frequency_preferences(test_client: TestClient):
    response = test_client.get('/dashboard/overview_payment_frequency_preferences')

    category_to_count = response.json()['category_to_count']
    assert len(category_to_count.keys()) > 0
    assert len(list(filter(lambda x: x > 0, category_to_count.values()))) == len(category_to_count.values())


def test_overview_all_preferences(test_client: TestClient):
    response = test_client.get('/dashboard/overview_all_preferences')
    response_body = response.json()

    overview_timeline_preferences = response_body['overview_timeline_preferences']
    assert overview_timeline_preferences['mean'] > 0
    assert overview_timeline_preferences['std'] > 0
    assert overview_timeline_preferences['median'] > 0
    assert overview_timeline_preferences['p25'] > 0
    assert overview_timeline_preferences['p75'] > 0

    for categorical_var in ['overview_plan_type_option_preferences', 'overview_payment_frequency_preferences']:
        category_to_count = response_body[categorical_var]['category_to_count']
        assert len(category_to_count.keys()) > 0
        assert len(list(filter(lambda x: x > 0, category_to_count.values()))) == len(category_to_count.values())





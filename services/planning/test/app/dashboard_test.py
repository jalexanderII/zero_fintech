import os
from datetime import datetime, timedelta
from typing import Generator, Any
from unittest.mock import MagicMock

import pytest
from bson import ObjectId
from dotenv import load_dotenv
from motor.motor_asyncio import AsyncIOMotorCollection, AsyncIOMotorClient
from pytest_mock import MockerFixture
from starlette.testclient import TestClient

from services.planning.app.api import dashboard
from services.planning.app.depends import get_async_mongo_collection
from services.planning.app.main import create_application
from services.planning.test.helpers.shared_objects import MOCK_CHASE_ACC, MOCK_AMEX_ACC


# async def get_async_mongo_test_collection() -> Generator[AsyncIOMotorCollection, Any, None]:
#     try:
#         load_dotenv()
#         client = AsyncIOMotorClient(os.getenv("MONGOURI"))
#         db = client[os.getenv("PLANNING_DB_NAME")]
#         yield db[f'{os.getenv("PLANNING_COLLECTION")}_TEST']
#     finally:
#         client.close()


@pytest.fixture
def test_client() -> Generator[TestClient, Any, None]:
    application = create_application()
    # application.dependency_overrides[get_async_mongo_collection] = get_async_mongo_test_collection
    with TestClient(application) as test_client:
        yield test_client


@pytest.fixture
def patch__fetch_accounts(mocker: MockerFixture) -> None:
    mocker.patch.object(dashboard, '_fetch_accounts').return_value = [MOCK_CHASE_ACC, MOCK_AMEX_ACC]


@pytest.fixture
def patch_planning_collection_find(mocker: MockerFixture) -> None:
    async def async_iterator(l):
        for item in l:
            yield item
    mocker.patch.object(AsyncIOMotorCollection, 'find').return_value = async_iterator([
        {'_id': ObjectId('62360df5526e120f5cb047fe'), 'userId': 'test',
         'paymentTaskId': ['62360df57e70ed11a0b2109b'], 'amount': 31.61, 'timeline': 4.0,
         'paymentFreq': 'PAYMENT_FREQUENCY_MONTHLY', 'amountPerPayment': 7.91, 'planType': 'PLAN_TYPE_MIN_FEES',
         'endDate': '2022-07-19T00:00:00Z', 'active': True, 'status': 'PAYMENT_STATUS_CURRENT', 'paymentAction': [
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.91, 'transactionDate': '2022-04-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.91, 'transactionDate': '2022-05-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.91, 'transactionDate': '2022-06-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.879999999999999,
             'transactionDate': '2022-07-19T00:00:00Z', 'status': 'PAYMENT_ACTION_STATUS_PENDING'}]},
        {'_id': ObjectId('62360df6526e120f5cb047ff'), 'userId': 'test',
         'paymentTaskId': ['62360df57e70ed11a0b2109b'], 'amount': 31.61, 'timeline': 3.0,
         'paymentFreq': 'PAYMENT_FREQUENCY_MONTHLY', 'amountPerPayment': 7.91,
         'planType': 'PLAN_TYPE_OPTIM_CREDIT_SCORE', 'endDate': '2022-05-19T00:00:00Z', 'active': True,
         'status': 'PAYMENT_STATUS_IN_DEFAULT', 'paymentAction': [
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.91, 'transactionDate': '2022-03-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_COMPLETED'},
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.91, 'transactionDate': '2022-04-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_IN_DEFAULT'},
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.91, 'transactionDate': '2022-05-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_IN_DEFAULT'},
            {'accountId': '62360c447e70ed11a0b21092', 'amount': 7.879999999999999,
             'transactionDate': '2022-05-19T00:00:00Z', 'status': 'PAYMENT_ACTION_STATUS_IN_DEFAULT'}]},
        {'_id': ObjectId('623613bbd01120921c0bcd0f'), 'userId': 'test',
         'paymentTaskId': ['623613bb8b4b5d9d3e81a1d5'], 'amount': 201.38, 'timeline': 4.0,
         'paymentFreq': 'PAYMENT_FREQUENCY_MONTHLY', 'amountPerPayment': 50.35,
         'planType': 'PLAN_TYPE_OPTIM_CREDIT_SCORE', 'endDate': '2022-07-19T00:00:00Z', 'active': True,
         'status': 'PAYMENT_STATUS_CURRENT', 'paymentAction': [
            {'accountId': '623612478b4b5d9d3e81a1c5', 'amount': 50.35, 'transactionDate': '2022-04-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '623612478b4b5d9d3e81a1c5', 'amount': 50.35, 'transactionDate': '2022-05-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '623612478b4b5d9d3e81a1c5', 'amount': 50.35, 'transactionDate': '2022-06-19T00:00:00Z',
             'status': 'PAYMENT_ACTION_STATUS_PENDING'},
            {'accountId': '623612478b4b5d9d3e81a1c5', 'amount': 50.330000000000005,
             'transactionDate': '2022-07-19T00:00:00Z', 'status': 'PAYMENT_ACTION_STATUS_PENDING'}]},
        ]
    )


def test_overview_amount_per_payment(test_client: TestClient, patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_amount_per_payment')

    response_body = response.json()
    assert response_body['mean'] == 22.05666666666667
    assert response_body['std'] == 20.006407862371386
    assert response_body['median'] == 7.91
    assert response_body['p25'] == 7.91
    assert response_body['p75'] == 29.130000000000003


def test_overview_payment_plan_amounts(test_client: TestClient):
    response = test_client.get('/dashboard/overview_payment_plan_amounts')

    response_body = response.json()
    assert response_body['mean'] == 88.2
    assert response_body['std'] == 80.03034549469345
    assert response_body['median'] == 31.61
    assert response_body['p25'] == 31.61
    assert response_body['p75'] == 116.495


def test_overview_percent_covered_by_plans(test_client: TestClient, patch__fetch_accounts: None,
                                           patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_percent_covered_by_plans')

    response_body = response.json()
    assert response_body['mean'] == 232.99 / 1800
    assert response_body['std'] == 0.0
    assert response_body['median'] == 232.99 / 1800
    assert response_body['p25'] == 232.99 / 1800
    assert response_body['p75'] == 232.99 / 1800


def test_overview_current_vs_defaulted_plans(test_client: TestClient, patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_current_vs_defaulted_plans')

    category_to_count = response.json()['category_to_count']
    assert category_to_count['PAYMENT_STATUS_CURRENT'] == 2
    assert category_to_count['PAYMENT_STATUS_IN_DEFAULT'] == 1


def test_overview_amounts_completed_pending_defaulted(test_client: TestClient, patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_amounts_completed_pending_defaulted')

    category_to_count = response.json()['category_to_count']
    category_to_count['completed'] == 7.91
    category_to_count['pending'] == 232.99
    category_to_count['in_default'] == 23.7


def test_overview_timeline_preferences(test_client: TestClient, patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_timeline_preferences')

    response_body = response.json()
    assert response_body['mean'] == 3.6666666666666665
    assert response_body['std'] == 0.4714045207910317
    assert response_body['median'] == 4.0
    assert response_body['p25'] == 3.5
    assert response_body['p75'] == 4.0


def test_overview_plan_type_preferences(test_client: TestClient, patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_plan_type_preferences')

    category_to_count = response.json()['category_to_count']
    assert category_to_count['PLAN_TYPE_OPTIM_CREDIT_SCORE'] == 2
    assert category_to_count['PLAN_TYPE_MIN_FEES'] == 1


def test_overview_payment_frequency_preferences(test_client: TestClient, patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_payment_frequency_preferences')

    category_to_count = response.json()['category_to_count']
    assert category_to_count['PAYMENT_FREQUENCY_MONTHLY'] == 3


def test_overview_all_preferences(test_client: TestClient, patch_planning_collection_find: None):
    response = test_client.get('/dashboard/overview_all_preferences')
    response_body = response.json()

    overview_timeline_preferences = response_body['overview_timeline_preferences']
    assert overview_timeline_preferences['mean'] == 3.6666666666666665
    assert overview_timeline_preferences['std'] == 0.4714045207910317
    assert overview_timeline_preferences['median'] == 4.0
    assert overview_timeline_preferences['p25'] == 3.5
    assert overview_timeline_preferences['p75'] == 4.0

    overview_plan_type_preferences = response_body['overview_plan_type_option_preferences']['category_to_count']
    assert overview_plan_type_preferences['PLAN_TYPE_OPTIM_CREDIT_SCORE'] == 2
    assert overview_plan_type_preferences['PLAN_TYPE_MIN_FEES'] == 1

    overview_payment_frequency_preferences = response_body['overview_payment_frequency_preferences']['category_to_count']
    assert overview_payment_frequency_preferences['PAYMENT_FREQUENCY_MONTHLY'] == 3

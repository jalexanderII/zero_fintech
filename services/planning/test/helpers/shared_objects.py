from gen.Python.core.accounts_pb2 import Account, AnnualPercentageRates

MOCK_USER_ID = "test"

MOCK_CHASE_ACC = Account(
    account_id="1",
    user_id=MOCK_USER_ID,
    name="Chase",
    available_balance=2500,
    current_balance=500,
    credit_limit=3000,
    annual_percentage_rate=[
        AnnualPercentageRates(apr_percentage=22, apr_type="purchase_apr")
    ],
)

MOCK_AMEX_ACC = Account(
    account_id="2",
    user_id=MOCK_USER_ID,
    name="Amex",
    available_balance=4000,
    current_balance=1000,
    credit_limit=5000,
    annual_percentage_rate=[
        AnnualPercentageRates(apr_percentage=42, apr_type="purchase_apr")
    ],
)

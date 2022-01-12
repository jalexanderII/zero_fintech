CREATE TYPE annual_percentage_rate AS
(
    low_end  float,
    high_end float
);

CREATE TYPE penalty_reason AS ENUM (
    'PENALTY_REASON_UNKNOWN',
    'LATE_PAYMENT'
    );

CREATE TYPE penalty_apr AS
(
    penalty_apr    float,
    penalty_reason penalty_reason
);

CREATE TYPE promotional_rate AS
(
    temporary_apr   float,
    expiration_date timestamp
);

CREATE TABLE IF NOT EXISTS accounts
(
    account_id              SERIAL PRIMARY KEY,
    user_id                 integer,
    name                    varchar(255) NOT NULL,
    created_at              timestamp    NOT NULL DEFAULT NOW(),
    annual_percentage_rate  annual_percentage_rate,
    penalty_apr             penalty_apr,
    due_day                 integer,
    minimum_interest_charge float,
    annual_account_fee      float,
    foreign_transaction_fee float,
    promotional_rate        promotional_rate,
    minimum_payment_due     float,
    current_balance         float,
    pending_transactions    float,
    credit_limit            float
);

CREATE INDEX accounts_name_idx ON accounts (name);

CREATE TYPE transaction_details AS
(
    address           text,
    doing_business_as text,
    date_processed    timestamp
);

CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id      SERIAL PRIMARY KEY,
    user_id             integer,
    account_id          integer      NOT NULL REFERENCES accounts (account_id),
    name                varchar(255) NOT NULL,
    amount              float        NOT NULL,
    date                timestamp,
    rewards_earned      integer,
    transaction_details transaction_details
);

CREATE INDEX transactions_name_idx ON transactions (name);

CREATE TYPE plan_type AS ENUM (
    'PLANTYPE_UNKNOWN',
    'OPTIM_CREDIT_SCORE',
    'MIN_FEES'
    );

CREATE TYPE payment_frequency AS ENUM (
    'PAYMENTFREQ_UNKNOWN',
    'WEEKLY',
    'BIWEEKLY',
    'MONTHLY',
    'QUARTERLY'
    );

CREATE TYPE meta_data AS
(
    preferred_plan_type    plan_type,
    preferred_timeline     float,
    preferred_payment_freq payment_frequency
);

CREATE TABLE IF NOT EXISTS payment_task
(
    payment_task_id SERIAL PRIMARY KEY,
    user_id         integer,
    transaction_id  integer NOT NULL REFERENCES transactions (transaction_id),
    account_id      integer NOT NULL REFERENCES accounts (account_id),
    meta_data       meta_data
);







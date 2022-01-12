-- name: CreatePaymentTask :one
INSERT INTO payment_task (user_id,
                          transaction_id,
                          account_id,
                          meta_data)
VALUES ($1,
        $2,
        $3,
        $4) RETURNING *;

-- name: CreateAccount :one
INSERT INTO accounts (user_id,
                      name,
                      annual_percentage_rate,
                      penalty_apr,
                      due_day,
                      minimum_interest_charge,
                      annual_account_fee,
                      foreign_transaction_fee,
                      promotional_rate,
                      minimum_payment_due,
                      current_balance,
                      pending_transactions,
                      credit_limit)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13) RETURNING *;


-- name: CreateTransaction :one
INSERT INTO transactions (user_id,
                          account_id,
                          name,
                          amount,
                          date,
                          rewards_earned,
                          transaction_details)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7) RETURNING *;
-- name: CreatePaymentTask :one
INSERT INTO payment_task (user_id,
                          transaction_id,
                          account_id,
                          meta_data)
VALUES ($1,
        $2,
        $3,
        $4) RETURNING *;
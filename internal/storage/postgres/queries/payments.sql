-- name: CreatePayment :one
  INSERT INTO payments (
      invoice,
      status,
      amount,
      currency,
      provider,
      provider_payment_id
  )
  VALUES (
      sqlc.arg(invoice),
      sqlc.arg(status),
      sqlc.arg(amount),
      sqlc.arg(currency),
      sqlc.arg(provider),
      sqlc.arg(provider_payment_id)
  )
  RETURNING *;

-- name: GetPaymentByID :one
  SELECT *
  FROM payments
  WHERE id = $1;
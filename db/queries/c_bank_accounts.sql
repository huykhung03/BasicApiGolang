-- name: ListBankAccountsByUsername :many
SELECT * FROM bank_accounts
WHERE username = $1
ORDER BY username;

-- name: CreateBankAccount :one
INSERT INTO bank_accounts (
  username, card_number, currency, balance
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: DeleteBankAccount :exec
DELETE FROM bank_accounts
WHERE card_number = $1;

-- name: GetCardNumberByUserNameAndCurrency :one
SELECT * FROM bank_accounts
WHERE username = $1 AND currency = $2;
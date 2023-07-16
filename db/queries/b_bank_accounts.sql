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

-- name: GetBankAccountByUserNameAndCurrency :one
SELECT * FROM bank_accounts
WHERE username = $1 AND currency = $2;

-- name: GetBankAccountByUserNameAndCurrencyForUpdate :one
SELECT * FROM bank_accounts
WHERE username = $1 AND currency = $2
FOR UPDATE;

-- name: AddBankAccountBalance :one
UPDATE bank_accounts
SET balance = $1
WHERE card_number = $2 AND currency = $3
RETURNING *;
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
SET balance = balance + sqlc.arg(amount)
WHERE card_number = sqlc.arg(card_number) AND currency = $1
RETURNING *;

-- name: GetCurrencyByCardNumberAndUsername :one
SELECT currency FROM bank_accounts
WHERE username = $1 AND card_number = $2;

-- name: ListBankAccounts :many
SELECT * FROM bank_accounts
WHERE username = $1
ORDER BY username;

-- name: DeleteBankAccountByCardNumberAndUserName :exec
DELETE FROM bank_accounts
WHERE card_number = $1 AND username = $2;

-- name: GetCardNumberByCardNumberAndUsername :one
SELECT card_number FROM bank_accounts
WHERE username = $1 AND card_number = $2;

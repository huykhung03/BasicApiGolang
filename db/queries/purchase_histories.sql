-- name: ListPuschaseHistories :many
SELECT * FROM purchase_history
ORDER BY created_at;

-- name: CreatePurchaseHistory :one
INSERT INTO purchase_history (
  id_product, buyer, card_number_of_buyer
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetPurchaseHistory :one
SELECT * FROM purchase_history
WHERE id_purchase_history = $1;

-- name: ListPuschaseHistories :many
SELECT * FROM purchase_history
ORDER BY created_at;

-- name: CreatePuschaseHistory :one
INSERT INTO purchase_history (
  id_product, buyer, card_number
) VALUES (
  $1, $2, $3
) RETURNING *;

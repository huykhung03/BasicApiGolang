-- name: GetProduct :one
SELECT * FROM products
WHERE id_product = $1 
LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id_product;

-- name: CreateProduct :one
INSERT INTO products (
  id_product, product_name, kind_of_product, owner, currency, price
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id_product = $1;
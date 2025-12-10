-- CRUD --
-- name: GetAll :many
SELECT * FROM sales;

-- name: GetById :one
SELECT * FROM sales
WHERE id = $1 LIMIT 1;

-- name: CreateSale :one
INSERT INTO sales (
  type, date, amount, category
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetDateInterval :many
SELECT * FROM sales
WHERE date >= $1 AND date <= $2;

-- name: GetByCategory :many
SELECT * FROM sales
WHERE category = $1;

-- name: UpdateItem :one
UPDATE sales
  set category = $2,
  date = $3,
  amount = $4
WHERE id = $1
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM sales
WHERE id = $1;

-- ANALYTICS --
-- name: GetSum :one
SELECT SUM(amount) FROM sales
WHERE date >= $1 AND date <= $2;

-- name: GetAverage :one
SELECT AVG(amount) FROM sales
WHERE date >= $1 AND date <= $2;

-- name: GetMedian :one
SELECT
percentile_cont(0.5) WITHIN GROUP (ORDER BY amount) AS median
FROM sales
WHERE date >= $1 AND date <= $2;

-- name: GetPercentile :one
SELECT
percentile_cont($1) WITHIN GROUP (ORDER BY amount) AS median
FROM sales
WHERE date >= $2 AND date <= $3;

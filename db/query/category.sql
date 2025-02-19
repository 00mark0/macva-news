-- name: CreateCategory :one
INSERT INTO category (category_name)
VALUES ($1)
RETURNING *;

-- name: GetCategory :one
SELECT *
FROM category
WHERE category_id = $1;

-- name: GetCategoryByName :one
SELECT *
FROM category
WHERE category_name = $1;

-- name: GetCategoryCount :one
SELECT COUNT(*) AS count
FROM category;

-- name: ListCategories :many
SELECT *
FROM category
ORDER BY category_name ASC
LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE category
SET category_name = $2
WHERE category_id = $1
RETURNING *;

-- name: DeleteCategory :one
DELETE FROM category
WHERE category_id = $1
RETURNING *;



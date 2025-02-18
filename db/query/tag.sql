-- name: CreateTag :one
INSERT INTO tag (tag_name)
VALUES ($1)
RETURNING tag_id, tag_name;

-- name: UpdateTag :one
UPDATE tag
SET tag_name = $1
WHERE tag_id = $2
RETURNING tag_id, tag_name;

-- name: DeleteTag :exec
WITH deleted_tag AS (
  DELETE FROM tag
  WHERE tag.tag_id = $1
  RETURNING tag_id
)
DELETE FROM content_tag
WHERE tag_id IN (SELECT tag_id FROM deleted_tag);

-- name: ListTags :many
SELECT tag_id, tag_name
FROM tag
WHERE lower(tag_name) LIKE lower($1)
ORDER BY tag_name ASC;

-- name: AddTagToContent :exec
INSERT INTO content_tag (content_id, tag_id)
VALUES ($1, $2);

-- name: RemoveTagFromContent :exec
DELETE FROM content_tag
WHERE content_id = $1 AND tag_id = $2;

-- name: GetContentByTag :many
SELECT c.content_id, c.user_id, c.category_id, c.title, c.content_description, c.status, c.published_at
FROM content c
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON t.tag_id = ct.tag_id
WHERE lower(t.tag_name) = lower($1)
  AND c.is_deleted = false
ORDER BY c.published_at DESC;


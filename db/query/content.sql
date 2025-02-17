-- name: CreateContent :one
INSERT INTO content (
    user_id,
    category_id,
    title,
    content_description,
    comments_enabled,
    view_count_enabled,
    like_count_enabled,
    dislike_count_enabled
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateContent :one
UPDATE content
SET
    title = COALESCE($2, title),
    content_description = COALESCE($3, content_description),
    category_id = COALESCE($4, category_id),
    comments_enabled = COALESCE($5, comments_enabled),
    view_count_enabled = COALESCE($6, view_count_enabled),
    like_count_enabled = COALESCE($7, like_count_enabled),
    dislike_count_enabled = COALESCE($8, dislike_count_enabled),
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: PublishContent :one
UPDATE content
SET
    status = 'published',
    published_at = now(),
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: SoftDeleteContent :one
UPDATE content
SET
    is_deleted = true,
    updated_at = now()
WHERE content_id = $1
RETURNING *;

-- name: GetContentDetails :one
SELECT
  c.*,
  row_to_json(u) AS author,
  row_to_json(cat) AS category,
  (
    SELECT array_agg(t.tag_name)
    FROM content_tag ct
    JOIN tag t ON ct.tag_id = t.tag_id
    WHERE ct.content_id = c.content_id
  ) AS tags,
  (
    SELECT json_agg(m)
    FROM (
      SELECT media_id, media_type, media_url, media_caption, media_order
      FROM media
      WHERE content_id = c.content_id
      ORDER BY media_order
    ) m
  ) AS media,
  (
    SELECT count(*)
    FROM content_reaction cr
    WHERE cr.content_id = c.content_id
  ) AS reaction_count,
  (
    SELECT count(*)
    FROM comment cm
    WHERE cm.content_id = c.content_id
      AND cm.is_deleted = false
  ) AS comment_count_sync
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.content_id = $1;

-- name: ListPublishedContent :many
SELECT
  c.*,
  row_to_json(u) AS author,
  row_to_json(cat) AS category
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $1 OFFSET $2;

-- name: ListContentByCategory :many
SELECT
  c.*,
  row_to_json(u) AS author
FROM content c
JOIN "user" u ON c.user_id = u.user_id
WHERE c.category_id = $1
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $2 OFFSET $3;

-- name: ListContentByTag :many
SELECT DISTINCT
  c.*,
  row_to_json(u) AS author,
  row_to_json(cat) AS category
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON ct.tag_id = t.tag_id
WHERE t.tag_name = $1
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchContent :many
SELECT DISTINCT
  c.*,
  row_to_json(u) AS author,
  row_to_json(cat) AS category
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
LEFT JOIN content_tag ct ON c.content_id = ct.content_id
LEFT JOIN tag t ON ct.tag_id = t.tag_id
WHERE c.status = 'published'
  AND c.is_deleted = false
  AND (
    c.title ILIKE '%' || $1 || '%'
    OR c.content_description ILIKE '%' || $1 || '%'
    OR t.tag_name ILIKE '%' || $1 || '%'
  )
ORDER BY c.published_at DESC
LIMIT $2 OFFSET $3;










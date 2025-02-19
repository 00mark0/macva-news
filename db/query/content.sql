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

-- name: HardDeleteContent :one
DELETE FROM content
WHERE content_id = $1
RETURNING *;

-- name: GetContentDetails :one
SELECT
  c.*,
  u.username AS author_username,
  cat.category_name,
  (
    SELECT array_agg(t.tag_name)::text[]
    FROM content_tag ct
    JOIN tag t ON ct.tag_id = t.tag_id
    WHERE ct.content_id = c.content_id
  ) AS tags,
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

-- name: GetPublishedContentCount :one
SELECT count(*)
FROM content
WHERE status = 'published'
  AND is_deleted = false;

-- name: ListPublishedContent :many
SELECT
  c.*,
  u.username AS author_username,
  cat.category_name
FROM content c
JOIN "user" u ON c.user_id = u.user_id
JOIN category cat ON c.category_id = cat.category_id
WHERE c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $1 OFFSET $2;

-- name: GetContentByCategoryCount :one
SELECT count(*)
FROM content
WHERE category_id = $1
  AND status = 'published'
  AND is_deleted = false;

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

-- name: GetContentByTagCount :one
SELECT count(DISTINCT c.content_id)
FROM content c
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON ct.tag_id = t.tag_id
WHERE t.tag_name = $1
  AND c.status = 'published'
  AND c.is_deleted = false;

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

-- name: IncrementViewCount :one
UPDATE content
SET
  view_count = view_count + 1
WHERE content_id = $1
RETURNING view_count;

-- name: InsertOrUpdateContentReaction :one
INSERT INTO content_reaction (content_id, user_id, reaction)
VALUES ($1, $2, $3)
ON CONFLICT (content_id, user_id)
DO UPDATE SET reaction = EXCLUDED.reaction
RETURNING content_id;

-- name: DeleteContentReaction :one
DELETE FROM content_reaction
WHERE content_id = $1 AND user_id = $2
RETURNING content_id;

-- name: UpdateContentLikeDislikeCount :one
UPDATE content c
SET
  like_count = (
    SELECT count(*) 
    FROM content_reaction 
    WHERE content_id = c.content_id AND reaction = 'like'
  ),
  dislike_count = (
    SELECT count(*) 
    FROM content_reaction 
    WHERE content_id = c.content_id AND reaction = 'dislike'
  ),
  updated_at = now()
WHERE c.content_id = $1  -- You can use the content_id returned from the previous action here
RETURNING *;

-- name: FetchContentReactions :many
SELECT
  cr.*,
  row_to_json(u) AS user_info
FROM content_reaction cr
JOIN "user" u ON cr.user_id = u.user_id
WHERE cr.content_id = $1;

-- name: GetTrendingContentCount :one
SELECT count(*)
FROM content
WHERE status = 'published'
  AND is_deleted = false
  AND published_at >= $1;

-- name: ListTrendingContent :many
SELECT 
  c.*,
  (c.view_count + c.like_count + c.comment_count) AS total_interactions
FROM content c
WHERE c.status = 'published'
  AND c.is_deleted = false
  AND c.published_at >= $1
ORDER BY total_interactions DESC
LIMIT $2;

-- name: GetContentOverview :one
SELECT 
  COUNT(*) FILTER (WHERE status = 'draft' AND is_deleted = false) AS draft_count,
  COUNT(*) FILTER (WHERE status = 'published' AND is_deleted = false) AS published_count,
  COUNT(*) FILTER (WHERE is_deleted = true) AS deleted_count
FROM content;

-- name: ListContentForModerationCount :one
SELECT COUNT(*) AS count
FROM content c
JOIN "user" u ON c.user_id = u.user_id
WHERE u.banned = true;

-- name: ListContentForModeration :many
SELECT c.*, row_to_json(u) AS author
FROM content c
JOIN "user" u ON c.user_id = u.user_id
WHERE u.banned = true
ORDER BY c.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListRelatedContentByCategoryCount :one
SELECT COUNT(*) AS count
FROM content c
JOIN "user" u ON c.user_id = u.user_id
WHERE c.category_id = $1
  AND c.content_id <> $2
  AND c.status = 'published'
  AND c.is_deleted = false;


-- name: ListRelatedContentByCategory :many
SELECT c.*, row_to_json(u) AS author
FROM content c
JOIN "user" u ON c.user_id = u.user_id
WHERE c.category_id = $1
  AND c.content_id <> $2
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $3;

-- name: ListRelatedContentByTagCount :one
SELECT COUNT(DISTINCT c.content_id) AS count
FROM content c
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON ct.tag_id = t.tag_id
JOIN "user" u ON c.user_id = u.user_id
WHERE t.tag_id = $1
  AND c.content_id <> $2
  AND c.status = 'published'
  AND c.is_deleted = false;

-- name: ListRelatedContentByTag :many
SELECT DISTINCT c.*, row_to_json(u) AS author
FROM content c
JOIN content_tag ct ON c.content_id = ct.content_id
JOIN tag t ON ct.tag_id = t.tag_id
JOIN "user" u ON c.user_id = u.user_id
WHERE t.tag_id = $1
  AND c.content_id <> $2
  AND c.status = 'published'
  AND c.is_deleted = false
ORDER BY c.published_at DESC
LIMIT $3;














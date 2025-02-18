-- name: ListContentComments :many
SELECT
  cm.*,
  row_to_json(u) AS author
FROM comment cm
JOIN "user" u ON cm.user_id = u.user_id
WHERE cm.content_id = $1
  AND cm.is_deleted = false
ORDER BY cm.created_at ASC
LIMIT $2;

-- name: InsertOrUpdateCommentReaction :one
INSERT INTO comment_reaction (comment_id, user_id, reaction)
VALUES ($1, $2, $3)
ON CONFLICT (comment_id, user_id)
DO UPDATE SET reaction = EXCLUDED.reaction
RETURNING comment_id;

-- name: DeleteCommentReaction :one
DELETE FROM comment_reaction
WHERE comment_id = $1 AND user_id = $2
RETURNING comment_id;

-- name: UpdateCommentLikeDislikeCount :one
UPDATE comment c
SET
  like_count = (
    SELECT count(*) 
    FROM comment_reaction 
    WHERE comment_id = c.comment_id AND reaction = 'like'
  ),
  dislike_count = (
    SELECT count(*) 
    FROM comment_reaction 
    WHERE comment_id = c.comment_id AND reaction = 'dislike'
  ),
  updated_at = now()
WHERE c.comment_id = $1
RETURNING *;

-- name: FetchCommentReactions :many
SELECT
  cr.*,
  row_to_json(u) AS user_info
FROM comment_reaction cr
JOIN "user" u ON cr.user_id = u.user_id
WHERE cr.comment_id = $1;


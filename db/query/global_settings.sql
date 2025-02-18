-- name: GetGlobalSettings :one
SELECT * FROM "global_settings" LIMIT 1;

-- name: UpdateGlobalSettings :exec
UPDATE "global_settings"
SET
    "disable_comments" = $1,
    "disable_likes" = $2,
    "disable_dislikes" = $3,
    "disable_views" = $4,
    "disable_ads" = $5,
    "updated_at" = now()
WHERE "global_settings_id" = (SELECT "global_settings_id" FROM "global_settings" LIMIT 1);

-- name: ResetGlobalSettings :exec
UPDATE "global_settings"
SET
    "disable_comments" = false,
    "disable_likes" = false,
    "disable_dislikes" = false,
    "disable_views" = false,
    "disable_ads" = false,
    "updated_at" = now()
WHERE "global_settings_id" = (SELECT "global_settings_id" FROM "global_settings" LIMIT 1);


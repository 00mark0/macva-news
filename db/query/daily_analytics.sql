-- name: GetDailyAnalytics :many
SELECT *
FROM analytics_daily
WHERE analytics_date BETWEEN $1 AND $2
ORDER BY analytics_date DESC
LIMIT $3;

-- name: CreateDailyAnalytics :one
INSERT INTO analytics_daily (
  analytics_date, total_views, total_likes, total_dislikes, total_comments, total_ads_clicks
)
VALUES ($1, 0, 0, 0, 0, 0)
RETURNING *;

-- name: UpdateDailyAnalytics :one
UPDATE analytics_daily
SET 
  total_views = $2,
  total_likes = $3,
  total_dislikes = $4,
  total_comments = $5,
  total_ads_clicks = $6,
  updated_at = now()
WHERE analytics_date = $1
RETURNING *;

-- name: AggregateAnalytics :one
SELECT
    SUM("total_views") AS "total_views",
    SUM("total_likes") AS "total_likes",
    SUM("total_dislikes") AS "total_dislikes",
    SUM("total_comments") AS "total_comments",
    SUM("total_ads_clicks") AS "total_ads_clicks"
FROM "analytics_daily"
WHERE "analytics_date" BETWEEN $1 AND $2;



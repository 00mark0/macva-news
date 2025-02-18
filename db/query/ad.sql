-- name: CreateAd :one
INSERT INTO "ads" 
("title", "description", "image_url", "target_url", "placement", "status", "start_date", "end_date")
VALUES 
($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING "id";

-- name: UpdateAd :one
UPDATE "ads"
SET 
  "title" = $1,
  "description" = $2,
  "image_url" = $3,
  "target_url" = $4,
  "placement" = $5,
  "status" = $6,
  "start_date" = $7,
  "end_date" = $8,
  "updated_at" = now()
WHERE "id" = $9
RETURNING "id";

-- name: DeactivateAd :one
UPDATE "ads"
SET 
  "status" = 'inactive', 
  "updated_at" = now()
WHERE "id" = $1
RETURNING "id";

-- name: DeleteAd :one
DELETE FROM "ads"
WHERE "id" = $1
RETURNING "id";

-- name: ListActiveAds :many
SELECT * 
FROM "ads"
WHERE "status" = 'active'
  AND "start_date" <= now() 
  AND "end_date" >= now();

-- name: ListAdsByPlacement :many
SELECT * 
FROM "ads"
WHERE "placement" = $1
  AND "status" = 'active'
  AND "start_date" <= now() 
  AND "end_date" >= now();

-- name: IncrementAdClicks :one
UPDATE "ads"
SET 
  "clicks" = "clicks" + 1, 
  "updated_at" = now()
WHERE "id" = $1
RETURNING "clicks";


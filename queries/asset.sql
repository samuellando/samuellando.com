-- name: GetAssets :many
SELECT id, name, created
FROM asset;

-- name: GetAsset :one
SELECT sqlc.embed(asset)
FROM asset
WHERE id = $1
LIMIT 1;

-- name: GetAssetByName :one
SELECT sqlc.embed(asset)
FROM asset
WHERE name = $1
LIMIT 1;

-- name: GetAssetContent :one
SELECT content
FROM asset
WHERE id = $1
LIMIT 1;

-- name: CreateAsset :one
INSERT INTO asset (name, content, created)
VALUES ($1, $2, DEFAULT)
ON CONFLICT (name) DO UPDATE
SET content = $2
RETURNING sqlc.embed(asset);

-- name: DeleteAsset :exec
DELETE FROM asset
WHERE id = $1;

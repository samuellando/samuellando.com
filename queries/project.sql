-- name: GetProject :many
SELECT 
    sqlc.embed(p),
    t.id as tag_id,
    t.value as tag_value,
    t.color as tag_color
FROM project p
LEFT JOIN project_tag pt ON pt.project = p.id
LEFT JOIN tag t ON pt.tag = t.id
WHERE p.id = $1
ORDER BY p.id, t.value;

-- name: GetProjects :many
SELECT
    sqlc.embed(p),
    t.id as tag_id,
    t.value as tag_value,
    t.color as tag_color
FROM project p
LEFT JOIN project_tag pt ON pt.project = p.id
LEFT JOIN tag t ON pt.tag = t.id
ORDER BY p.id, t.value;

-- name: UpdateProject :exec
INSERT INTO project (id, description) 
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
SET description = $2;

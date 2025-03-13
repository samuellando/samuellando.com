-- name: GetDocument :many
SELECT 
    sqlc.embed(d),
    t.id as tag_id,
    t.value as tag_value,
    t.color as tag_color
FROM document d
LEFT JOIN document_tag dt ON dt.document = d.id
LEFT JOIN tag t ON dt.tag = t.id
WHERE d.id = $1
ORDER BY d.id, t.value;

-- name: GetDocuments :many
SELECT
    sqlc.embed(d),
    t.id as tag_id,
    t.value as tag_value,
    t.color as tag_color
FROM document d
LEFT JOIN document_tag dt ON dt.document = d.id
LEFT JOIN tag t ON dt.tag = t.id
ORDER BY d.id, t.value;

-- name: CreateDocument :one
INSERT INTO document (title, content, created) VALUES ($1, $2, $3)
RETURNING id;

-- name: UpdateDocument :exec
UPDATE document SET 
    title = $1,
    content = $2,
    created = $3
WHERE 
    id = $4;

-- name: DeleteDocument :exec
DELETE FROM document WHERE id = $1;

-- name: SetDocumentTags :many
WITH clear_document_tags AS (
    DELETE FROM document_tag
    WHERE document_tag.document = $1
    RETURNING document_tag.document
),
tags AS (
    INSERT INTO tag (value)
    SELECT unnest(sqlc.arg(tag_values)::text[])
    ON CONFLICT (value) DO UPDATE
    SET value = tag.value
    RETURNING id, value, color
),
document_tags AS (
    INSERT INTO document_tag (document, tag)
    SELECT $1, tags.id FROM tags
    ON CONFLICT DO NOTHING
    RETURNING document, tag
)
SELECT tags.id, tags.value, tags.color 
FROM tags
ORDER BY value;

-- name: GetAllDocumentTags :many
SELECT DISTINCT sqlc.embed(t)
FROM document_tag dt
INNER JOIN tag t ON dt.tag = t.id;

-- name: GetSharedDocumentTags :many
WITH main_tag AS (
    SELECT tag.id
    FROM tag
    WHERE tag.value = $1
    LIMIT 1
),
documents AS (
    SELECT document
    FROM document_tag, main_tag
    WHERE tag = main_tag.id
),
tags AS (
    SELECT
        tag as id
    FROM document_tag dt, documents
    WHERE dt.document = documents.document
)
SELECT DISTINCT sqlc.embed(t)
FROM tags
INNER JOIN tag t ON t.id = tags.id;

-- name: SetProjectTags :many
WITH clear_project_tags AS (
    DELETE FROM project_tag
    WHERE project_tag.project = $1
    RETURNING project_tag.project
),
tags AS (
    INSERT INTO tag (value)
    SELECT unnest(sqlc.arg(tag_values)::text[])
    ON CONFLICT (value) DO UPDATE
    SET value = tag.value
    RETURNING id, value, color
),
project_tags AS (
    INSERT INTO project_tag (project, tag)
    SELECT $1, tags.id FROM tags
    ON CONFLICT DO NOTHING
    RETURNING project, tag
)
SELECT tags.id, tags.value, tags.color 
FROM tags
ORDER BY value;

-- name: GetAllProjectTags :many
SELECT DISTINCT sqlc.embed(t)
FROM project_tag pt
INNER JOIN tag t ON pt.tag = t.id;

-- name: GetSharedProjectTags :many
WITH main_tag AS (
    SELECT tag.id
    FROM tag
    WHERE tag.value = $1
    LIMIT 1
),
projects AS (
    SELECT project
    FROM project_tag, main_tag
    WHERE tag = main_tag.id
),
tags AS (
    SELECT
        tag as id
    FROM project_tag dt, projects
    WHERE dt.project = projects.project
)
SELECT DISTINCT sqlc.embed(t)
FROM tags
INNER JOIN tag t ON t.id = tags.id;

-- name: GetTags :many
SELECT sqlc.embed(tag)
FROM tag;

-- name: GetTag :one
SELECT sqlc.embed(tag)
FROM tag
WHERE id = $1
LIMIT 1;

-- name: GetTagByValue :one
SELECT sqlc.embed(tag)
FROM tag
WHERE value = $1
LIMIT 1;

-- name: CreateOrUpdateTag :one
INSERT INTO tag (value, color)
VALUES ($1, $2) 
ON CONFLICT (value) DO UPDATE
SET color = $2
RETURNING sqlc.embed(tag);

-- name: DeleteTag :exec
DELETE FROM tag
WHERE id = $1;

-- name: GetCacheByKey :one
SELECT 
    cache_value, valid_to 
FROM cache 
WHERE cache_key = $1;

-- name: SetCacheByKey :exec
INSERT INTO cache (cache_key, cache_value, valid_to) VALUES
    ($1, $2, $3)
ON CONFLICT (cache_key) DO UPDATE 
SET cache_value = EXCLUDED.cache_value,
    valid_to = EXCLUDED.valid_to;

CREATE TABLE IF NOT EXISTS cache (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cache_key text UNIQUE,
    cache_value bytea,
    validto timestamp with time zone
);

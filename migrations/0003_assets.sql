CREATE TABLE IF NOT EXISTS asset (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text UNIQUE,
    created timestamp with time zone DEFAULT now(),
    content bytea
);

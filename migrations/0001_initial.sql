CREATE TABLE IF NOT EXISTS document (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    published BOOL NOT NULL DEFAULT false,
    created timestamp with time zone DEFAULT now()
);
CREATE TABLE IF NOT EXISTS tag (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    value text UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS document_tag (
    document bigint NOT NULL REFERENCES document (id) ON DELETE CASCADE,
    tag bigint NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
    PRIMARY KEY (document, tag)
);

CREATE TABLE IF NOT EXISTS project (
    id bigint NOT NULL PRIMARY KEY,
    description text
);
CREATE TABLE IF NOT EXISTS project_tag (
    project bigint NOT NULL REFERENCES project (id) ON DELETE CASCADE,
    tag bigint NOT NULL REFERENCES tag (id) ON DELETE CASCADE,
    PRIMARY KEY (project, tag)
);

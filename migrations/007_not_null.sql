ALTER TABLE document
ALTER COLUMN created SET NOT NULL;

ALTER TABLE tag
ALTER COLUMN color SET DEFAULT 'white',
ALTER COLUMN color SET NOT NULL;

ALTER TABLE asset
ALTER COLUMN name SET NOT NULL,
ALTER COLUMN created SET NOT NULL;

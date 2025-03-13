ALTER TABLE project
ADD COLUMN image_link text,
ADD COLUMN hidden boolean NOT NULL DEFAULT false;

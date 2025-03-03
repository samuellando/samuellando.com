-- Drop the existing foreign key constraint
ALTER TABLE document_tag DROP CONSTRAINT IF EXISTS document_tag_tag_fkey;

-- Add the new foreign key constraint with ON DELETE RESTRICT
ALTER TABLE document_tag
ADD CONSTRAINT document_tag_tag_fkey
FOREIGN KEY (tag) REFERENCES tag (id) ON DELETE RESTRICT;

-- Drop the existing foreign key constraint
ALTER TABLE project_tag DROP CONSTRAINT IF EXISTS project_tag_tag_fkey;

-- Add the new foreign key constraint with ON DELETE RESTRICT
ALTER TABLE project_tag
ADD CONSTRAINT project_tag_tag_fkey
FOREIGN KEY (tag) REFERENCES tag (id) ON DELETE RESTRICT;

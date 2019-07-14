CREATE TYPE repo_type as ENUM('private', 'public');

ALTER TABLE repos DROP COLUMN type;

ALTER TABLE repos ADD COLUMN type repo_type;

UPDATE repos SET type = 'public'
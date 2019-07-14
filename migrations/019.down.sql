ALTER TABLE repos DROP COLUMN type;

DROP TYPE repo_type;

ALTER TABLE repos ADD COLUMN type VARCHAR(300);

UPDATE repos SET type = 'public'

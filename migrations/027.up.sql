ALTER TABLE labels ADD COLUMN repo_id UUID REFERENCES repos(id) ;

CREATE TABLE issue_labels(
	issue_id BIGINT REFERENCES issues(id),
  label_id BIGINT REFERENCES labels(id)
);

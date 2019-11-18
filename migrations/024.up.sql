create table assignees(
  issue_id BIGINT REFERENCES issues(id),
  user_id  UUID REFERENCES users(id)
);


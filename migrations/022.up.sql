CREATE TYPE read AS ENUM ('read','unread');

CREATE TABLE notifications (
id BIGSERIAL NOT NULL PRIMARY KEY,
read read NOT NULL ,
created_at TIMESTAMP NOT NULL,
issue_id BIGINT REFERENCES issues(id), 
user_id  UUID REFERENCES users(id),
repo_id  UUID REFERENCES repos(id)
);
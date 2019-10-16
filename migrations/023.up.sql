CREATE TABLE logins(
user_id UUID REFERENCES users(id),
secret  VARCHAR(32) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL
);
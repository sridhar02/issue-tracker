CREATE TABLE cookies (
	id UUID NOT NULL PRIMARY KEY,
	user_id UUID REFERENCES users(id),
	created_at TIMESTAMP NOT NULL,
	expires_at TIMESTAMP NOT NULL
);
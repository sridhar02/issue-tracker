CREATE TABLE collaborators (
	repo_id UUID REFERENCES repos(id),
	user_id UUID REFERENCES users(id)
);
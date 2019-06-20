ALTER TABLE users ADD CONSTRAINT unique_email_address UNIQUE(email);

ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE(username);

ALTER TABLE repos ADD CONSTRAINT unique_user_id UNIQUE(user_id,name);

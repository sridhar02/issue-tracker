ALTER TABLE users DROP CONSTRAINT unique_email_address;

ALTER TABLE users DROP CONSTRAINT unique_username;


ALTER TABLE repos DROP CONSTRAINT unique_user_id;
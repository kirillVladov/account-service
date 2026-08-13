ALTER TABLE auth_tokens
DROP COLUMN token_type;

DROP TYPE IF EXISTS token_type;

ALTER TABLE auth_tokens
DROP COLUMN parent_token_id;

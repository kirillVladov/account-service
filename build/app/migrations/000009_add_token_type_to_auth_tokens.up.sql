ALTER TABLE auth_tokens
    ADD COLUMN token_type TEXT NOT NULL DEFAULT 'refresh_token';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_auth_tokens_user_id ON auth_tokens (user_id)
WHERE
    revoked = FALSE;
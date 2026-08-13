CREATE TABLE IF NOT EXISTS account (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid (),
    idempotency_key UUID NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);
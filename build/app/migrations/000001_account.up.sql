CREATE TABLE IF NOT EXISTS account (
    id               UUID           PRIMARY KEY NOT NULL,
    name             TEXT           NOT NULL,
    email            TEXT           NOT NULL UNIQUE,
    password_hash    TEXT           NOT NULL,
    telegram_id      TEXT,
    created_at       TIMESTAMPTZ    NOT NULL,
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
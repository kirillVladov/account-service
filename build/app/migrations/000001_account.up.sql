CREATE TABLE IF NOT EXISTS account (
    id               UUID           PRIMARY KEY NOT NULL,
    -- name             TEXT           NOT NULL,
    -- email            TEXT           NOT NULL UNIQUE,
    -- password_hash    TEXT           NOT NULL,
    -- token            TEXT           NOT NULL,
    telegram_id      TEXT,
    -- telegram_username
    -- refresh_token    TEXT           NOT NULL,
    created_at       TIMESTEMPTZ    NOT NULL,
    updated_at       TIMESTEMPTZ    NOT NULL DEFAULT NOW()
);
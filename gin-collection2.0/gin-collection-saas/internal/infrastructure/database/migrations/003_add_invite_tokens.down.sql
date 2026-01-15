-- Migration: add_invite_tokens (down)
-- Created at: 2026-01-15T11:57:22+01:00

DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS invite_tokens;

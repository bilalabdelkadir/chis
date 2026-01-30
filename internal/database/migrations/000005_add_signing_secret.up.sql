CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE organizations
ADD COLUMN signing_secret TEXT NOT NULL DEFAULT 'whsec_' || encode(gen_random_bytes(32), 'hex');

ALTER TABLE organizations ALTER COLUMN signing_secret DROP DEFAULT;

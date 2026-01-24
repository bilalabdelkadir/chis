ALTER TABLE messages
DROP COLUMN IF EXISTS attempt_count,
DROP COLUMN IF EXISTS next_retry_at;

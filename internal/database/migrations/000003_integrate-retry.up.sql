DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_enum WHERE enumlabel = 'retry' AND enumtypid = 'message_status'::regtype) THEN
        ALTER TYPE message_status ADD VALUE 'retry';
    END IF;
END
$$;
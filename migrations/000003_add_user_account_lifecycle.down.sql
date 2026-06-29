ALTER TABLE users
    DROP COLUMN IF EXISTS deactivated_at;

ALTER TABLE users
    DROP COLUMN IF EXISTS account_status;

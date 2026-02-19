-- Drop indexes
DROP INDEX IF EXISTS idx_accounts_user_is_active;
DROP INDEX IF EXISTS idx_accounts_user_role;
DROP INDEX IF EXISTS idx_accounts_user_email;

-- Drop table
DROP TABLE IF EXISTS accounts_user;

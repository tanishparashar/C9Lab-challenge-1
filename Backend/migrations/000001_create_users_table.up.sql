-- Create accounts_user table
CREATE TABLE IF NOT EXISTS accounts_user (
    id SERIAL PRIMARY KEY,
    password VARCHAR(128) NOT NULL,
    last_login TIMESTAMP WITH TIME ZONE NULL,
    is_superuser BOOLEAN NOT NULL DEFAULT FALSE,
    email VARCHAR(254) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL DEFAULT 'User',
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_staff BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_accounts_user_email ON accounts_user(email);

-- Create index on role for filtering
CREATE INDEX IF NOT EXISTS idx_accounts_user_role ON accounts_user(role);

-- Create index on is_active for filtering
CREATE INDEX IF NOT EXISTS idx_accounts_user_is_active ON accounts_user(is_active);

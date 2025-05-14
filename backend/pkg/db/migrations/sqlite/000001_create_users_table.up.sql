CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    nickname TEXT,
    date_of_birth TEXT, -- ISO 8601 format (e.g., '1990-01-01')
    about_me TEXT,
    avatar_url TEXT, -- Stores path to JPEG/PNG/GIF
    is_private INTEGER NOT NULL DEFAULT 0, -- 0=false, 1=true
    created_at INTEGER NOT NULL, -- Unix timestamp
    updated_at INTEGER NOT NULL,
    deleted_at INTEGER
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
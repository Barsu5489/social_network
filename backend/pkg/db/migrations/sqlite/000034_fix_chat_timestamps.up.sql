-- Create a temporary table with correct schema (matching original)
CREATE TABLE chats_temp (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL CHECK(type IN ('direct', 'group')),
    created_at INTEGER NOT NULL,
    deleted_at INTEGER
);

-- Convert and copy data
INSERT INTO chats_temp (id, type, created_at, deleted_at)
SELECT 
    id, 
    type, 
    CASE 
        WHEN typeof(created_at) = 'text' THEN strftime('%s', created_at)
        ELSE created_at 
    END as created_at,
    deleted_at
FROM chats;

-- Drop old table and rename
DROP TABLE chats;
ALTER TABLE chats_temp RENAME TO chats;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_chats_type ON chats(type);
CREATE INDEX IF NOT EXISTS idx_chats_created_at ON chats(created_at);

CREATE TABLE IF NOT EXISTS group_chats (
    id TEXT PRIMARY KEY,
    group_id TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    deleted_at INTEGER,
    
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    
    UNIQUE (group_id, chat_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_group_chats_group_id ON group_chats(group_id);
CREATE INDEX IF NOT EXISTS idx_group_chats_chat_id ON group_chats(chat_id);
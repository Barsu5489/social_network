CREATE TABLE chats (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL CHECK(type IN ('direct', 'group')),
    created_at INTEGER NOT NULL,
    deleted_at INTEGER
);
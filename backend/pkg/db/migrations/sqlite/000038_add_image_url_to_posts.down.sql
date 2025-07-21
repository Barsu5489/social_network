-- SQLite doesn't support DROP COLUMN directly, so we need to recreate the table
CREATE TABLE posts_temp (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    group_id TEXT,
    content TEXT NOT NULL,
    privacy TEXT NOT NULL CHECK(privacy IN ('public', 'almost_private', 'private')),
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    deleted_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
);

INSERT INTO posts_temp SELECT id, user_id, group_id, content, privacy, created_at, updated_at, deleted_at FROM posts;

DROP TABLE posts;
ALTER TABLE posts_temp RENAME TO posts;

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_group_id ON posts(group_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
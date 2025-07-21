-- SQLite doesn't support DROP COLUMN directly, so we need to recreate the table
-- Create a new table without actor_id
CREATE TABLE notifications_temp (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('follow_request', 'new_follower', 'new_comment', 'group_invite', 'event_created', 'new_like', 'group_join_request', 'group_join_response', 'group_invitation_response', 'new_message')),
    reference_id TEXT NOT NULL,
    is_read INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL,
    deleted_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Copy data from original table (excluding actor_id)
INSERT INTO notifications_temp (id, user_id, type, reference_id, is_read, created_at, deleted_at)
SELECT id, user_id, type, reference_id, is_read, created_at, deleted_at FROM notifications;

-- Drop original table
DROP TABLE notifications;

-- Rename temp table
ALTER TABLE notifications_temp RENAME TO notifications;

-- Recreate indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_reference_id ON notifications(reference_id);
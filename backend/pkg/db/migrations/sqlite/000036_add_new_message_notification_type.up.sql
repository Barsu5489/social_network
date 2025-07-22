-- Add 'new_message' to the notification type constraint
-- First, create a new table with the updated constraint
CREATE TABLE notifications_new (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('follow_request', 'new_follower', 'new_comment', 'group_invite', 'event_created', 'new_like', 'group_join_request', 'group_join_response', 'group_invitation_response', 'new_message')),
    reference_id TEXT NOT NULL,
    is_read INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL,
    deleted_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Copy existing data
INSERT INTO notifications_new SELECT * FROM notifications;

-- Drop old table and rename new one
DROP TABLE notifications;
ALTER TABLE notifications_new RENAME TO notifications;

-- Recreate indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_reference_id ON notifications(reference_id);
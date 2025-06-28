CREATE TABLE notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('follow_request', 'new_follower', 'new_comment', 'group_invite', 'event_created', 'new_like', 'group_join_request', 'group_join_response', 'group_invitation_response')),
    reference_id TEXT NOT NULL, -- Polymorphic UUID, validated in app
    is_read INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL,
    deleted_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_reference_id ON notifications(reference_id);
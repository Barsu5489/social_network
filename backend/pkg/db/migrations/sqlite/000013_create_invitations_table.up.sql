CREATE TABLE invitations (
    id TEXT PRIMARY KEY,
    inviter_id TEXT NOT NULL,
    invitee_id TEXT NOT NULL,
    entity_type TEXT NOT NULL CHECK(entity_type IN ('group', 'event')),
    entity_id TEXT NOT NULL, -- Polymorphic UUID, validated in app
    status TEXT NOT NULL CHECK(status IN ('pending', 'accepted', 'declined')),
    created_at INTEGER NOT NULL,
    deleted_at INTEGER,
    FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (invitee_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (inviter_id, invitee_id, entity_type, entity_id)
);
CREATE INDEX idx_invitations_inviter_id ON invitations(inviter_id);
CREATE INDEX idx_invitations_invitee_id ON invitations(invitee_id);
CREATE INDEX idx_invitations_entity_id ON invitations(entity_id);
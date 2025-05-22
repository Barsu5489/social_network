CREATE TABLE likes (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
     likeable_type TEXT CHECK (likeable_type IN ('post', 'comment')),
    likeable_id UUID,
    created_at INTEGER,
    deleted_at INTEGER,
    CONSTRAINT likes_user_content_unique UNIQUE(user_id, likeable_type, likeable_id)
);
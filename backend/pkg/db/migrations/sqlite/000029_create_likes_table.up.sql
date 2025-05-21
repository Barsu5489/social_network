CREATE TABLE likes (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    likeable_type ENUM('post', 'comment'),
    likeable_id UUID,
    created_at TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT likes_user_content_unique UNIQUE(user_id, likeable_type, likeable_id)
);
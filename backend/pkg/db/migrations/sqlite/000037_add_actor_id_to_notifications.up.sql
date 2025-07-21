-- Add actor_id column to notifications table
ALTER TABLE notifications ADD COLUMN actor_id TEXT;

-- Create index for better performance on actor_id lookups
CREATE INDEX IF NOT EXISTS idx_notifications_actor_id ON notifications(actor_id);

-- Convert existing string timestamps to Unix timestamps in chats table
UPDATE chats 
SET created_at = strftime('%s', created_at) 
WHERE typeof(created_at) = 'text';

-- Convert existing string timestamps to Unix timestamps in chat_participants table
UPDATE chat_participants 
SET joined_at = strftime('%s', joined_at) 
WHERE typeof(joined_at) = 'text';

-- Convert existing string timestamps to Unix timestamps in messages table
UPDATE messages 
SET sent_at = strftime('%s', sent_at) 
WHERE typeof(sent_at) = 'text';
package repository

import "social-nework/pkg/models"

// GetMessageByID retrieves a message by its ID with sender details
func (r *MessageRepository) GetMessageByID(messageID string) (*models.Message, error) {
	query := `
        SELECT m.id, m.chat_id, m.sender_id, m.content, m.sent_at,
               u.first_name, u.last_name, u.avatar_url
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.id = ? AND m.deleted_at IS NULL
    `

	var message models.Message
	var sender models.User

	err := r.DB.QueryRow(query, messageID).Scan(
		&message.ID,
		&message.ChatID,
		&message.SenderID,
		&message.Content,
		&message.SentAt,
		&sender.FirstName,
		&sender.LastName,
		&sender.AvatarURL,
	)

	if err != nil {
		return nil, err
	}

	message.Sender = sender
	return &message, nil
}

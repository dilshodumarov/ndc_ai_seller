package entity

import "time"

type ChatHistory struct {
	GUID             string    `json:"guid"`
	IntegrationID    string    `json:"integration_id"`
	Message          string    `json:"message"`
	ChatID           int64     `json:"chat_id"`
	PlatformID       string    `json:"platform_id"`
	BusinessId       string
	AIResponse       string    `json:"ai_response"`
	ReplyToMessageID *int64    `json:"reply_to_message_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}



package entity

import "time"

type ChatHistory struct {
	GUID             string     `json:"guid"`
	MessageId        int        `json:"message_id"`
	Phone            string     `json:"phone"`
	Message          string     `json:"message"`
	BusinessId       string     `json:"business_id"`
	ChatID           int64      `json:"chat_id"`
	PlatformID       string     `json:"platform_id"`
	AIResponse       string     `json:"ai_response"`
	ReplyToMessageID *int64     `json:"reply_to_message_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}



type SendMessageResponse struct {
	Message          string `json:"message"`
	AIResponse       string `json:"ai_response"`
	BusinessId       string `json:"business_id"`
	From             string `json:"from"`
	Platform         string `json:"platform"`
	Timestamp        string `json:"timestamp"`
	MessageId        int    `json:"message_id"`
	ReplyToMessageID int    `json:"reply_to_message_id"`
	Chatid           int64  `json:"chatid"`
}




type ListChatHistoryRequest struct {
	ChatID     int64  `json:"chat_id"`
	BusinessID string `json:"business_id"`
	Limit      int    `json:"limit"` 
}



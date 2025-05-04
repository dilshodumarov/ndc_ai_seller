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
	Type          string       `json:"type"`
	Notifications *Notification `json:"notifications,omitempty"`
	ChatMessage   *SendMessage  `json:"chat_message,omitempty"`
}

type Notification struct {
	UserId    string `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type SendMessage struct {
	Message          string `json:"message"`
	AIResponse       string `json:"ai_response"`
	UserId           string `json:"user_id"`
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



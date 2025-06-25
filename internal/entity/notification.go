package entity

import (
	"time"
)

type CreateNotificationRequest struct {
	UserID  string `json:"user_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type UpdateNotificationRequest struct {
	GUID    string     `json:"guid"`
	Title   string     `json:"title"`
	Message string     `json:"message"`
	Type    string     `json:"type"`
	IsRead  bool       `json:"is_read"`
	ReadAt  *time.Time `json:"read_at"`
}

type GetNotification struct {
	GUID      string    `json:"guid"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    time.Time `json:"read_at"`
}

type ListNotification struct {
	Notification []GetNotification `json:"notification"`
	Count        int               `json:"count"`
}

type ListNotificationRequest struct {
	UserID string `json:"user_id"`
	Limit  int    `json:"limit"`
	Page   int    `json:"page"` 
}

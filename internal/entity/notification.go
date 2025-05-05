package entity

import (
	"time"
)

type CreateNotificationRequest struct {
	UserID  string
	Title   string
	Message string
	Type    string
}

type UpdateNotificationRequest struct {
	GUID    string
	Title   string
	Message string
	Type    string
	IsRead  bool
	ReadAt  *time.Time
}

type GetNotification struct {
	GUID      string
	UserID    string
	Title     string
	Message   string
	Type      string
	IsRead    bool
	CreatedAt time.Time
	ReadAt    time.Time
}

type ListNotification struct{
	Notification []GetNotification
	Count        int
}

type ListNotificationRequest struct {
	UserID string `json:"user_id"`
	Limit  int    `json:"limit"`
	Page   int    `json:"page"` 
}

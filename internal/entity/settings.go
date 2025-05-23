package entity

import (
	"time"
)

type Settings struct {
	GUID              string
	Name              string
	BrandName         string
	BusinessName      string
	Status            bool
	BusinessID        string
	PromptText        string
	PromptOrder       string
	WaitingTime       int
	PromptProduct     string
	TokenLimit        int
	IntelligenceLevel int
	ErrorMessage      string
	FirstMessage      string
	IsStop            bool
	StopUntil         int
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
	ChatToken         string `json:"chat_token"`
}

type CreateSettingsRequest struct {
	BusinessID string `json:"business_id"`
}

type UpdateSettingsRequest struct {
	GUID              string `json:"guid"`
	Name              string `json:"name"`
	Status            *bool  `json:"status"`
	PromptText        string `json:"prompt_text"`
	WaitingTime       int    `json:"waiting_time"`
	PromptProduct     string `json:"prompt_product"`
	TokenLimit        int    `json:"token_limit"`
	IntelligenceLevel int    `json:"intelligence_level"`
	StopUntil         int    `json:"stop_until"`
	ChatToken         string `json:"chat_token"`
	ChatTokenInt      int
	BrandName         string
	BusinessName      string
	ErrorMessage      string
	FirstMessage      string
	IsStop            *bool
}

type UpdatePromptOrdersRequest struct {
	OrderStatus []UpdateOrderStatusRequest
}

type PromptOrderResponse struct {
	Guid      string `json:"guid"`
	Number    string `json:"number" example:"2"`
	Prompt    string `json:"prompt" example:"Hello there"`
	IsHave    bool   `json:"is_have"`
	PromtJson string `json:"promt_json"`
}

type GetPromptOrdersResponse struct {
	OrderStatus []*OrderStatus `json:"order_status"`
	Id          string         `json:"id"`
}

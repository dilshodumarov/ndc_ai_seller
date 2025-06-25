package entity

import (
	"time"
)

type Settings struct {
	GUID              string     `json:"guid"`
	Name              string     `json:"name"`
	BrandName         string     `json:"brand_name"`
	BusinessName      string     `json:"business_name"`
	Status            bool       `json:"status"`
	BusinessID        string     `json:"business_id"`
	PromptText        string     `json:"prompt_text"`
	PromptOrder       string     `json:"prompt_order"`
	WaitingTime       int        `json:"waiting_time"`
	PromptProduct     string     `json:"prompt_product"`
	TokenLimit        int        `json:"token_limit"`
	IntelligenceLevel int        `json:"intelligence_level"`
	ErrorMessage      string     `json:"error_message"`
	FirstMessage      string     `json:"first_message"`
	IsStop            bool       `json:"is_stop"`
	StopUntil         int        `json:"stop_until"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	ChatToken         string     `json:"chat_token"`
}

type CreateSettingsRequest struct {
	BusinessID string `json:"business_id"`
}

type UpdateSettingsRequest struct {
	GUID              string  `json:"guid"`
	Name              string  `json:"name"`
	Status            *bool   `json:"status"`
	PromptText        string  `json:"prompt_text"`
	WaitingTime       int     `json:"waiting_time"`
	PromptProduct     string  `json:"prompt_product"`
	TokenLimit        int     `json:"token_limit"`
	IntelligenceLevel int     `json:"intelligence_level"`
	StopUntil         int     `json:"stop_until"`
	ChatToken         string  `json:"chat_token"`
	ChatTokenInt      int     `json:"chat_token_int"`
	BrandName         string  `json:"brand_name"`
	BusinessName      string  `json:"business_name"`
	ErrorMessage      string  `json:"error_message"`
	FirstMessage      string  `json:"first_message"`
	IsStop            *bool   `json:"is_stop"`
}

type UpdatePromptOrdersRequest struct {
	OrderStatus []UpdateOrderStatusRequest `json:"order_status"`
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

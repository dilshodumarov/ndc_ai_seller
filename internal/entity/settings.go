package entity

import "time"

type Settings struct {
	GUID              string     `json:"guid"`
	Name              string     `json:"name"`
	Status            bool       `json:"status"`
	BusinessID        string     `json:"business_id"`
	PromptText        string     `json:"prompt_text"`
	PromptOrder       string     `json:"prompt_order"`
	WaitingTime       int        `json:"waiting_time"`
	PromptProduct     string     `json:"prompt_product"`
	TokenLimit        int        `json:"token_limit"`
	IntelligenceLevel int        `json:"intelligence_level"`
	StopUntil         int        `json:"stop_until"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type CreateSettingsRequest struct {
	Name              string `json:"name"`
	BrandName         string `json:"nabrand_nameme"`
	BusinessName      string `json:"business_name"`
	BusinessID        string `json:"business_id"`
	PromptText        string `json:"prompt_text"`
	PromptOrder       string `json:"prompt_order"`
	ErrorMessage      string `json:"error_message"`
	FirstMessage      string `json:"first_message"`
	IsStop            bool   `json:"is_stop"`
	Status            bool   `json:"status"`
	WaitingTime       int    `json:"waiting_time"`
	PromptProduct     string `json:"prompt_product"`
	TokenLimit        int    `json:"token_limit"`
	IntelligenceLevel int    `json:"intelligence_level"`
	StopUntil         int    `json:"stop_until"`
}

type CreateSettingsRequestForSwagger struct {
	Name              string `json:"name"`
	BrandName         string `json:"nabrand_nameme"`
	BusinessName      string `json:"business_name"`
	PromptText        string `json:"prompt_text"`
	PromptOrder       string `json:"prompt_order"`
	ErrorMessage      string `json:"error_message"`
	FirstMessage      string `json:"first_message"`
	IsStop            bool   `json:"is_stop"`
	Status            bool   `json:"status"`
	WaitingTime       int    `json:"waiting_time"`
	PromptProduct     string `json:"prompt_product"`
	TokenLimit        int    `json:"token_limit"`
	IntelligenceLevel int    `json:"intelligence_level"`
	StopUntil         int    `json:"stop_until"`
}

type UpdateSettingsRequest struct {
	GUID              string `json:"guid"`
	Name              string `json:"name"`
	Status            *bool   `json:"status"`
	PromptText        string `json:"prompt_text"`
	PromptOrder       string `json:"prompt_order"`
	WaitingTime       int    `json:"waiting_time"`
	PromptProduct     string `json:"prompt_product"`
	TokenLimit        int    `json:"token_limit"`
	IntelligenceLevel int    `json:"intelligence_level"`
	StopUntil         int    `json:"stop_until"`
}

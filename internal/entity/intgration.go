package entity

import "time"

type IntegrationCreate struct {
	BusinessId       string `json:"business_id"`
	IntegrationToken string `json:"integration_token"`
}


type IntegrationCreateForSwagger struct {
	IntegrationToken string `json:"integration_token"`
}

type IntegrationUpdate struct {
	ID                string `json:"id"`
	Token             string `json:"token"`
}

type IntegrationUpdateForSwagger struct {
	Token             string `json:"token"`
}

type IntegrationUpdateStatus struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

type IntegrationUpdateStatusResponse struct {
	BusinessId       string `json:"business_id"`
	IntegrationToken string `json:"integration_token"`
}

type IntegrationRequest struct {
	BusinessId string `json:"business_id"`
}

type IntegrationListRequest struct {
	BusinessID string    `json:"business_id"`
	SourceType string    `json:"source_type,omitempty"`
	FromDate   time.Time `json:"from_date,omitempty"`
	ToDate     time.Time `json:"to_date,omitempty"`
}

type IntegrationListResponse struct {
	Usages      []TokenUsage `json:"usages"`
	TotalTokens int          `json:"total_tokens"`
}

type TokenUsage struct {
	ID             string    `json:"id"`
	SourceType     string    `json:"source_type"`
	UsedFor        string    `json:"used_for"`
	RequestTokens  int       `json:"request_tokens"`
	ResponseTokens int       `json:"response_tokens"`
	TotalTokens    int       `json:"total_tokens"`
	CreatedAt      time.Time `json:"created_at"`
}

type IntegrationGetResponse struct {
	Guid              string    `json:"guid"`
	Status            string    `json:"status"`
	StartedAt         time.Time `json:"started_at"`
	StoppedAt         time.Time `json:"stopped_at"`
}

type IntegrationUpdateResponse struct {
	GUID  string `json:"guid"`
	Itype string `json:"itype"`
}

type BotIntegration struct {
	Token string `json:"token"`
	Guid  string `json:"guid"`
}

type BotNotification struct {
	Guid      string `json:"guid"`
	ProductId string `json:"product_id"`
}

type BotIntegrationResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}




type IntegrationExistenceResponse struct {
	TelegramAccount bool `json:"telegram_account"`
	TelegramBot     bool `json:"telegram_bot"`
	Instagram       bool `json:"instagram"`
}



type MessageRequest struct {
	Phone  string `json:"phone"`
	UserID string  `json:"user_id"`
	Text   string `json:"text"`
}


type MessageResponse struct {
	Phone  string `json:"phone"`
	Fromid string `json:"fromid"` // `fromid` nomini to'g'ri qilib o'zgartirdim
	Text   string `json:"text"`
	Code   int `json:"code"`  // "Code" ni string sifatida o'zgartirdim
	Message string `json:"message"` // Add a message field for better logging
}
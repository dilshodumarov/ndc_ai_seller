package entity

import "time"

type IntegrationCreate struct {
	BusinessId          string
	IntegrationToken string
	IntegrationType  string
}

type IntegrationCreateForSwagger struct {
	IntegrationToken string
	IntegrationType  string
}
type IntegrationUpdate struct {
	ID                string
	Token             string
	PromptText        string
	PromptOrder        string
	TokenLimit        int
	IntelligenceLevel int
	StopUntil         int
}

type IntegrationUpdateStatus struct {
	Status string
	ID     string
}


type IntegrationUpdateStatusResponse struct {
	IntegrationType string
	BusinessId         string
	IntegrationToken string
}


type IntegrationRequest struct {
	BusinessId string
}

type IntegrationListRequest struct {
	BusinessID  string
	SourceType  string    // optional
	FromDate    time.Time // optional
	ToDate      time.Time // optional
}

type IntegrationListResponse struct {
	Usages      []TokenUsage
	TotalTokens int
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
	Guid string
	IntegrationToken string
	IntegrationType  string
	Status           string
	StartedAt        time.Time
	StoppedAt        time.Time 
	PromptText       string
	PromtOrder       string
	TokenLimit        int
	IntelligenceLevel int
}

type IntegrationUpdateResponse struct {
	GUID    string
	Itype string
}

type BotIntegration struct {
	Token   string
	Guid    string
}

type BotNotification struct {
	Guid string `json:"Guid"`
	ProductId string
}

type BotIntegrationResponse struct {
	Code    int
	Message string
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
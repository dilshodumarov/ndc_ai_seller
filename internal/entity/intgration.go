package entity

import "time"

type IntegrationCreate struct {
	OwnerID          string
	IntegrationToken string
	IntegrationType  string
}

type IntegrationCreateForSwagger struct {
	OwnerID          string
	IntegrationToken string
	IntegrationType  string
}
type IntegrationUpdate struct {
	ID    string
	Token string
	Status string
}

type IntegrationUpdateStatus struct {
	Status string
	ID     string
}


type IntegrationUpdateStatusResponse struct {
	IntegrationType string
	OwnerID         string
	IntegrationToken string
}


type IntegrationRequest struct {
	OwnerID string
}

type IntegrationGetResponse struct {
	Guid string
	IntegrationToken string
	IntegrationType  string
	Status           string
	StartedAt        time.Time
	StoppedAt        time.Time 
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





type PhoneNumber struct {
	Phone string `json:"phone"`
}

type CodeInput struct {
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Password string `json:"password,omitempty"`
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
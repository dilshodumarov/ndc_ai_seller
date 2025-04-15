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

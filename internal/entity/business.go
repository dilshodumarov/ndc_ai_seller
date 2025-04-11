package entity

import "time"

type (
	CreateBusinessRequest struct {
		OwnerID          string
		Name             string
		IntegrationToken string
		IntegrationType  string
		Description      string
	}

	Business struct {
		ID          string    `json:"id"`
		OwnerID     string    `json:"owner_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	UpdateBusinessRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IntegrationType  string
		IntegrationToken string
	}

	UpdateBusinessRequestForSwagger struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IntegrationType  string
		IntegrationToken string
	}

	GetAllBusinessesRequest struct {
		OwnerID string `json:"owner_id"`
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
	}

	GetAllBusinessesResponse struct {
		Itmes []Business `json:"itmes"`
		Total int        `json:"total"`
	}
)

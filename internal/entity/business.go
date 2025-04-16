package entity

import "time"

type (
	CreateBusinessRequest struct {
		OwnerID          string
		Name             string
		Description      string
	}
	CreateBusinessRequestForSwagger struct {
		Name             string
		Description      string
	}
	StartBot struct {
		BusinessId string
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
	}

	UpdateBusinessRequestForSwagger struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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

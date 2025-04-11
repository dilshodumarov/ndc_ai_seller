package entity

import "time"

// CreateClientTypeRequest defines the payload for creating a client type.
type CreateClientTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateClientType struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// ClientTypeResponse defines the response returned for client type operations.
type ClientType struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ClientTypeListResponse defines the response structure for lists (optional)
type ClientTypeListResponse struct {
	Count int64        `json:"count"`
	Items []ClientType `json:"items"`
}

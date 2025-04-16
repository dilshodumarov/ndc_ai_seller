package entity

import "time"

type BotCommandRequest struct {
	IntegrationID string `json:"integration_id"`
	Command       string `json:"command"`
	Response      string `json:"response"`
}

type BotCommandResponse struct {
	Guid          string    `json:"guid"`
	IntegrationID string    `json:"integration_id"`
	Command       string    `json:"command"`
	Response      string    `json:"response"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}


type BotCommandUpdateRequest struct {
	Guid          string    `json:"guid"`
	Command       string    `json:"command"`
	Response      string    `json:"response"`
}


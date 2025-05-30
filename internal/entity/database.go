package entity

import "time"

type Database struct {
	Guid        string    `json:"guid"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Tokens      *int      `json:"tokens,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Databaselist struct {
	Items []Database `json:"items"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Count int        `json:"count"`
}

type Filter struct {
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type CreateDatabaseRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Tokens      *int    `json:"tokens"`
}

type UpdateDatabaseRequest struct {
	Guid        string  `json:"guid"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Tokens      *int    `json:"tokens"`
}

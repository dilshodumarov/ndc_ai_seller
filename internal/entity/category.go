package entity

import "time"

type (
	//Category

	CreateCategoryRequest struct {
		Name       string `json:"name"`
	}
	CreateCategoryRequestForSwagger struct {
		Name string `json:"name"`
	}
	Category struct {
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	UpdateCategoryRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	GetAllCategoriesRequest struct {
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
	}

	GetAllCategoriesResponse struct {
		Items []Category `json:"items"`
		Total uint64     `json:"total"`
	}
	CategoryFilter struct {
		Name       string `json:"name,omitempty"`
		Page       uint64 `json:"page,omitempty"`
		Limit      uint64 `json:"limit,omitempty"`
	}
)

package entity

import "time"

type (
	//Product

	CreateProductRequest struct {
		BusinessID   string `json:"business_id"`
		Name         string `json:"name"`
		CategoryID   string `json:"category_id"`
		ShortInfo    string `json:"short_info"`
		Description  string `json:"description"`
		Cost         int    `json:"cost"`
		Count        int    `json:"count"`
		DiscountCost int    `json:"discount_cost"`
		Discount     int    `json:"discount"`
	}

	Product struct {
		ID           string    `json:"id"`
		BusinessID   string    `json:"business_id"`
		Name         string    `json:"name"`
		CategoryID   string    `json:"category_id"`
		ShortInfo    string    `json:"short_info"`
		Description  string    `json:"description"`
		Cost         int       `json:"cost"`
		Count        int       `json:"count"`
		DiscountCost int       `json:"discount_cost"`
		Discount     int       `json:"discount"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	UpdateProductRequest struct {
		ID           string `json:"id"`
		BusinessID   string `json:"business_id"`
		Name         string `json:"name"`
		CategoryID   string `json:"category_id"`
		ShortInfo    string `json:"short_info"`
		Description  string `json:"description"`
		Cost         int    `json:"cost"`
		Count        int    `json:"count"`
		DiscountCost int    `json:"discount_cost"`
		Discount     int    `json:"discount"`
	}

	GetAllProductsRequest struct {
		BusinessID string `json:"business_id"`
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
	}

	GetAllProductsResponse struct {
		Items []Product `json:"items"`
		Total uint64    `json:"total"`
	}
)

type (
	//Category

	CreateCategoryRequest struct {
		BusinessID string `json:"business_id"`
		Name       string `json:"name"`
	}
	Category struct {
		ID         string    `json:"id"`
		BusinessID string    `json:"business_id"`
		Name       string    `json:"name"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	UpdateCategoryRequest struct {
		ID         string `json:"id"`
		BusinessID string `json:"business_id"`
		Name       string `json:"name"`
	}

	GetAllCategoriesRequest struct {
		BusinessID string `json:"business_id"`
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
	}

	GetAllCategoriesResponse struct {
		Items []Category `json:"items"`
		Total uint64     `json:"total"`
	}
)


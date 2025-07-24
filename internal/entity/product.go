package entity

import "time"

type (
	//Product

	CreateProductRequest struct {
		BusinessID   string   `json:"business_id"`
		Name         string   `json:"name"`
		CategoryID   *string  `json:"category_id"`
		ShortInfo    string   `json:"short_info"`
		Description  string   `json:"description"`
		Cost         int      `json:"cost"`
		Count        int      `json:"count"`
		DiscountCost int      `json:"discount_cost"`
		Discount     int      `json:"discount"`
		Status       bool     `json:"status"`
		Image_url    []string `json:"image_url"`
	}
	CreateProductImage struct {
		ProductId string `json:"product_id"`
		ImageUrl  string `json:"image_url"`
	}

	CreateProductRequestForSwagger struct {
		Name         string   `json:"name"`
		CategoryID   *string  `json:"category_id"`
		ShortInfo    string   `json:"short_info"`
		Description  string   `json:"description"`
		Cost         int      `json:"cost"`
		Count        int      `json:"count"`
		DiscountCost int      `json:"discount_cost"`
		Discount     int      `json:"discount"`
		Status       bool     `json:"status"`
		Image_url    []string `json:"image_url"`
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
		Image_urls   []string  `json:"image_url"`
		ProductId    int       `json:"product_id"`
		Status       bool      `json:"status"`
		CategoryName string    `json:"category"`
	}

	UpdateProductRequest struct {
		ID           string   `json:"id"`
		Name         string   `json:"name"`
		CategoryID   string   `json:"category_id"`
		ShortInfo    string   `json:"short_info"`
		Description  string   `json:"description"`
		Cost         int      `json:"cost"`
		Count        int      `json:"count"`
		DiscountCost int      `json:"discount_cost"`
		Discount     int      `json:"discount"`
		Status       *bool    `json:"status"`
		Image_url    []string `json:"image_url"`
	}

	UpdateProductRequestForSwagger struct {
		Name         string   `json:"name"`
		CategoryID   string   `json:"category_id"`
		ShortInfo    string   `json:"short_info"`
		Description  string   `json:"description"`
		Cost         int      `json:"cost"`
		Count        int      `json:"count"`
		DiscountCost int      `json:"discount_cost"`
		Discount     int      `json:"discount"`
		Status       *bool    `json:"status"`
		Image_url    []string `json:"image_url"`
	}

	GetAllProductsRequest struct {
		BusinessID string `json:"business_id"`
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
	}

	GetAllProductsResponse struct {
		Items      []Product `json:"items"`
		Count      uint64    `json:"count"`
		TotalCount uint64    `json:"totalcount"`
	}

	ProductFilter struct {
		OwnerID      string `json:"owner_id"`
		CategoryID   string `json:"category_id"`
		Search       string `json:"search"`
		Limit        uint64 `json:"limit"`
		Page         uint64 `json:"page"`
		ProductId    int    `json:"product_id"`
		ProductCount int    `json:"product_count"`
		Status       string `json:"status"`
	}
)

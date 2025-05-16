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
		Image_url    []string
	}
	CreateProductImage struct {
		ProductId string
		ImageUrl  string
	}
	CreateProductRequestForSwagger struct {
		Name         string `json:"name"`
		CategoryID   string `json:"category_id"`
		ShortInfo    string `json:"short_info"`
		Description  string `json:"description"`
		Cost         int    `json:"cost"`
		Count        int    `json:"count"`
		DiscountCost int    `json:"discount_cost"`
		Discount     int    `json:"discount"`
		Image_url    []string
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
		Image_urls   []string
		ProductId    int       `json:"product_id"`
		Status       string    `json:"status"`
		CategoryName string    `json:"category"`
	}

	UpdateProductRequest struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		CategoryID   string `json:"category_id"`
		ShortInfo    string `json:"short_info"`
		Description  string `json:"description"`
		Cost         int    `json:"cost"`
		Count        int    `json:"count"`
		DiscountCost int    `json:"discount_cost"`
		Discount     int    `json:"discount"`
		Image_url    []string
	}

	UpdateProductRequestForSwagger struct {
		Name         string `json:"name"`
		CategoryID   string `json:"category_id"`
		ShortInfo    string `json:"short_info"`
		Description  string `json:"description"`
		Cost         int    `json:"cost"`
		Count        int    `json:"count"`
		DiscountCost int    `json:"discount_cost"`
		Discount     int    `json:"discount"`
		Image_url    []string
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
	 ProductFilter struct {
		OwnerID    string  // majburiy
		CategoryID string  // optional
		Search     string  // optional
		Limit      uint64  // required
		Page       uint64  // required
		ProductId  int
		ProductCount int
		Status     string
	}
	
)


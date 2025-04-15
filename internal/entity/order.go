package entity

import "time"

type (
	// Attribute

	Attribute struct {
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		CategoryID string    `json:"category_id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
)

type (
	//Order

	CreateOrderRequest struct {
		BusinessID    string `json:"business_id"`
		ClientID      string `json:"client_id"`
		IntegrationID string `json:"integration_id"`
		Status        string `json:"status"`
		TotalCost     int    `json:"total_cost"`
	}
	Order struct {
		ID                string    `json:"id"`
		BusinessID        string    `json:"business_id"`
		ClientID          string    `json:"client_id"`
		Status            string    `json:"status"`
		StatusChangedTime *time.Time    `json:"status_chaged_time"`
		TotalCost         int       `json:"total_cost"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}

	OrderUpdate struct {
		ID                string    `json:"id"`
		ClientID          string    `json:"client_id"`
		IntegrationID     string    `json:"integration_id"`
		Status            string    `json:"status"`
	}

	GetAllOrdersResponse struct {
		Items []Order `json:"items"`
		Total uint64  `json:"total"`
	}
)

type (
	// OrderProducts

	CreateOrderProductsRequest struct {
		BusinessID string `json:"business_id"`
		OrderID    string `json:"order_id"`
		ProductID  string `json:"product_id"`
		Count      int    `json:"count"`
		Cost       int    `json:"cost"`
	}
	UpdateOrderProductsRequest struct {
		ID         string `json:"id"`
		BusinessID string `json:"business_id"`
		OrderID    string `json:"order_id"`
		ProductID  string `json:"product_id"`
		Count      int    `json:"count"`
		Cost       int    `json:"cost"`
	}
	OrderProducts struct {
		ID         string    `json:"id"`
		BusinessID string    `json:"business_id"`
		OrderID    string    `json:"order_id"`
		ProductID  string    `json:"product_id"`
		Count      int       `json:"count"`
		Cost       int       `json:"cost"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
)

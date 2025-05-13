package entity

import "time"



type (
	// Order

	CreateOrderRequest struct {
		ClientID      string               `json:"client_id"`
		BusinessID    string               `json:"business_id"`
		LocationURL   string               `json:"location_url"`
		Status        string               `json:"status"`
		TotalPrice    float64              `json:"total_price"`
		PaymentMethod string               `json:"payment_method"`
		Products      []CreateOrderProduct `json:"products"`
	}

	CreateOrderProduct struct {
		ProductID string  `json:"product_id"`
		Count     int     `json:"count"`
		Price     float64 `json:"price"`
	}

	Order struct {
		ID                string
		OrderId           string
		Client          ClientInfo
		BusinessID        string
		LocationURL       string
		Status            string
		AdminStatus       string
		TotalPrice        float64
		PaymentMethod     string
		Platform          string
		ImageUrl          string
		StatusChangedTime *time.Time
		CreatedAt         time.Time
		UpdatedAt         time.Time
		Products          []OrderProduct // <--- Products field qo‘shamiz
	}

	ClientInfo struct {
		GUID  string
		Name  string
		Phone string
	}
	
	OrderProduct struct {
		ProductID        string
		Name             string
		ImageURL         string
		Cost             int
		Count            int
		ProductTotalPrice float64 // order_products.price
	}
	

	OrderUpdate struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		LocationURL   string `json:"location_url"`
		PaymentMethod string `json:"payment_method"`
	}
	OrderUpdateForSwagger struct {
		Status        string `json:"status"`
		LocationURL   string `json:"location_url"`
		PaymentMethod string `json:"payment_method"`
	}

	GetAllOrdersResponse struct {
		Items []Order `json:"items"`
		Total uint64  `json:"total"`
	}

	OrderFilter struct {
		ID            string
		ClientID      string
		BusinessID    string
		Status        string
		PaymentMethod string
		Platform      string
		Search        string
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

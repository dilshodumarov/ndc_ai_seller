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
		ID                string         `json:"id"`
		OrderId           string         `json:"order_id"`
		Client            ClientInfo     `json:"client"`
		BusinessID        string         `json:"business_id"`
		StatusNumber      int            `json:"status_number"`
		LocationURL       string         `json:"location_url"`
		Status            string         `json:"status"`
		AdminStatus       string         `json:"admin_status"`
		TotalPrice        float64        `json:"total_price"`
		PaymentMethod     string         `json:"payment_method"`
		Platform          string         `json:"platform"`
		ImageUrl          string         `json:"image_url"`
		StatusChangedTime *time.Time     `json:"status_changed_time"`
		CreatedAt         time.Time      `json:"created_at"`
		UpdatedAt         time.Time      `json:"updated_at"`
		Location          string         `json:"location"`
		Description       string         `json:"description"`
		Products          []OrderProduct `json:"products"`
	}

	ClientInfo struct {
		GUID     string `json:"guid"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		UserName string `json:"user_name"`
	}

	OrderProduct struct {
		ProductID         string   `json:"product_id"`
		Name              string   `json:"name"`
		ImageURL          []string `json:"image_urls"`
		Cost              int      `json:"cost"`
		Count             int      `json:"count"`
		ProductTotalPrice int      `json:"product_total_price"`
	}

	OrderUpdate struct {
		ID            string  `json:"id"`
		Status        string  `json:"status"`
		LocationURL   string  `json:"location_url"`
		PaymentMethod string  `json:"payment_method"`
		BussnesId     string  `json:"bussnes_id"`
		StatusNumber  int     `json:"status_number"`
		StatusID      *string `json:"status_id"`
	}

	OrderUpdateForSwagger struct {
		Status string `json:"status"`
	}

	GetAllOrdersResponse struct {
		Items []Order `json:"items"`
		Total uint64  `json:"total"`
	}

	OrderFilter struct {
		ID            string `json:"id"`
		ClientID      string `json:"client_id"`
		BusinessID    string `json:"business_id"`
		Status        string `json:"status"`
		PaymentMethod string `json:"payment_method"`
		Platform      string `json:"platform"`
		Search        string `json:"search"`
		Daye          int    `json:"daye"`
	}

	OrderProductBuOrderID struct {
		ProductID             string    `json:"product_id"`
		Name                  string    `json:"name"`
		ImageURL              string    `json:"image_url"`
		Cost                  int       `json:"cost"`
		Status                bool      `json:"status"`
		Discount              int       `json:"discount"`
		DiscountCost          int       `json:"discount_cost"`
		ShortInfo             string    `json:"short_info"`
		Description           string    `json:"description"`
		CreatedAt             time.Time `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
		Count                 int       `json:"count"`
		Price                 float64   `json:"price"`
		ProductTotalPrice     float64   `json:"product_total_price"`
		OrderProductCreatedAt time.Time `json:"order_product_created_at"`
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

type OrderStatus struct {
	GUID         string              `json:"guid"`
	BusinessID   string              `json:"business_id"`
	TypeID       string              `json:"type_id"`
	CustomName   string              `json:"custom_name"`
	TypeName     string              `json:"type_name"`
	StatusNumber int                 `json:"status_number"`
	FonColor     string              `json:"fon_color"`
	OrderCount   int                 `json:"order_count"`
	CreatedAt    time.Time           `json:"created_at"`
	Prompts      PromptOrderResponse `json:"prompts"`
}

type OrderStatusFilter struct {
	BusinessID string `json:"business_id"`
	Status     string `json:"status"`
	Days       int    `json:"daye"`
}

type CreateOrderStatusRequest struct {
	BusinessID string `json:"business_id"`
	TypeID     string `json:"type_id"`
	CustomName string `json:"custom_name"`
}

type CreateOrderStatusRequestForswagger struct {
	TypeID     string `json:"type_id"`
	CustomName string `json:"custom_name"`
}

type UpdateOrderStatusRequest struct {
	GUID        string `json:"guid"`
	CustomName  string `json:"custom_name"`
	PromtNumber int    `json:"promt_number"`
	Promt       string `json:"promt"`
	FonColor    string `json:"fon_color"`
}

package entity

import (
	"time"
)

// Authentication request and response structures
type (
	// RegisterRequest represents the user registration request payload
	RegisterRequest struct {
		FirstName   string `json:"first_name" binding:"required,min=2,max=30"`
		LastName    string `json:"last_name"`
		Email       string `json:"email" binding:"required,email"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password" binding:"required,min=6,max=16"`
	}

	// AuthResponse represents the response after successful authentication
	User struct {
		ID           string    `json:"id"`
		FirstName    string    `json:"first_name"`
		LastName     string    `json:"last_name"`
		Email        string    `json:"email"`
		PhoneNumber  string    `json:"phone_number"`
		Password     string    `json:"-"`
		IsActive     *bool     `json:"is_active"`
		RoleID       string    `json:"role_id"`
		RoleData     Role      `json:"role_data"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token,omitempty"`
		BusinessID   string    `json:"business_id,omitempty"`
	}
	UserFilter struct {
		IsActive  *bool
		RoleID    string
		CreatedAt string
		Offset    uint64
		Limit     uint64
	}

	ListUsers struct {
		Items []User `json:"itmes"`
		Total int    `json:"total"`
	}

	// LoginRequest represents the login request payload
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// VerifyEmail represents the email verification request
	VerifyEmail struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	// VerifyRequest represents the request to verify a user with a code
	VerifyRequest struct {
		Email       string `json:"email" binding:"required,email"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	// ForgotPasswordRequest represents the forgot password request
	ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email"`
	}

	// UpdatePasswordRequest represents the update password request
	UpdatePasswordRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// AccessToken represents a token for API access
	AccessToken struct {
		AccessToken string `json:"access_token"`
	}

	// JWTClaims represents the JWT token claims
	JWTClaims struct {
		Sub        string `json:"sub"`
		Role       string `json:"role"`
		Exp        int64  `json:"exp"`
		BusinessId string `json:"bussnes_id"`
	}

	// UpdatePassword represents a request to update a password
	UpdatePassword struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

// entity.ClientFilter
type ClientFilter struct {
	Name        string `json:"name"`
	ClientId    int `json:"client_id"`
	Phone       string `json:"phone"`
	From        string `json:"from"`
	Goal        string `json:"goal"`
	OrderStatus string `json:"order_status"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

// entity.Client
type Client struct {
	ID          string    `json:"id"`
	PlatformID  string    `json:"platform_id"`
	ClientId    int    `json:"client_id"`
	FirstName   string    `json:"first_name"`
	Phone       string    `json:"phone"`
	CreatedAt   time.Time `json:"created_at"`
	UserName    string    `json:"user_name"`
	From        string    `json:"from"`
	OrderStatus string    `json:"order_status"`
	Goal        string    `json:"goal"`
	IsBlock     bool      `json:"is_block"`
	Location    string    `json:"location"`
}

type ListClients struct {
	Items   []Client `json:"items"`
	Count   int      `json:"count"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}


type BlockUser struct {
	PlatformId string `json:"platform_id"`
	BusinessID string `json:"business_id"`
	Block      bool   `json:"block"`
}


type PauzeChat struct {
	PlatformId string `json:"platform_id"`
	BusinessID string `json:"business_id"`
	Pauze      bool   `json:"pauze"`
	Type       string `json:"type"`
}

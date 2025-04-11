package entity

import "time"

// CreateRoleRequest is the structure for creating a new role.
type CreateRoleRequest struct {
	Name         string `json:"name" binding:"required"`
	ClientTypeId string `json:"client_type_id" binding:"required"`
}

type UpdateRoleRequest struct {
	ID           string `json:"id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	ClientTypeId string `json:"client_type_id" binding:"required"`
}

// RoleResponse is the structure for sending the role data in response.
type Role struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	ClientTypeId    string     `json:"client_type_id"`
	ClientTypeData ClientType `json:"client_type_data"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type RoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ClientTypeId string    `json:"client_type_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleListResponse struct {
	Items []RoleResponse `json:"items"`
	Count uint64          `json:"count"`
}

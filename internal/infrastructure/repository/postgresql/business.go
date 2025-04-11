package postgresql

import (
	"context"
	"fmt"
	"strings"
	"sugurta/internal/pkg/postgres"
	"time"

	"sugurta/internal/entity"
)

const (
	businessTableName      = "business"
	businessServiceName    = "businessService"
	businessSpanRepoPrefix = "businessRepo"
)

type businessRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewBusinessRepo(db *postgres.Postgres) *businessRepo {
	return &businessRepo{
		tableName: businessTableName,
		db:        db,
	}
}

// CreateBusiness creates a new business record.
func (p *businessRepo) Create(ctx context.Context, b *entity.CreateBusinessRequest) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (owner_id, name, integration_token, integration_type, description) 
		 VALUES ($1, $2, $3, $4, $5)`,
		p.tableName,
	)

	_, err := p.db.Exec(ctx, query,
		b.OwnerID,
		b.Name,
		b.IntegrationToken,
		b.IntegrationType,
		b.Description,
	)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

// GetBusiness retrieves a business record by ID.
func (p *businessRepo) Get(ctx context.Context, id string) (*entity.Business, error) {
	var business entity.Business

	query := fmt.Sprintf(
		"SELECT id, owner_id, name, description, created_at, updated_at FROM %s WHERE id = $1 AND deleted_at IS NULL",
		p.tableName,
	)

	err := p.db.QueryRow(ctx, query, id).Scan(
		&business.ID,
		&business.OwnerID,
		&business.Name,
		&business.Description,
		&business.CreatedAt,
		&business.UpdatedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return &business, nil
}

// List retrieves a list of businesses with filtering and pagination.
func (p *businessRepo) List(ctx context.Context, req entity.GetAllBusinessesRequest) (*entity.GetAllBusinessesResponse, error) {
	// Default to pagination values if not provided
	offset := (req.Page - 1) * req.Limit

	// Build the query for businesses with filtering and pagination
	query := fmt.Sprintf(
		"SELECT id, owner_id, name, description, created_at, updated_at FROM %s WHERE deleted_at IS NULL",
		p.tableName,
	)

	// Apply filtering by OwnerID if provided
	if req.OwnerID != "" {
		query += " AND owner_id = $1"
	}

	// Apply pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(req.OwnerID)+1, len(req.OwnerID)+2)

	// Prepare the arguments for pagination and filtering
	args := []interface{}{req.OwnerID, req.Limit, offset}

	// Total count query
	countQuery := fmt.Sprintf(
		"SELECT COUNT(1) FROM %s WHERE deleted_at IS NULL",
		p.tableName,
	)

	if req.OwnerID != "" {
		countQuery += " AND owner_id = $1"
	}

	// Get total count
	var total int
	err := p.db.QueryRow(ctx, countQuery, req.OwnerID).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Execute the main query
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var businesses entity.GetAllBusinessesResponse
	for rows.Next() {
		var business entity.Business
		if err = rows.Scan(
			&business.ID,
			&business.OwnerID,
			&business.Name,
			&business.Description,
			&business.CreatedAt,
			&business.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		businesses.Itmes = append(businesses.Itmes, business)
	}

	// Set total count
	businesses.Total = total

	return &businesses, nil
}

// UpdateBusiness updates an existing business record.
func (p *businessRepo) Update(ctx context.Context, b *entity.UpdateBusinessRequest) error {
	setClauses := []string{}
	args := []interface{}{}
	argID := 1

	if b.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argID))
		args = append(args, b.Name)
		argID++
	}

	if b.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argID))
		args = append(args, b.Description)
		argID++
	}

	if b.IntegrationType != "" {
		setClauses = append(setClauses, fmt.Sprintf("integration_type = $%d", argID))
		args = append(args, b.IntegrationType)
		argID++
	}

	if b.IntegrationToken != "" {
		setClauses = append(setClauses, fmt.Sprintf("integration_token = $%d", argID))
		args = append(args, b.IntegrationToken)
		argID++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now().UTC())
	argID++

	args = append(args, b.ID)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE guid = $%d", p.tableName, strings.Join(setClauses, ", "), argID)

	_, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

// DeleteBusiness deletes a business record by ID.
func (p *businessRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(
		"UPDATE %s SET deleted_at = $1 WHERE id = $2",
		p.tableName,
	)

	_, err := p.db.Exec(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

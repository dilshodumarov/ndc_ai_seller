package postgresql

import (
	"context"
	"fmt"
	"time"

	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
)

const (
	categoryTableName = "category"
)

// categoryRepo -.
type categoryRepo struct {
	tableName string
	db        *postgres.Postgres
}

// NewCategoryRepo -.
func NewCategoryRepo(db *postgres.Postgres) *categoryRepo {
	return &categoryRepo{
		tableName: categoryTableName,
		db:        db,
	}
}

// Create -.
func (r *categoryRepo) Create(ctx context.Context, category *entity.CreateCategoryRequest) error {
	query := `
		INSERT INTO category (name, business_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	now := time.Now()
	_, err := r.db.Exec(ctx, query, category.Name, category.BusinessID, now, now)
	if err != nil {
		return r.db.Error(err)
	}
	return nil
}

// Get -.
func (r *categoryRepo) Get(ctx context.Context, id string) (*entity.Category, error) {
	var category entity.Category

	query := `
		SELECT guid, name, business_id, created_at, updated_at
		FROM category
		WHERE guid = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&category.ID,
		&category.Name,
		&category.BusinessID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("categoryRepo - Get - row.Scan: %w", err)
	}

	return &category, nil
}

// List -.
func (r *categoryRepo) List(ctx context.Context, filter entity.CategoryFilter) (*entity.GetAllCategoriesResponse, error) {
	var (
		categories  []entity.Category
		args        []interface{}
		whereClause = " WHERE true"
	)
	argIndex := 1
	
	if filter.BusinessID != "" {
		whereClause += fmt.Sprintf(" AND business_id = $%d", argIndex)
		args = append(args, filter.BusinessID)
		argIndex++
	}
	
	if filter.Name != "" {
		whereClause += fmt.Sprintf(" AND name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.Name+"%")
		argIndex++
	}
	
	// Default values
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.Limit == 0 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit
	
	baseQuery := `
		SELECT guid, name, business_id, created_at, updated_at
		FROM category
	`
	countQuery := `SELECT COUNT(*) FROM category`
	
	limitOffset := fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, offset)
	
	finalQuery := baseQuery + whereClause + limitOffset
	rows, err := r.db.Query(ctx, finalQuery, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()
	
	categories = make([]entity.Category, 0)
	for rows.Next() {
		var category entity.Category
		err = rows.Scan(
			&category.ID,
			&category.Name,
			&category.BusinessID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, r.db.Error(err)
		}
		categories = append(categories, category)
	}
	
	// Count query
	countQueryFinal := countQuery + whereClause
	var totalCount uint64
	err = r.db.QueryRow(ctx, countQueryFinal, args...).Scan(&totalCount)
	if err != nil {
		return nil, r.db.Error(err)
	}
	
	return &entity.GetAllCategoriesResponse{
		Items: categories,
		Total: totalCount,
	}, nil
}	

// Update -.
func (r *categoryRepo) Update(ctx context.Context, category *entity.UpdateCategoryRequest) error {
	query := `
		UPDATE category
		SET name = $1, updated_at = $2
		WHERE guid = $3
	`
	res, err := r.db.Exec(ctx, query, category.Name, time.Now(), category.ID)
	if err != nil {
		return r.db.Error(err)
	}
	if res.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no rows affected"))
	}
	return nil
}

// Delete -.
func (r *categoryRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM category WHERE guid = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.db.Error(err)
	}
	if res.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no rows affected"))
	}
	return nil
}

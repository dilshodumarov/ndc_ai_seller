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
		SELECT id, name, business_id, created_at, updated_at
		FROM category
		WHERE id = $1
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
func (r *categoryRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.GetAllCategoriesResponse, error) {
	var (
		categories []entity.Category
		args       []interface{}
		where      string
	)

	baseQuery := `
		SELECT id, name, business_id, created_at, updated_at
		FROM category
	`
	countQuery := `SELECT COUNT(*) FROM category`

	i := 1
	for key, value := range filter {
		if key == "id" || key == "business_id" || key == "name" {
			if where == "" {
				where = fmt.Sprintf(" WHERE %s = $%d", key, i)
			} else {
				where += fmt.Sprintf(" AND %s = $%d", key, i)
			}
			args = append(args, value)
			i++
		}
	}

	limitOffset := ""
	if limit > 0 {
		limitOffset = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}

	finalQuery := baseQuery + where + limitOffset
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

	// Count
	countQueryFinal := countQuery + where
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
		SET name = $1, business_id = $2, updated_at = $3
		WHERE id = $4
	`
	res, err := r.db.Exec(ctx, query, category.Name, category.BusinessID, time.Now(), category.ID)
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
	query := `DELETE FROM category WHERE id = $1`
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.db.Error(err)
	}
	if res.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no rows affected"))
	}
	return nil
}

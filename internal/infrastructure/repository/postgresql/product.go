package postgresql

import (
	"context"
	"fmt"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
	"time"
)

const productTableName = "product"

type productRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewProductRepo(db *postgres.Postgres) *productRepo {
	return &productRepo{
		tableName: productTableName,
		db:        db,
	}
}

func (p *productRepo) Create(ctx context.Context, product *entity.CreateProductRequest) error {
	query := `
		INSERT INTO product (
			business_id, name, category_id, short_info, description,
			cost, count, discount_cost, discount, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	_, err := p.db.Exec(ctx, query,
		product.BusinessID,
		product.Name,
		product.CategoryID,
		product.ShortInfo,
		product.Description,
		product.Cost,
		product.Count,
		product.DiscountCost,
		product.Discount,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p *productRepo) Get(ctx context.Context, id string) (*entity.Product, error) {
	query := `
		SELECT id, business_id, name, category_id, short_info, description,
		       cost, count, discount_cost, discount, created_at, updated_at
		FROM product
		WHERE id = $1
	`

	var product entity.Product
	err := p.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.BusinessID,
		&product.Name,
		&product.CategoryID,
		&product.ShortInfo,
		&product.Description,
		&product.Cost,
		&product.Count,
		&product.DiscountCost,
		&product.Discount,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return &product, nil
}

func (p *productRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.GetAllProductsResponse, error) {
	baseQuery := `
		SELECT id, business_id, name, category_id, short_info, description,
		       cost, count, discount_cost, discount, created_at, updated_at
		FROM product
		WHERE 1=1
	`
	args := []any{}
	argID := 1

	for key, value := range filter {
		if key == "category_id" || key == "owner_id" {
			baseQuery += fmt.Sprintf(" AND %s = $%d", key, argID)
			args = append(args, value)
			argID++
		}
		if key == "created_at" {
			baseQuery += fmt.Sprintf(" AND created_at::date = $%d", argID)
			args = append(args, value)
			argID++
		}
	}

	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
		args = append(args, limit, offset)
	}

	rows, err := p.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var products entity.GetAllProductsResponse
	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.ID,
			&product.BusinessID,
			&product.Name,
			&product.CategoryID,
			&product.ShortInfo,
			&product.Description,
			&product.Cost,
			&product.Count,
			&product.DiscountCost,
			&product.Discount,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		products.Items = append(products.Items, product)
	}

	// Count query
	countQuery := `SELECT COUNT(*) FROM product WHERE 1=1`
	countArgs := []any{}
	argID = 1

	for key, value := range filter {
		if key == "category_id" || key == "owner_id" {
			countQuery += fmt.Sprintf(" AND %s = $%d", key, argID)
			countArgs = append(countArgs, value)
			argID++
		}
		if key == "created_at" {
			countQuery += fmt.Sprintf(" AND created_at::date = $%d", argID)
			countArgs = append(countArgs, value)
			argID++
		}
	}

	var totalCount uint64
	err = p.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, p.db.Error(err)
	}

	products.Total = totalCount
	return &products, nil
}

func (p *productRepo) Update(ctx context.Context, product *entity.UpdateProductRequest) error {
	query := `
		UPDATE product
		SET business_id=$1, name=$2, category_id=$3, short_info=$4, description=$5,
		    cost=$6, count=$7, discount_cost=$8, discount=$9, updated_at=$10
		WHERE id=$11
	`

	res, err := p.db.Exec(ctx, query,
		product.BusinessID,
		product.Name,
		product.CategoryID,
		product.ShortInfo,
		product.Description,
		product.Cost,
		product.Count,
		product.DiscountCost,
		product.Discount,
		time.Now(),
		product.ID,
	)

	if err != nil {
		return p.db.Error(err)
	}

	if res.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *productRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM product WHERE id = $1`

	res, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return p.db.Error(err)
	}

	if res.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

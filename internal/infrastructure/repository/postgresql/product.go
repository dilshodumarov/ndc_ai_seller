package postgresql

import (
	"context"
	"fmt"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
	"time"

	"github.com/google/uuid"
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

func (p *productRepo) Create(ctx context.Context, product *entity.CreateProductRequest) (string,error) {
	id:=uuid.New().String()
	query := `
		INSERT INTO product (
			guid,business_id, name, category_id, short_info, description,
			cost, count, discount_cost, discount, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`
	
	if product.Discount > 0 {
		product.DiscountCost = product.Cost - (product.Cost * product.Discount / 100)
	} else {
		product.DiscountCost = product.Cost
	}
	
	_, err := p.db.Exec(ctx, query,
		id,
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
		fmt.Println(err)
		return "",p.db.Error(err)
	}

	return id,nil
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
	setParts := []string{}
	args := []interface{}{}
	argID := 1

	if product.BusinessID != "" {
		setParts = append(setParts, fmt.Sprintf("business_id=$%d", argID))
		args = append(args, product.BusinessID)
		argID++
	}
	if product.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name=$%d", argID))
		args = append(args, product.Name)
		argID++
	}
	if product.CategoryID != "" {
		setParts = append(setParts, fmt.Sprintf("category_id=$%d", argID))
		args = append(args, product.CategoryID)
		argID++
	}
	if product.ShortInfo != "" {
		setParts = append(setParts, fmt.Sprintf("short_info=$%d", argID))
		args = append(args, product.ShortInfo)
		argID++
	}
	if product.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description=$%d", argID))
		args = append(args, product.Description)
		argID++
	}
	if product.Cost != 0 {
		setParts = append(setParts, fmt.Sprintf("cost=$%d", argID))
		args = append(args, product.Cost)
		argID++
	}
	if product.Count != 0 {
		setParts = append(setParts, fmt.Sprintf("count=$%d", argID))
		args = append(args, product.Count)
		argID++
	}
	if product.Discount != 0 {
		setParts = append(setParts, fmt.Sprintf("discount=$%d", argID))
		args = append(args, product.Discount)
		argID++

		// discount_cost hisoblash
		discountCost := product.Cost - (product.Cost * product.Discount / 100)
		setParts = append(setParts, fmt.Sprintf("discount_cost=$%d", argID))
		args = append(args, discountCost)
		argID++
	}

	setParts = append(setParts, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now())
	argID++

	args = append(args, product.ID) // WHERE id=$n

	query := fmt.Sprintf(`UPDATE product SET %s WHERE guid=$%d`, 
	                     joinStrings(setParts, ", "), argID)

	res, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return p.db.Error(err)
	}
	if res.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}
	return nil
}

func joinStrings(parts []string, sep string) string {
	return fmt.Sprintf("%s", join(parts, sep))
}

func join(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for _, str := range s[1:] {
		result += sep + str
	}
	return result
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

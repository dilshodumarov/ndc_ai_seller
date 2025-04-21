package postgresql

import (
	"context"
	"database/sql"
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

func (p *productRepo) Create(ctx context.Context, product *entity.CreateProductRequest) (string, error) {
	id := uuid.New().String()
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
		return "", p.db.Error(err)
	}

	return id, nil
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

func (p *productRepo) List(ctx context.Context, filter entity.ProductFilter) (*entity.GetAllProductsResponse, error) {
	var (
		args      []any
		argID     = 1
		where     = "WHERE true"
		limitStmt string
	)

	// Ixtiyoriy OwnerID bo‘lsa, filter qilamiz
	if filter.OwnerID != "" {
		where += fmt.Sprintf(" AND business_id = $%d", argID)
		args = append(args, filter.OwnerID)
		argID++
	}

	// Ixtiyoriy CategoryID bo‘lsa, filter qilamiz
	if filter.CategoryID != "" {
		where += fmt.Sprintf(" AND category_id = $%d", argID)
		args = append(args, filter.CategoryID)
		argID++
	}

	// Qidiruv so‘zi bo‘lsa, name/description/short_info bo‘yicha qidiramiz
	if filter.Search != "" {
		where += fmt.Sprintf(` AND (name ILIKE $%d OR description ILIKE $%d OR short_info ILIKE $%d)`, argID, argID+1, argID+2)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
		argID += 3
	}

	// Pagination
	offset := (filter.Page - 1) * filter.Limit
	if filter.Limit > 0 {
		limitStmt = fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
		args = append(args, filter.Limit, offset)
		argID += 2
	}

	// Mahsulotlar ro‘yxatini olish uchun query
	query := fmt.Sprintf(`
		SELECT guid, business_id, name, category_id, short_info, description,
		       cost, count, discount_cost, discount, created_at, updated_at
		FROM product
		%s
		%s
	`, where, limitStmt)

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var products entity.GetAllProductsResponse
	for rows.Next() {
		var (
			product        entity.Product
			discountCostDB sql.NullInt64
			discountDB     sql.NullInt64
		)

		if err := rows.Scan(
			&product.ID,
			&product.BusinessID,
			&product.Name,
			&product.CategoryID,
			&product.ShortInfo,
			&product.Description,
			&product.Cost,
			&product.Count,
			&discountCostDB,
			&discountDB,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		if discountCostDB.Valid {
			product.DiscountCost = int(discountCostDB.Int64)
		}
		if discountDB.Valid {
			product.Discount = int(discountDB.Int64)
		}

		products.Items = append(products.Items, product)
	}

	// Umumiy sonini hisoblash (limit va offsetni olib tashlaymiz)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM product %s`, where)
	countArgs := args
	if filter.Limit > 0 {
		countArgs = args[:len(args)-2]
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

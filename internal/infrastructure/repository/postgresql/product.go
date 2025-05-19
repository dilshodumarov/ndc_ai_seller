package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
			guid,business_id, status,name, category_id, short_info, description,
			cost, count, discount_cost, discount, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`

	if product.Discount > 0 {
		product.DiscountCost = product.Cost - (product.Cost * product.Discount / 100)
	} else {
		product.DiscountCost = product.Cost
	}

	_, err := p.db.Exec(ctx, query,
		id,
		product.BusinessID,
		product.Status,
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
	var (
		picturesting string
	)
	query := `
	SELECT 
		p.guid, 
		p.business_id, 
		p.status,
		p.product_id,
		p.name, 
		p.category_id, 
		p.short_info, 
		p.description,
		p.cost, 
		p.count, 
		p.discount_cost, 
		p.discount, 
		p.created_at, 
		p.updated_at,
		c.name AS category_name, 
		COALESCE(STRING_AGG(pp.image_url, ','), '') AS image_urls
	FROM product p
	LEFT JOIN product_pictures pp ON p.guid = pp.product_id
	LEFT JOIN category c ON p.category_id = c.guid  -- category jadvali bilan bog'lanmoqda
	WHERE p.guid = $1 AND p.deleted_at is null
	GROUP BY p.guid, c.name
`

	var product entity.Product
	err := p.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.BusinessID,
		&product.Status,
		&product.ProductId,
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
		&product.CategoryName,
		&picturesting,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	// image_urlsni vergul bo'yicha ajratamiz
	if len(picturesting) > 0 {
		product.Image_urls = strings.Split(picturesting, ",")
	}

	return &product, nil
}

func (p *productRepo) List(ctx context.Context, filter entity.ProductFilter) (*entity.GetAllProductsResponse, error) {
	var (
		args      []any
		argID     = 1
		where     = "WHERE p.deleted_at is null"
		limitStmt string
	)
	fmt.Println(filter.OwnerID)

	// Ixtiyoriy OwnerID bo‘lsa, filter qilamiz
	if filter.OwnerID != "" {
		where += fmt.Sprintf(" AND p.business_id = $%d", argID)
		args = append(args, filter.OwnerID)
		argID++
	}

	if filter.ProductCount > 0 {
		where += fmt.Sprintf(" AND p.count = $%d", argID)
		args = append(args, filter.ProductCount)
		argID++
	}

	if filter.Status != "" {
		where += fmt.Sprintf(" AND p.status = $%d", argID)
		args = append(args, filter.Status)
		argID++
	}

	if filter.ProductId > 0 {
		where += fmt.Sprintf(" AND p.product_id = $%d", argID)
		args = append(args, filter.ProductId)
		argID++
	}

	// Ixtiyoriy CategoryID bo‘lsa, filter qilamiz
	if filter.CategoryID != "" {
		where += fmt.Sprintf(" AND p.category_id = $%d", argID)
		args = append(args, filter.CategoryID)
		argID++
	}

	// Qidiruv so‘zi bo‘lsa, name/description/short_info bo‘yicha qidiramiz
	if filter.Search != "" {
		where += fmt.Sprintf(` AND (p.name ILIKE $%d OR p.description ILIKE $%d OR p.short_info ILIKE $%d)`, argID, argID+1, argID+2)
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

	query := fmt.Sprintf(`
	SELECT p.guid, p.business_id, p.status, p.product_id, p.name, p.category_id, p.short_info, p.description,
	       p.cost, p.count, p.discount_cost, p.discount, p.created_at, p.updated_at, c.name,
	       COALESCE(STRING_AGG(pp.image_url, ','), '') AS image_urls
	FROM product p
	LEFT JOIN product_pictures pp ON p.guid = pp.product_id
	LEFT JOIN category c ON p.category_id = c.guid
	%s
	GROUP BY p.guid, p.business_id, p.status, p.product_id, p.name, p.category_id, p.short_info, p.description,
	         p.cost, p.count, p.discount_cost, p.discount, p.created_at, p.updated_at, c.name
	ORDER BY GREATEST(p.created_at, p.updated_at) DESC
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
			imageUrlsStr   string
		)

		if err := rows.Scan(
			&product.ID,
			&product.BusinessID,
			&product.Status,
			&product.ProductId,
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
			&product.CategoryName,
			&imageUrlsStr,
		); err != nil {
			return nil, p.db.Error(err)
		}

		// Discount values
		if discountCostDB.Valid {
			product.DiscountCost = int(discountCostDB.Int64)
		}
		if discountDB.Valid {
			product.Discount = int(discountDB.Int64)
		}

		// Split image URLs if they exist
		if imageUrlsStr != "" {
			product.Image_urls = strings.Split(imageUrlsStr, ",")
		}

		products.Items = append(products.Items, product)
	}

	// Count filtered results
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM product p %s`, where)
	var countArgs []any
	if filter.OwnerID != "" {
		countArgs = append(countArgs, filter.OwnerID)
	}
	if filter.ProductCount > 0 {
		countArgs = append(countArgs, filter.ProductCount)
	}
	if filter.Status != "" {
		countArgs = append(countArgs, filter.Status)
	}
	if filter.ProductId > 0 {
		countArgs = append(countArgs, filter.ProductId)
	}
	if filter.CategoryID != "" {
		countArgs = append(countArgs, filter.CategoryID)
	}
	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		countArgs = append(countArgs, searchTerm, searchTerm, searchTerm)
	}

	var Count uint64
	err = p.db.QueryRow(ctx, countQuery, countArgs...).Scan(&Count)
	if err != nil {
		return nil, p.db.Error(err)
	}

	// TotalCount: all non-deleted products (optionally filtered by OwnerID)
	var (
		TotalCount      uint64
		totalCountQuery string
	)

	if filter.OwnerID != "" {
		totalCountQuery = `SELECT COUNT(*) FROM product WHERE deleted_at IS NULL AND business_id = $1`
		err = p.db.QueryRow(ctx, totalCountQuery, filter.OwnerID).Scan(&TotalCount)
	} else {
		totalCountQuery = `SELECT COUNT(*) FROM product WHERE deleted_at IS NULL`
		err = p.db.QueryRow(ctx, totalCountQuery).Scan(&TotalCount)
	}
	if err != nil {
		return nil, p.db.Error(err)
	}

	products.Count = Count
	products.TotalCount = TotalCount

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
	if product.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status=$%d", argID))
		args = append(args, product.Status)
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

	args = append(args, product.ID)

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
	query := `UPDATE  product set deleted_at =$2 WHERE guid = $1`

	res, err := p.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return p.db.Error(err)
	}

	if res.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *productRepo) AddPicture(ctx context.Context, image *entity.CreateProductImage) (string, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO product_pictures (
			guid,product_id, image_url
		) VALUES ($1,$2,$3)
	`

	_, err := p.db.Exec(ctx, query,
		id,
		image.ProductId,
		image.ImageUrl,
	)
	if err != nil {
		fmt.Println(err)
		return "", p.db.Error(err)
	}

	return id, nil
}

func (p *productRepo) DeletePicture(ctx context.Context, id string) error {
	query := `
		delete from product_pictures
		where product_id =$1
	`

	_, err := p.db.Exec(ctx, query,
		id,
	)
	if err != nil {
		fmt.Println(err)
		return p.db.Error(err)
	}

	return nil
}

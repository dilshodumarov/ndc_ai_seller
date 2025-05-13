package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
	"time"
)

const (
	orderTableName = `"order"`
)

type OrderRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewOrderRepo(db *postgres.Postgres) *OrderRepo {
	return &OrderRepo{
		tableName: orderTableName,
		db:        db,
	}
}

func (r *OrderRepo) Create(ctx context.Context, o *entity.CreateOrderRequest) (id string, err error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (client_id,  business_id, location_url, status, total_price, payment_method, status_changed_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING guid
	`, r.tableName)

	now := time.Now().UTC()
	err = r.db.QueryRow(
		ctx, query,
		o.ClientID,
		o.BusinessID,
		o.LocationURL,
		o.Status,
		o.TotalPrice,
		o.PaymentMethod,
		now,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("OrderRepo - Create - insert order: %w", err)
	}

	// Insert Order Products
	for _, p := range o.Products {
		productQuery := `
			INSERT INTO order_products (order_id, product_id, count, price, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
		`
		_, err := r.db.Exec(ctx, productQuery, id, p.ProductID, p.Count, p.Price)
		if err != nil {
			return "", fmt.Errorf("OrderRepo - Create - insert order product: %w", err)
		}
	}

	return id, nil
}

func (r *OrderRepo) Get(ctx context.Context, id string) (*entity.Order, error) {
	query := `
	SELECT 
		o.guid, o.order_id,o.image_url,o.business_id,o.platform,o.location_url, o.status, o.total_price, 
		o.payment_method, o.status_changed_time, o.created_at, o.updated_at,

		c.guid, c.first_name, c.phone, 

		p.guid, p.name, p.image_url, p.cost,
		op.count, op.price
	FROM "order" o
	LEFT JOIN order_products op ON o.guid = op.order_id
	LEFT JOIN product p ON p.guid = op.product_id
	LEFT JOIN client c ON o.client_id = c.guid
	WHERE o.guid = $1
`


	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - Get - query order with products: %w", err)
	}
	defer rows.Close()

	var order *entity.Order
	products := make([]entity.OrderProduct, 0)

	for rows.Next() {
		var (
			o                     entity.Order
			product               entity.OrderProduct
			nullStatusChangedTime sql.NullTime
		
			clientGUID  sql.NullString
			clientName  sql.NullString
			clientPhone sql.NullString
			imageurl           sql.NullString
		)
		

		err = rows.Scan(
			&o.ID,  &o.OrderId,&imageurl,&o.BusinessID, &o.Platform,&o.LocationURL, &o.Status, &o.TotalPrice,
			&o.PaymentMethod, &nullStatusChangedTime, &o.CreatedAt, &o.UpdatedAt,
		
			&clientGUID, &clientName, &clientPhone, // client
		
			&product.ProductID, &product.Name, &product.ImageURL, &product.Cost,
			&product.Count, &product.ProductTotalPrice,
		)
		
		products = append(products, product)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - Get - scan row: %w", err)
		}

		if order == nil {
			order = &o
			if nullStatusChangedTime.Valid {
				order.StatusChangedTime = &nullStatusChangedTime.Time
			}
		}
		if clientGUID.Valid{
			order.Client.GUID=clientGUID.String
		}
		if imageurl.Valid{
			order.ImageUrl=imageurl.String
		}

		if clientName.Valid{
			order.Client.Name=clientName.String
		}
		if clientPhone.Valid{
			order.Client.Phone=clientPhone.String
		}

	}

	if order == nil {
		return nil, fmt.Errorf("OrderRepo - Get - no order found with id: %s", id)
	}

	order.Products = products

	return order, nil
}

func (r *OrderRepo) List(ctx context.Context, filter *entity.OrderFilter, limit, offset uint64) (*entity.GetAllOrdersResponse, error) {
	query := ""
	var where []string
	var args []interface{}
	argPos := 1
	
	if filter.ID != "" {
		where = append(where, fmt.Sprintf("o.guid = $%d", argPos))
		args = append(args, filter.ID)
		argPos++
	}
	
	if filter.Platform != "" {
		where = append(where, fmt.Sprintf("o.platform = $%d", argPos))
		args = append(args, filter.Platform)
		argPos++
	}
	
	if filter.ClientID != "" {
		where = append(where, fmt.Sprintf("o.client_id = $%d", argPos))
		args = append(args, filter.ClientID)
		argPos++
	}
	
	if filter.BusinessID != "" {
		where = append(where, fmt.Sprintf("o.business_id = $%d", argPos))
		args = append(args, filter.BusinessID)
		argPos++
	}
	
	if filter.Status != "" {
		where = append(where, fmt.Sprintf("o.status = $%d", argPos))
		args = append(args, filter.Status)
		argPos++
	}
	
	if filter.PaymentMethod != "" {
		where = append(where, fmt.Sprintf("o.payment_method = $%d", argPos))
		args = append(args, filter.PaymentMethod)
		argPos++
	}
	
	if filter.Search != "" {
		
		where = append(where, fmt.Sprintf("(p.name ILIKE $%d OR c.first_name ILIKE $%d)", argPos, argPos))
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}
	
	// WHERE blokini qo‘shish
	if len(where) > 0 {
		query = "WHERE " + strings.Join(where, " AND ")
	}
	
	// Yakuniy order ID'larini olish so‘rovi
	orderIDQuery := fmt.Sprintf(`
		SELECT o.guid
		FROM %s o
		LEFT JOIN order_products op ON o.guid = op.order_id
		LEFT JOIN product p ON p.guid = op.product_id
		LEFT JOIN client c ON o.client_id = c.guid
		%s
		ORDER BY o.created_at DESC
	`, r.tableName, query)
	
	if limit > 0 {
		orderIDQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, limit, offset)
	}

	orderIDRows, err := r.db.Query(ctx, orderIDQuery, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer orderIDRows.Close()

	var orderIDs []string
	for orderIDRows.Next() {
		var id string
		if err := orderIDRows.Scan(&id); err != nil {
			return nil, r.db.Error(err)
		}
		orderIDs = append(orderIDs, id)
	}
	if len(orderIDs) == 0 {
		return &entity.GetAllOrdersResponse{
			Items: []entity.Order{},
			Total: 0,
		}, nil
	}

	inClause := ""
	inArgs := make([]interface{}, len(orderIDs))
	for i, id := range orderIDs {
		inClause += fmt.Sprintf("$%d,", i+1)
		inArgs[i] = id
	}
	inClause = inClause[:len(inClause)-1]

	fullQuery := fmt.Sprintf(`
	SELECT 
		o.guid, o.order_id,o.image_url,
		c.guid AS client_guid, c.first_name AS client_name, c.phone AS client_phone, 
		o.business_id, o.platform,o.location_url, o.status, o.total_price, 
		o.payment_method, o.status_changed_time, o.created_at, o.updated_at,
		os.custom_name AS status_name,
		p.guid, p.name, p.image_url, p.cost,
		op.count, op.price
	FROM %s o
	LEFT JOIN order_status os ON o.order_status_id = os.guid
	INNER JOIN order_products op ON o.guid = op.order_id
	INNER JOIN product p ON p.guid = op.product_id
	INNER JOIN client c ON o.client_id = c.guid
	WHERE o.guid IN (%s)
	ORDER BY o.created_at DESC
`, r.tableName, inClause)

	rows, err := r.db.Query(ctx, fullQuery, inArgs...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()

	orderMap := make(map[string]*entity.Order)

	for rows.Next() {
		var (
			order                 entity.Order
			product               entity.OrderProduct
			nullStatusChangedTime sql.NullTime
			nullStatusName        sql.NullString
			clientGUID            sql.NullString
			clientName            sql.NullString
			clientPhone           sql.NullString
			imageurl              sql.NullString
		)

		if err := rows.Scan(
			&order.ID,
			&order.OrderId,
			&imageurl,
			&clientGUID, &clientName, &clientPhone,
			&order.BusinessID,
			&order.Platform,
			&order.LocationURL,
			&order.Status,
			&order.TotalPrice,
			&order.PaymentMethod,
			&nullStatusChangedTime,
			&order.CreatedAt,
			&order.UpdatedAt,
			&nullStatusName,
			&product.ProductID,
			&product.Name,
			&product.ImageURL,
			&product.Cost,
			&product.Count,
			&product.ProductTotalPrice,
		); err != nil {
			return nil, r.db.Error(err)
		}

		if nullStatusName.Valid {
			order.AdminStatus = nullStatusName.String
		}

		if clientGUID.Valid{
			order.Client.GUID=clientGUID.String
		}

		if clientName.Valid{
			order.Client.Name=clientName.String
		}
		if clientPhone.Valid{
			order.Client.Phone=clientPhone.String
		}
		if imageurl.Valid{
			order.ImageUrl=imageurl.String
		}
		if existingOrder, ok := orderMap[order.ID]; ok {
			existingOrder.Products = append(existingOrder.Products, product)
		} else {
			if nullStatusChangedTime.Valid {
				order.StatusChangedTime = &nullStatusChangedTime.Time
			}
			order.Products = append(order.Products, product)
			orderMap[order.ID] = &order
		}
	}

	var orders []entity.Order
	for _, id := range orderIDs {
		if order, ok := orderMap[id]; ok {
			orders = append(orders, *order)
		}
	}

	// 4. Count query
	countWhere, countArgs := buildWhereClause(filter, "")
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, r.tableName, countWhere)

	var totalCount uint64
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount); err != nil {
		return nil, r.db.Error(err)
	}

	return &entity.GetAllOrdersResponse{
		Items: orders,
		Total: totalCount,
	}, nil
}

// Update - Update Order
func (r *OrderRepo) Update(ctx context.Context, o *entity.OrderUpdate) error {
	var updates []string
	var args []interface{}
	argPos := 1

	if o.Status != "" {
		updates = append(updates, fmt.Sprintf("status = $%d", argPos))
		args = append(args, o.Status)
		argPos++
	}
	if o.LocationURL != "" {
		updates = append(updates, fmt.Sprintf("location_url = $%d", argPos))
		args = append(args, o.LocationURL)
		argPos++
	}
	if o.PaymentMethod != "" {
		updates = append(updates, fmt.Sprintf("payment_method = $%d", argPos))
		args = append(args, o.PaymentMethod)
		argPos++
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now().UTC())
	argPos++

	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE %s SET %s WHERE guid = $%d
	`, r.tableName, strings.Join(updates, ", "), argPos)

	args = append(args, o.ID)
	cmd, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return r.db.Error(err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("OrderRepo - Update - no rows affected")
	}

	return nil
}

// Delete - Delete Order
func (r *OrderRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE guid = $1`, r.tableName)
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.db.Error(err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("OrderRepo - Delete - no rows deleted")
	}
	return nil
}

// Helper function to build where clause
func buildWhereClause(filter *entity.OrderFilter, alias string) (string, []interface{}) {
	var where []string
	var args []interface{}
	argPos := 1

	col := func(name string) string {
		if alias != "" {
			return fmt.Sprintf("%s.%s", alias, name)
		}
		return name
	}

	if filter.ID != "" {
		where = append(where, fmt.Sprintf(`%s = $%d`, col("guid"), argPos))
		args = append(args, filter.ID)
		argPos++
	}
	if filter.Platform != "" {
		where = append(where, fmt.Sprintf(`%s = $%d`, col("platform"), argPos))
		args = append(args, filter.Platform)
		argPos++
	}

	if filter.ClientID != "" {
		where = append(where, fmt.Sprintf(`%s = $%d`, col("client_id"), argPos))
		args = append(args, filter.ClientID)
		argPos++
	}

	if filter.BusinessID != "" {
		where = append(where, fmt.Sprintf(`%s = $%d`, col("business_id"), argPos))
		args = append(args, filter.BusinessID)
		argPos++
	}

	if filter.Status != "" {
		where = append(where, fmt.Sprintf(`%s = $%d`, col("status"), argPos))
		args = append(args, filter.Status)
		argPos++
	}

	if filter.PaymentMethod != "" {
		where = append(where, fmt.Sprintf(`%s = $%d`, col("payment_method"), argPos))
		args = append(args, filter.PaymentMethod)
		argPos++
	}

	if len(where) > 0 {
		return "WHERE " + strings.Join(where, " AND "), args
	}
	return "", args
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

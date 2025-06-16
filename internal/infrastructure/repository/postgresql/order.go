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

func (r *OrderRepo) Create(ctx context.Context,o *entity.CreateOrderRequest) (id string, err error) {
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
		o.guid, o.order_id,o.status_number,o.image_url,o.business_id,o.platform,o.location_url, o.status, o.total_price, 
		o.payment_method, o.status_changed_time, o.created_at, o.updated_at,o.location,o.user_note,

		c.guid, c.first_name, c.phone, c.user_name,

		p.guid, p.name, p.image_url, p.cost,
		op.count, op.price
	FROM "order" o
	LEFT JOIN order_products op ON o.guid = op.order_id
	LEFT JOIN product p ON p.guid = op.product_id
	LEFT JOIN client c ON o.client_id = c.guid
	WHERE o.guid = $1 and o.deleted_at is null
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

			clientGUID    sql.NullString
			clientName    sql.NullString
			clientPhone   sql.NullString
			imageurl      sql.NullString
			primageurl    sql.NullString
			username      sql.NullString
			paymentMethod sql.NullString
			location      sql.NullString
			description   sql.NullString
		)

		err = rows.Scan(
			&o.ID, &o.OrderId, &o.StatusNumber, &imageurl, &o.BusinessID, &o.Platform, &o.LocationURL, &o.Status, &o.TotalPrice,
			&paymentMethod, &nullStatusChangedTime, &o.CreatedAt, &o.UpdatedAt, &location, &description,

			&clientGUID, &clientName, &clientPhone, &username, // client

			&product.ProductID, &product.Name, &primageurl, &product.Cost,
			&product.Count, &product.ProductTotalPrice,
		)

		if primageurl.Valid {
			product.ImageURL = primageurl.String
		}
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
		if clientGUID.Valid {
			order.Client.GUID = clientGUID.String
		}
		if location.Valid {
			order.Location = location.String
		}
		if description.Valid {
			order.Description = description.String
		}
		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		if username.Valid {
			order.Client.UserName = username.String
		}
		if imageurl.Valid {
			order.ImageUrl = imageurl.String
		}

		if clientName.Valid {
			order.Client.Name = clientName.String
		}
		if clientPhone.Valid {
			order.Client.Phone = clientPhone.String
		}

	}

	if order == nil {
		return nil, fmt.Errorf("OrderRepo - Get - no order found with id: %s", id)
	}

	order.Products = products

	return order, nil
}

func (r *OrderRepo) List(ctx context.Context, filter *entity.OrderFilter, limit, offset uint64) (*entity.GetAllOrdersResponse, error) {
	var where []string
	var args []interface{}
	argPos := 1
	// Filterlar
	if filter.ID != "" {
		where = append(where, fmt.Sprintf("o.guid = $%d", argPos))
		args = append(args, filter.ID)
		argPos++
	}
	if filter.Daye > 0 {
		if filter.Daye > 0 {
			where = append(where, fmt.Sprintf("o.created_at >= NOW() - INTERVAL '%d days'", filter.Daye))
		}

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

	// WHERE clause
	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	// Step 1: Faqat order_id'larni olish
	orderIDQuery := fmt.Sprintf(`
	SELECT DISTINCT ON (o.guid) o.guid
	FROM %s o
	LEFT JOIN order_products op ON o.guid = op.order_id
	LEFT JOIN product p ON p.guid = op.product_id
	LEFT JOIN client c ON o.client_id = c.guid
	%s 
	ORDER BY o.guid, o.created_at DESC
`, r.tableName, whereClause)

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

	// Step 2: IN clause bilan to‘liq ma'lumotlarni olish
	inClause := ""
	inArgs := make([]interface{}, len(orderIDs))
	for i, id := range orderIDs {
		inClause += fmt.Sprintf("$%d,", i+1)
		inArgs[i] = id
	}
	inClause = inClause[:len(inClause)-1]

	fullQuery := fmt.Sprintf(`
	SELECT 
		o.guid, o.order_id, o.status_number, o.image_url,
		c.guid, c.first_name, c.phone, c.user_name,
		o.business_id, o.platform, o.location_url, o.status, o.total_price,
		o.payment_method, o.status_changed_time, o.created_at, o.updated_at,o.location,o.user_note,
		os.custom_name,
		p.guid, p.name, p.image_url, p.cost,
		op.count, op.total_price
	FROM %s o
	LEFT JOIN order_status os ON o.order_status_id = os.guid
	INNER JOIN order_products op ON o.guid = op.order_id
	INNER JOIN product p ON p.guid = op.product_id
	INNER JOIN client c ON o.client_id = c.guid
	WHERE o.guid IN (%s) AND o.deleted_at IS NULL
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
			order                               entity.Order
			product                             entity.OrderProduct
			nullStatusChangedTime               sql.NullTime
			nullStatusName                      sql.NullString
			clientGUID, clientName, clientPhone sql.NullString
			imageURL                            sql.NullString
			primageURL                          sql.NullString
			username                            sql.NullString
			paymentMethod                       sql.NullString
			location                            sql.NullString
			description                         sql.NullString
		)

		if err := rows.Scan(
			&order.ID,
			&order.OrderId,
			&order.StatusNumber,
			&imageURL,
			&clientGUID, &clientName, &clientPhone, &username,
			&order.BusinessID, &order.Platform, &order.LocationURL, &order.Status,
			&order.TotalPrice, &paymentMethod,
			&nullStatusChangedTime, &order.CreatedAt, &order.UpdatedAt, &location, &description,
			&nullStatusName,
			&product.ProductID, &product.Name, &primageURL, &product.Cost,
			&product.Count, &product.ProductTotalPrice,
		); err != nil {
			return nil, r.db.Error(err)
		}

		if existingOrder, ok := orderMap[order.ID]; ok {
			existingOrder.Products = append(existingOrder.Products, product)
			continue
		}

		// yangi order
		if nullStatusChangedTime.Valid {
			order.StatusChangedTime = &nullStatusChangedTime.Time
		}
		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		if location.Valid {
			order.Location = location.String
		}
		if description.Valid {
			order.Description = description.String
		}

		if nullStatusName.Valid {
			order.AdminStatus = nullStatusName.String
		}
		if primageURL.Valid {
			product.ImageURL = primageURL.String
		}
		if imageURL.Valid {
			order.ImageUrl = imageURL.String
		}
		order.Client = entity.ClientInfo{}
		if clientGUID.Valid {
			order.Client.GUID = clientGUID.String
		}
		if username.Valid {
			order.Client.UserName = username.String
		}
		if clientName.Valid {
			order.Client.Name = clientName.String
		}
		if clientPhone.Valid {
			order.Client.Phone = clientPhone.String
		}

		order.Products = []entity.OrderProduct{product}
		orderMap[order.ID] = &order
	}

	var orders []entity.Order
	for _, id := range orderIDs {
		if order, ok := orderMap[id]; ok {
			orders = append(orders, *order)
		}
	}

	// Step 3: total count
	countWhere, countArgs := buildWhereClause(filter, "o")
	countQuery := fmt.Sprintf(`SELECT COUNT(DISTINCT o.guid) FROM %s o
	LEFT JOIN order_products op ON o.guid = op.order_id
	LEFT JOIN product p ON p.guid = op.product_id
	LEFT JOIN client c ON o.client_id = c.guid
	%s AND o.deleted_at is null`, r.tableName, countWhere)

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
	fmt.Println(1111, o)
	if o.Status != "" {
		updates = append(updates, fmt.Sprintf("status = $%d", argPos))
		args = append(args, o.Status)
		argPos++
		updates = append(updates, fmt.Sprintf("status_number = $%d", argPos))
		args = append(args, o.StatusNumber)
		argPos++
		updates = append(updates, fmt.Sprintf("order_status_id = $%d", argPos))
		args = append(args, o.StatusID)
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

func (r *OrderRepo) GetProductsByOrderID(ctx context.Context, orderID string) ([]entity.OrderProductBuOrderID, error) {
	query := `
	SELECT 
	p.guid, p.name, p.image_url, p.cost, p.status, p.discount_cost, p.discount,
		p.short_info, p.description, p.created_at, p.updated_at,
		op.count, op.price, op.total_price, op.created_at
	FROM "order" o
	JOIN "order" parent_order ON parent_order.guid = o.order_guid
	JOIN order_products op ON parent_order.guid = op.order_id
	JOIN product p ON op.product_id = p.guid
	WHERE o.guid = $1;
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetProductsByOrderID - query: %w", err)
	}
	defer rows.Close()

	var products []entity.OrderProductBuOrderID

	for rows.Next() {
		var (
			product entity.OrderProductBuOrderID

			imageURL, shortInfo, description                sql.NullString
			status                                          sql.NullBool
			discountCost, discount                          sql.NullInt64
			productCreatedAt, productUpdatedAt, opCreatedAt sql.NullTime
		)

		err := rows.Scan(
			&product.ProductID, &product.Name, &imageURL, &product.Cost, &status, &discountCost, &discount,
			&shortInfo, &description, &productCreatedAt, &productUpdatedAt,

			&product.Count, &product.Price, &product.ProductTotalPrice, &opCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("OrderRepo - GetProductsByOrderID - scan: %w", err)
		}

		if imageURL.Valid {
			product.ImageURL = imageURL.String
		}
		if shortInfo.Valid {
			product.ShortInfo = shortInfo.String
		}
		if description.Valid {
			product.Description = description.String
		}
		if status.Valid {
			product.Status = status.Bool
		}
		if discount.Valid {
			product.Discount = int(discount.Int64)
		}
		if discountCost.Valid {
			product.DiscountCost = int(discountCost.Int64)
		}
		if productCreatedAt.Valid {
			product.CreatedAt = productCreatedAt.Time
		}
		if productUpdatedAt.Valid {
			product.UpdatedAt = productUpdatedAt.Time
		}
		if opCreatedAt.Valid {
			product.OrderProductCreatedAt = opCreatedAt.Time
		}

		products = append(products, product)
	}

	return products, nil
}

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

// CreateOrder -.
func (r *OrderRepo) Create(ctx context.Context, o *entity.CreateOrderRequest) (id string, err error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (client_id, integration_id, status, status_changed_time, total_cost, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, r.tableName)

	now := time.Now().UTC()
	err = r.db.QueryRow(
		ctx, query,
		o.ClientID, o.IntegrationID, o.Status, now, o.TotalCost, now, now,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("OrderRepo - CreateOrder - r.db.QueryRow: %w", err)
	}

	return id, nil
}

// GetOrder -.
func (r *OrderRepo) Get(ctx context.Context, params map[string]string) (*entity.Order, error) {
	var o entity.Order

	whereClause := ""
	args := []interface{}{}
	argPos := 1
	for key, value := range params {
		if whereClause == "" {
			whereClause = fmt.Sprintf(`WHERE %s = $%d`, key, argPos)
		} else {
			whereClause += fmt.Sprintf(` AND %s = $%d`, key, argPos)
		}
		args = append(args, value)
		argPos++
	}

	query := fmt.Sprintf(`
		SELECT id, client_id, integration_id, status, status_changed_time,  created_at, updated_at
		FROM %s %s
		LIMIT 1
	`, r.tableName, whereClause)

	row := r.db.QueryRow(ctx, query, args...)
	err := row.Scan(
		&o.ID,
		&o.ClientID,
		&o.IntegrationID,
		&o.Status,
		&o.StatusChangedTime,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetOrder - r.db.QueryRow: %w", err)
	}

	return &o, nil
}

// ListOrders -.
func (r *OrderRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.GetAllOrdersResponse, error) {
	var orders []entity.Order
	var args []interface{}
	query := `SELECT id, client_id, integration_id, status, status_changed_time,  created_at, updated_at FROM "order"`
	whereClauses := ""
	argPos := 1

	for key, value := range filter {
		if key == "id" || key == "client_id" || key == "integration_id" || key == "status" {
			if whereClauses == "" {
				whereClauses = fmt.Sprintf(` WHERE %s = $%d`, key, argPos)
			} else {
				whereClauses += fmt.Sprintf(` AND %s = $%d`, key, argPos)
			}
			args = append(args, value)
			argPos++
		}
	}

	if limit > 0 {
		query += whereClauses + fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argPos, argPos+1)
		args = append(args, limit, offset)
	} else {
		query += whereClauses
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var o entity.Order
		var nullStatusChangedTime sql.NullTime

		err = rows.Scan(
			&o.ID,
			&o.ClientID,
			&o.IntegrationID,
			&o.Status,
			&nullStatusChangedTime,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, r.db.Error(err)
		}

		if nullStatusChangedTime.Valid {
			o.StatusChangedTime = &nullStatusChangedTime.Time
		} else {
			o.StatusChangedTime = nil
		}

		orders = append(orders, o)
	}

	countQuery := `SELECT COUNT(*) FROM "order"` + whereClauses
	var totalCount uint64
	err = r.db.QueryRow(ctx, countQuery, args[:len(args)-(2*boolToInt(limit > 0))]...).Scan(&totalCount)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &entity.GetAllOrdersResponse{
		Items: orders,
		Total: totalCount,
	}, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// UpdateOrder -.
func (r *OrderRepo) Update(ctx context.Context, o *entity.Order) error {
	// Dinamik filterlar uchun
	var updateFields []string
	var args []interface{}
	argPos := 1

	// client_id, integration_id va boshqa atributlar uchun update qilish
	if o.ClientID != "" {
		updateFields = append(updateFields, fmt.Sprintf("client_id = $%d", argPos))
		args = append(args, o.ClientID)
		argPos++
	}

	if o.IntegrationID != "" {
		updateFields = append(updateFields, fmt.Sprintf("integration_id = $%d", argPos))
		args = append(args, o.IntegrationID)
		argPos++
	}

	if o.Status != "" {
		updateFields = append(updateFields, fmt.Sprintf("status = $%d", argPos))
		args = append(args, o.Status)
		argPos++
	}

	// updated_at maydoni har doim yangilanishi kerak
	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now().UTC())
	argPos++

	// IDni oxirgi argument sifatida qo'shish
	updateFieldsStr := strings.Join(updateFields, ", ")
	args = append(args, o.ID)

	// SQL so'rovini yaratish
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $%d`, r.tableName, updateFieldsStr, argPos)
	commandTag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return r.db.Error(err)
	}

	// Hech qanday satr yangilanmasa
	if commandTag.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no sql rows affected"))
	}

	return nil
}

// DeleteOrder -.
func (r *OrderRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, r.tableName)

	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

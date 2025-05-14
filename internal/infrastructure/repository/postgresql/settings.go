package postgresql

import (
	"context"
	"fmt"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
)

type settingsRepo struct {
	db *postgres.Postgres
}

func NewSettingsRepo(db *postgres.Postgres) *settingsRepo {
	return &settingsRepo{db: db}
}

// Create creates a new order status.
func (r *settingsRepo) Create(ctx context.Context, req *entity.CreateOrderStatusRequest) error {
	query := `
		INSERT INTO order_status (business_id, type_id, custom_name)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, req.BusinessID, req.TypeID, req.CustomName)
	if err != nil {
		return fmt.Errorf("Create order_status: %w", err)
	}
	return nil
}

// Get retrieves a specific order status with type name.
func (r *settingsRepo) Get(ctx context.Context, guid string) (*entity.OrderStatus, error) {
	query := `
		SELECT os.guid, os.business_id, os.type_id, os.custom_name, ost.name as type_name, os.created_at
		FROM order_status os
		JOIN order_status_type ost ON os.type_id = ost.guid
		WHERE os.guid = $1
	`
	row := r.db.QueryRow(ctx, query, guid)

	var status entity.OrderStatus
	err := row.Scan(
		&status.GUID, &status.BusinessID, &status.TypeID,
		&status.CustomName, &status.TypeName, &status.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Get order_status: %w", err)
	}
	return &status, nil
}

// Update modifies an existing order status.
func (r *settingsRepo) Update(ctx context.Context, req *entity.UpdateOrderStatusRequest) error {
	query := `
		UPDATE order_status
		SET custom_name = $1
		WHERE guid = $2
	`
	_, err := r.db.Exec(ctx, query, req.CustomName, req.GUID)
	if err != nil {
		return fmt.Errorf("Update order_status: %w", err)
	}
	return nil
}

// Delete removes an order status by ID.
func (r *settingsRepo) Delete(ctx context.Context, guid string) error {
	query := `DELETE FROM order_status WHERE guid = $1`
	_, err := r.db.Exec(ctx, query, guid)
	if err != nil {
		return fmt.Errorf("Delete order_status: %w", err)
	}
	return nil
}

// List returns all order statuses with their type names for a business.
func (r *settingsRepo) List(ctx context.Context, businessID string) ([]*entity.OrderStatus, error) {
	query := `
		SELECT 
			os.guid,
			os.business_id,
			os.type_id,
			os.custom_name,
			ost.name as type_name,
			ost.status_number,
			os.created_at,
			COALESCE(COUNT(o.guid), 0) AS order_count
		FROM order_status os
		JOIN order_status_type ost ON os.type_id = ost.guid
		LEFT JOIN "order" o ON 
			o.order_status_id = os.guid AND 
			o.status = ost.name AND 
			o.business_id = os.business_id
		WHERE os.business_id = $1
		GROUP BY os.guid, ost.name, ost.status_number
		ORDER BY os.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, businessID)
	if err != nil {
		return nil, fmt.Errorf("List order_status: %w", err)
	}
	defer rows.Close()

	var result []*entity.OrderStatus
	for rows.Next() {
		var status entity.OrderStatus
		err := rows.Scan(
			&status.GUID,
			&status.BusinessID,
			&status.TypeID,
			&status.CustomName,
			&status.TypeName,
			&status.StatusNumber,
			&status.CreatedAt,
			&status.OrderCount,
		)
		if err != nil {
			return nil, fmt.Errorf("List scan: %w", err)
		}
		result = append(result, &status)
	}
	return result, nil
}



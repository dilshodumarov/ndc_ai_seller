package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
	"time"
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
	// Step 1: Tekshir order_status mavjudmi
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM order_status WHERE business_id = $1
		)
	`
	var exists bool
	if err := r.db.QueryRow(ctx, checkQuery, businessID).Scan(&exists); err != nil {
		return nil, fmt.Errorf("check order_status exists: %w", err)
	}

	var query string
	if exists {
		// order_status mavjud bo‘lsa
		query = `
			SELECT 
				os.guid,
				os.business_id,
				os.type_id,
				os.custom_name,
				ost.name as type_name,
				ost.status_number,
				os.created_at,
				COUNT(o.guid) AS order_count
			FROM order_status_type ost
			LEFT JOIN order_status os 
				ON os.type_id = ost.guid AND os.business_id = $1
			LEFT JOIN "order" o 
				ON o.status = ost.name AND o.business_id = $1
			GROUP BY os.guid, os.business_id, os.type_id, os.custom_name, os.created_at,
			         ost.name, ost.status_number
			ORDER BY ost.status_number
		`
	} else {
		// order_status yo‘q bo‘lsa
		query = `
			SELECT 
				NULL AS guid,
				NULL AS business_id,
				NULL AS type_id,
				NULL AS custom_name,
				ost.name AS type_name,
				ost.status_number,
				NULL AS created_at,
				COUNT(o.guid) AS order_count
			FROM order_status_type ost
			LEFT JOIN "order" o 
				ON o.status = ost.name AND o.business_id = $1
			GROUP BY ost.name, ost.status_number
			ORDER BY ost.status_number
		`
	}

	rows, err := r.db.Query(ctx, query, businessID)
	if err != nil {
		return nil, fmt.Errorf("List order_status: %w", err)
	}
	defer rows.Close()

	var result []*entity.OrderStatus
	for rows.Next() {
		var (
			guid, businessID, typeID, customName sql.NullString
			createdAt                            sql.NullTime
			status                               entity.OrderStatus
		)

		err := rows.Scan(
			&guid,
			&businessID,
			&typeID,
			&customName,
			&status.TypeName,
			&status.StatusNumber,
			&createdAt,
			&status.OrderCount,
		)
		if err != nil {
			return nil, fmt.Errorf("List scan: %w", err)
		}

		// Null bo'lishini tekshiramiz va kerakli structga o‘tkazamiz
		if guid.Valid {
			status.GUID = guid.String
		}
		if businessID.Valid {
			status.BusinessID = businessID.String
		}
		if typeID.Valid {
			status.TypeID = typeID.String
		}
		if customName.Valid {
			status.CustomName = customName.String
		}
		if createdAt.Valid {
			status.CreatedAt = createdAt.Time
		}

		result = append(result, &status)
	}

	return result, nil
}

func (r *settingsRepo) GetStatusByName(ctx context.Context, name, bussnesid string) (*string, error) {
	query := `
		SELECT os.guid
		FROM order_status os
		JOIN order_status_type ost ON os.type_id = ost.guid
		WHERE ost.name = $1 and os.business_id=$2
	`
	row := r.db.QueryRow(ctx, query, name, bussnesid)

	var guid string
	err := row.Scan(&guid)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("Get order_status: %w", err)
	}

	return &guid, nil
}

func (r *settingsRepo) CreateSettings(ctx context.Context, req *entity.CreateSettingsRequest) error {
	query := `
		INSERT INTO settings (
			name, status, business_id, prompt_text, prompt_order,
			waiting_time, prompt_product, token_limit,
			intelligence_level, stop_until
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		req.Name, req.Status, req.BusinessID, req.PromptText,
		req.PromptOrder, req.WaitingTime, req.PromptProduct,
		req.TokenLimit, req.IntelligenceLevel, req.StopUntil,
	)
	if err != nil {
		return fmt.Errorf("Create settings: %w", err)
	}
	return nil
}

func (r *settingsRepo) GetSettings(ctx context.Context, guid string) (*entity.Settings, error) {
	query := `
		SELECT guid, name, status, business_id, prompt_text, prompt_order,
		       waiting_time, prompt_product, token_limit, intelligence_level,
		       stop_until, created_at, updated_at, deleted_at
		FROM settings
		WHERE guid = $1
	`

	row := r.db.QueryRow(ctx, query, guid)
	var s entity.Settings

	err := row.Scan(
		&s.GUID, &s.Name, &s.Status, &s.BusinessID, &s.PromptText,
		&s.PromptOrder, &s.WaitingTime, &s.PromptProduct,
		&s.TokenLimit, &s.IntelligenceLevel, &s.StopUntil,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("Get settings: %w", err)
	}
	return &s, nil
}

func (r *settingsRepo) UpdateSettings(ctx context.Context, req *entity.UpdateSettingsRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argID := 1

	if req.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name=$%d", argID))
		args = append(args, req.Name)
		argID++
	}
	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status=$%d", argID))
		args = append(args, req.Status)
		argID++
	}
	if req.PromptText != "" {
		setParts = append(setParts, fmt.Sprintf("prompt_text=$%d", argID))
		args = append(args, req.PromptText)
		argID++
	}
	if req.PromptOrder != "" {
		setParts = append(setParts, fmt.Sprintf("prompt_order=$%d", argID))
		args = append(args, req.PromptOrder)
		argID++
	}
	if req.WaitingTime != 0 {
		setParts = append(setParts, fmt.Sprintf("waiting_time=$%d", argID))
		args = append(args, req.WaitingTime)
		argID++
	}
	if req.PromptProduct != "" {
		setParts = append(setParts, fmt.Sprintf("prompt_product=$%d", argID))
		args = append(args, req.PromptProduct)
		argID++
	}
	if req.TokenLimit != 0 {
		setParts = append(setParts, fmt.Sprintf("token_limit=$%d", argID))
		args = append(args, req.TokenLimit)
		argID++
	}
	if req.IntelligenceLevel != 0 {
		setParts = append(setParts, fmt.Sprintf("intelligence_level=$%d", argID))
		args = append(args, req.IntelligenceLevel)
		argID++
	}
	if req.StopUntil > 0 {
		setParts = append(setParts, fmt.Sprintf("stop_until=$%d", argID))
		args = append(args, req.StopUntil)
		argID++
	}

	// always update updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, time.Now())
	argID++

	// WHERE guid = ...
	args = append(args, req.GUID)

	query := fmt.Sprintf(`
		UPDATE settings
		SET %s
		WHERE guid = $%d
	`, joinStrings(setParts, ", "), argID)

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Update settings: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (r *settingsRepo) DeleteSettings(ctx context.Context, guid string) error {
	query := `DELETE FROM settings WHERE guid = $1`
	_, err := r.db.Exec(ctx, query, guid)
	if err != nil {
		return fmt.Errorf("Delete settings: %w", err)
	}
	return nil
}

func (r *settingsRepo) ListSettingsByBusinessID(ctx context.Context, businessID string) ([]*entity.Settings, error) {
	query := `
		SELECT guid, name, status, business_id, prompt_text, prompt_order,
		       waiting_time, prompt_product, token_limit, intelligence_level,
		       stop_until, created_at, updated_at, deleted_at
		FROM settings
		WHERE business_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, businessID)
	if err != nil {
		return nil, fmt.Errorf("List settings: %w", err)
	}
	defer rows.Close()

	var result []*entity.Settings
	for rows.Next() {
		var s entity.Settings
		err := rows.Scan(
			&s.GUID, &s.Name, &s.Status, &s.BusinessID, &s.PromptText,
			&s.PromptOrder, &s.WaitingTime, &s.PromptProduct,
			&s.TokenLimit, &s.IntelligenceLevel, &s.StopUntil,
			&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("List scan: %w", err)
		}
		result = append(result, &s)
	}

	return result, nil
}

func (r *settingsRepo) GetSettingsByName(ctx context.Context, name, businessID string) (*entity.Settings, error) {
	query := `
		SELECT guid, name, status, business_id, prompt_text, prompt_order,
		       waiting_time, prompt_product, token_limit, intelligence_level,
		       stop_until, created_at, updated_at, deleted_at
		FROM settings
		WHERE name = $1 AND business_id = $2
	`

	row := r.db.QueryRow(ctx, query, name, businessID)
	var s entity.Settings

	err := row.Scan(
		&s.GUID, &s.Name, &s.Status, &s.BusinessID, &s.PromptText,
		&s.PromptOrder, &s.WaitingTime, &s.PromptProduct,
		&s.TokenLimit, &s.IntelligenceLevel, &s.StopUntil,
		&s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("Get settings by name: %w", err)
	}

	return &s, nil
}

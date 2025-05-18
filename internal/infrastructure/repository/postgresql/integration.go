package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
)

const integrationTable = "integration"

type integrationRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewIntegrationRepo(db *postgres.Postgres) *integrationRepo {
	return &integrationRepo{
		tableName: integrationTable,
		db:        db,
	}
}

func (r *integrationRepo) Create(ctx context.Context, req *entity.IntegrationCreate) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (owner_id, integration_token, integration_type, status) 
		VALUES ($1, $2, $3, $4)
	`, r.tableName)

	_, err := r.db.Exec(ctx, query,
		req.BusinessId,
		req.IntegrationToken,
		req.IntegrationType,
		"active",
	)
	if err != nil {
		return r.db.Error(err)
	}

	return nil
}

func (r *integrationRepo) Update(ctx context.Context, req *entity.IntegrationUpdate) (*entity.IntegrationUpdateResponse, error) {
	setClauses := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{}
	argPos := 1

	if req.Token != "" {
		setClauses = append(setClauses, fmt.Sprintf("integration_token = $%d", argPos))
		args = append(args, req.Token)
		argPos++
	}
	if req.PromptText != "" {
		setClauses = append(setClauses, fmt.Sprintf("prompt_text = $%d", argPos))
		args = append(args, req.PromptText)
		argPos++
	}

	if req.PromptOrder != "" {
		setClauses = append(setClauses, fmt.Sprintf("prompt_order = $%d", argPos))
		args = append(args, req.PromptOrder)
		argPos++
	}
	if req.PromptProduct != "" {
		setClauses = append(setClauses, fmt.Sprintf("prompt_product = $%d", argPos))
		args = append(args, req.PromptProduct)
		argPos++
	}
	if req.TokenLimit != 0 {
		setClauses = append(setClauses, fmt.Sprintf("token_limit = $%d", argPos))
		args = append(args, req.TokenLimit)
		argPos++
	}
	if req.IntelligenceLevel != 0 {
		setClauses = append(setClauses, fmt.Sprintf("intelligence_level = $%d", argPos))
		args = append(args, req.IntelligenceLevel)
		argPos++
	}

	if req.StopUntil > 0 {
		setClauses = append(setClauses, fmt.Sprintf("stop_until = $%d", argPos))
		args = append(args, req.StopUntil)
		argPos++
	}

	// Har doim kerak: guid va deleted_at IS NULL
	setClause := strings.Join(setClauses, ", ")
	args = append(args, req.ID)

	query := fmt.Sprintf(`
		UPDATE %s 
		SET %s
		WHERE guid = $%d AND deleted_at IS NULL
		RETURNING integration_type, owner_id
	`, r.tableName, setClause, argPos)

	var integration entity.IntegrationUpdateResponse
	err := r.db.QueryRow(ctx, query, args...).Scan(&integration.Itype, &integration.GUID)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &integration, nil
}

func (r *integrationRepo) UpdateStatus(ctx context.Context, req *entity.IntegrationUpdateStatus) (*entity.IntegrationUpdateStatusResponse, error) {

	if req.Status != "active" && req.Status != "stop" {
		return nil, fmt.Errorf("invalid status value: must be 'active' or 'stop'")
	}

	query := fmt.Sprintf(`
		UPDATE %s 
		SET 
			status = $1,
			updated_at = CURRENT_TIMESTAMP,
			started_at = CASE WHEN $2 = 'active' THEN CURRENT_TIMESTAMP ELSE started_at END,
			stoped_at  = CASE WHEN $2 = 'stop' THEN CURRENT_TIMESTAMP ELSE stoped_at END
		WHERE guid = $3 AND deleted_at IS NULL
		RETURNING integration_type, owner_id, integration_token
	`, r.tableName)

	var resp entity.IntegrationUpdateStatusResponse
	err := r.db.QueryRow(ctx, query, req.Status, req.Status, req.ID).Scan(
		&resp.IntegrationType,
		&resp.BusinessId,
		&resp.IntegrationToken,
	)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &resp, nil
}

func (r *integrationRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`
		UPDATE %s 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE guid = $1 AND deleted_at IS NULL
	`, r.tableName)

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.db.Error(err)
	}

	return nil
}

func (r *integrationRepo) GetByOwnerID(ctx context.Context, req *entity.IntegrationRequest) (*entity.IntegrationGetResponse, error) {
	query := fmt.Sprintf(`
		SELECT guid, integration_token, integration_type, status, started_at, stoped_at,
		       prompt_text, prompt_order,prompt_product,token_limit, intelligence_level
		FROM %s 
		WHERE owner_id = $1 AND deleted_at IS NULL
	`, r.tableName)

	var (
		startedAt     sql.NullTime
		stoppedAt     sql.NullTime
		prompt        sql.NullString
		promptOrder   sql.NullString
		promptProduct sql.NullString
		res           entity.IntegrationGetResponse
	)

	err := r.db.QueryRow(ctx, query, req.BusinessId).Scan(
		&res.Guid,
		&res.IntegrationToken,
		&res.IntegrationType,
		&res.Status,
		&startedAt,
		&stoppedAt,
		&prompt,
		&promptOrder,
		&promptProduct,
		&res.TokenLimit,
		&res.IntelligenceLevel,
	)
	if err != nil {
		return nil, r.db.Error(err)
	}

	if startedAt.Valid {
		res.StartedAt = startedAt.Time
	}
	if promptProduct.Valid {
		res.PromtProduct = promptProduct.String
	}
	if stoppedAt.Valid {
		res.StoppedAt = stoppedAt.Time
	}
	if prompt.Valid {
		res.PromptText = prompt.String
	}

	if promptOrder.Valid {
		res.PromtOrder = promptOrder.String
	}

	return &res, nil
}


func (r *integrationRepo) CheckIntegrationExistence(ctx context.Context, businessID string) (*entity.IntegrationExistenceResponse, error) {
	query := fmt.Sprintf(`
		SELECT integration_type
		FROM %s 
		WHERE owner_id = $1 AND deleted_at IS NULL
	`, r.tableName)

	rows, err := r.db.Query(ctx, query, businessID)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()

	var (
		integrationType string
		exists          = &entity.IntegrationExistenceResponse{}
	)

	for rows.Next() {
		if err := rows.Scan(&integrationType); err != nil {
			return nil, r.db.Error(err)
		}

		switch integrationType {
		case "telegram_account":
			exists.TelegramAccount = true
		case "bot":
			exists.TelegramBot = true
		case "instagram":
			exists.Instagram = true
		}
	}

	return exists, nil
}

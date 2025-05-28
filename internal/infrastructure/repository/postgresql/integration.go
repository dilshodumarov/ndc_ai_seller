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
		INSERT INTO %s (owner_id, integration_token, status) 
		VALUES ($1, $2, $3)
	`, r.tableName)

	_, err := r.db.Exec(ctx, query,
		req.BusinessId,
		req.IntegrationToken,
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
		RETURNING integration_type, owner_id
	`, r.tableName)

	var resp entity.IntegrationUpdateStatusResponse
	err := r.db.QueryRow(ctx, query, req.Status, req.Status, req.ID).Scan(
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
		SELECT guid,   status, started_at, stoped_at
		FROM %s 
		WHERE owner_id = $1 AND deleted_at IS NULL
	`, r.tableName)

	var (
		startedAt     sql.NullTime
		stoppedAt     sql.NullTime
		res           entity.IntegrationGetResponse
	)

	err := r.db.QueryRow(ctx, query, req.BusinessId).Scan(
		&res.Guid,
		&res.Status,
		&startedAt,
		&stoppedAt,

	)
	if err != nil {
		return nil, r.db.Error(err)
	}

	if startedAt.Valid {
		res.StartedAt = startedAt.Time
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

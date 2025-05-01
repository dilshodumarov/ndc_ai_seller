package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sugurta/internal/pkg/postgres"

	"sugurta/internal/entity"
)

const (
	telegramTableName = "telegram_accaunt"
)

type telegramRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewTelegramRepo(db *postgres.Postgres) *telegramRepo {
	return &telegramRepo{
		tableName: telegramTableName,
		db:        db,
	}
}

func (r *telegramRepo) Create(ctx context.Context, req entity.CreateTelegramAccountRequest) (string, error) {
	query := `
		INSERT INTO telegram_accaunt (number, business_id)
		VALUES ($1, $2)
		RETURNING guid
	`
	var id string
	err := r.db.QueryRow(ctx, query, req.Number, req.BusinessID).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("telegramRepo - Create: %w", err)
	}
	return id, nil
}

// Get By ID
func (r *telegramRepo) Get(ctx context.Context, id string) (*entity.TelegramAccount, error) {
	query := `
		SELECT guid, number, business_id, status, created_at, updated_at
		FROM telegram_accaunt
		WHERE guid = $1
	`
	var acc entity.TelegramAccount
	err := r.db.QueryRow(ctx, query, id).Scan(
		&acc.ID, &acc.Number, &acc.BusinessID, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("telegramRepo - Get: %w", err)
	}
	return &acc, nil
}

// Update
func (r *telegramRepo) Update(ctx context.Context, req entity.UpdateTelegramAccountRequest) error {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Number != "" {
		setClauses = append(setClauses, fmt.Sprintf("number = $%d", argPos))
		args = append(args, req.Number)
		argPos++
	}
	if req.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, req.Status)
		argPos++
	}

	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	// Add updated_at always
	setClauses = append(setClauses, fmt.Sprintf("updated_at = CURRENT_TIMESTAMP"))

	query := fmt.Sprintf(`
		UPDATE telegram_accaunt
		SET %s
		WHERE number = $%d
	`, strings.Join(setClauses, ", "), argPos)

	args = append(args, req.Phone)

	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("telegramRepo - Update: %w", err)
	}
	return nil
}

// Delete
func (r *telegramRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM telegram_accaunt WHERE guid = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("telegramRepo - Delete: %w", err)
	}
	return nil
}

// List (optional: business_id orqali filter bilan)
func (r *telegramRepo) List(ctx context.Context, businessID string) ([]entity.TelegramAccount, error) {
	query := `
		SELECT guid, number, business_id, status, created_at, updated_at
		FROM telegram_accaunt
		WHERE business_id = $1
	`
	rows, err := r.db.Query(ctx, query, businessID)
	if err != nil {
		return nil, fmt.Errorf("telegramRepo - List: %w", err)
	}
	defer rows.Close()

	var accounts []entity.TelegramAccount
	for rows.Next() {
		var acc entity.TelegramAccount
		if err := rows.Scan(
			&acc.ID, &acc.Number, &acc.BusinessID, &acc.Status, &acc.CreatedAt, &acc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("telegramRepo - List Scan: %w", err)
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

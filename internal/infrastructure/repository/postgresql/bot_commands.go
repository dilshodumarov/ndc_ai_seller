package postgresql

import (
	"context"
	"fmt"
	"strings"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"

	"github.com/jackc/pgx"
)

type BotCommandsRepo struct {
	db *postgres.Postgres
}

func NewBotCommandsRepo(db *postgres.Postgres) *BotCommandsRepo {
	return &BotCommandsRepo{
		db: db,
	}
}

func (p *BotCommandsRepo) CreateBotCommand(ctx context.Context, cmd entity.BotCommandRequest) (string, error) {
	query := `
		INSERT INTO bot_commands (integration_id, command, response)
		VALUES ($1, $2, $3)
		RETURNING guid
	`

	var guid string
	err := p.db.QueryRow(ctx, query, cmd.IntegrationID, cmd.Command, cmd.Response).Scan(&guid)
	if err != nil {
		return "", fmt.Errorf("CreateBotCommand: %w", err)
	}

	return guid, nil
}

func (p *BotCommandsRepo) GetBotCommand(ctx context.Context, guid string) (*entity.BotCommandResponse, error) {
	query := `
		SELECT guid, integration_id, command, response, created_at, updated_at
		FROM bot_commands
		WHERE guid = $1
	`

	var cmd entity.BotCommandResponse
	err := p.db.QueryRow(ctx, query, guid).Scan(
		&cmd.Guid, &cmd.IntegrationID, &cmd.Command, &cmd.Response, &cmd.CreatedAt, &cmd.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetBotCommand: %w", err)
	}

	return &cmd, nil
}

func (p *BotCommandsRepo) UpdateBotCommand(ctx context.Context, cmd entity.BotCommandUpdateRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argID := 1
	if cmd.Command != "" {
		setParts = append(setParts, fmt.Sprintf("command = $%d", argID))
		args = append(args, cmd.Command)
		argID++
	}

	if cmd.Response != "" {
		setParts = append(setParts, fmt.Sprintf("response = $%d", argID))
		args = append(args, cmd.Response)
		argID++
	}

	// Agar hech narsa update qilinmasa — chiqib ket
	if len(setParts) == 0 {
		return fmt.Errorf("UpdateBotCommand: no fields to update")
	}

	// updated_at har doim update qilinadi
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")

	query := fmt.Sprintf(`
		UPDATE bot_commands
		SET %s
		WHERE guid = $%d
	`, strings.Join(setParts, ", "), argID)

	args = append(args, cmd.Guid)

	_, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("UpdateBotCommand exec error: %w", err)
	}

	return nil
}

func (p *BotCommandsRepo) DeleteBotCommand(ctx context.Context, guid string) error {
	query := `DELETE FROM bot_commands WHERE guid = $1`

	_, err := p.db.Exec(ctx, query, guid)
	if err != nil {
		return fmt.Errorf("DeleteBotCommand: %w", err)
	}

	return nil
}

func (p *BotCommandsRepo) ListBotCommands(ctx context.Context, integrationID string) ([]entity.BotCommandResponse, error) {
	query := `
		SELECT guid, integration_id, command, response, created_at, updated_at
		FROM bot_commands
		WHERE integration_id = $1
	`

	rows, err := p.db.Query(ctx, query, integrationID)
	if err != nil {
		return nil, fmt.Errorf("ListBotCommands: %w", err)
	}
	defer rows.Close()

	var commands []entity.BotCommandResponse
	for rows.Next() {
		var cmd entity.BotCommandResponse
		err := rows.Scan(
			&cmd.Guid, &cmd.IntegrationID, &cmd.Command, &cmd.Response, &cmd.CreatedAt, &cmd.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListBotCommands Scan: %w", err)
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

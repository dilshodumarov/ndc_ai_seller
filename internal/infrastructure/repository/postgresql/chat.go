package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
)

type ChatRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewChatRepo(db *postgres.Postgres) *ChatRepo {
	return &ChatRepo{
		tableName: "chat_history",
		db:        db,
	}
}

func (p *ChatRepo) Create(ctx context.Context, h *entity.ChatHistory) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			integration_id, message, chat_id, platform_id, ai_response, reply_to_message_id
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING guid, created_at, updated_at`, p.tableName)

	err := p.db.QueryRow(ctx, query,
		h.IntegrationID,
		h.Message,
		h.ChatID,
		h.PlatformID,
		h.AIResponse,
		h.ReplyToMessageID,
	).Scan(&h.GUID, &h.CreatedAt, &h.UpdatedAt)

	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p *ChatRepo) Get(ctx context.Context, guid string) (*entity.ChatHistory, error) {
	query := fmt.Sprintf(`
		SELECT guid, integration_id, message, chat_id, platform_id, ai_response, reply_to_message_id, created_at, updated_at
		FROM %s WHERE guid = $1`, p.tableName)

	var h entity.ChatHistory
	err := p.db.QueryRow(ctx, query, guid).Scan(
		&h.GUID, &h.IntegrationID, &h.Message, &h.ChatID,
		&h.PlatformID, &h.AIResponse, &h.ReplyToMessageID,
		&h.CreatedAt, &h.UpdatedAt,
	)

	if err != nil {
		return nil, p.db.Error(err)
	}

	return &h, nil
}

func (p *ChatRepo) List(ctx context.Context, chatID int64) ([]*entity.ChatHistory, error) {
	query := fmt.Sprintf(`
		SELECT guid, integration_id, message, chat_id, platform_id, ai_response, reply_to_message_id, created_at, updated_at
		FROM %s WHERE chat_id = $1 ORDER BY created_at DESC`, p.tableName)

	rows, err := p.db.Query(ctx, query, chatID)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var results []*entity.ChatHistory
	for rows.Next() {
		var h entity.ChatHistory
		if err := rows.Scan(
			&h.GUID, &h.IntegrationID, &h.Message, &h.ChatID,
			&h.PlatformID, &h.AIResponse, &h.ReplyToMessageID,
			&h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		results = append(results, &h)
	}

	return results, nil
}

func (p *ChatRepo) Update(ctx context.Context, h *entity.ChatHistory) error {
	setClauses := []string{}
	args := []interface{}{}
	argID := 1

	if h.Message != "" {
		setClauses = append(setClauses, fmt.Sprintf("message = $%d", argID))
		args = append(args, h.Message)
		argID++
	}

	if h.AIResponse != "" {
		setClauses = append(setClauses, fmt.Sprintf("ai_response = $%d", argID))
		args = append(args, h.AIResponse)
		argID++
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now().UTC())
	argID++

	args = append(args, h.GUID)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE guid = $%d", p.tableName, strings.Join(setClauses, ", "), argID)

	_, err := p.db.Exec(ctx, query, args...)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

func (p *ChatRepo) Delete(ctx context.Context, guid string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE guid = $1", p.tableName)

	_, err := p.db.Exec(ctx, query, guid)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

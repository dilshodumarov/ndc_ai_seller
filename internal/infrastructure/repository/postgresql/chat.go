package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"

	"github.com/jackc/pgx/v5"
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
	query := `
		INSERT INTO chat_history (
			message_id, business_id, phone, platform_id, chat_id, message, ai_response, reply_to_message_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING guid, created_at, updated_at
	`

	err := p.db.QueryRow(ctx, query,
		h.MessageId,
		h.BusinessId,
		h.Phone,
		h.PlatformID,
		h.ChatID,
		h.Message,
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
		SELECT guid, business_id, message, chat_id, platform_id, ai_response, reply_to_message_id, created_at, updated_at
		FROM %s WHERE guid = $1`, p.tableName)

	var h entity.ChatHistory
	err := p.db.QueryRow(ctx, query, guid).Scan(
		&h.GUID, &h.BusinessId, &h.Message, &h.ChatID,
		&h.PlatformID, &h.AIResponse, &h.ReplyToMessageID,
		&h.CreatedAt, &h.UpdatedAt,
	)

	if err != nil {
		return nil, p.db.Error(err)
	}

	return &h, nil
}

func (p *ChatRepo) List(ctx context.Context, req *entity.ListChatHistoryRequest) ([]*entity.SendMessage, error) {
	baseQuery := `
		SELECT message_id, business_id, message, ai_response, chat_id, platform, reply_to_message_id, created_at
		FROM chat_history
		WHERE chat_id = $1 AND business_id = $2
		ORDER BY created_at DESC
	`

	var (
		rows pgx.Rows
		err  error
	)

	if req.Limit > 0 {
		baseQuery += " LIMIT $3"
		rows, err = p.db.Query(ctx, baseQuery, req.ChatID, req.BusinessID, req.Limit)
	} else {
		rows, err = p.db.Query(ctx, baseQuery, req.ChatID, req.BusinessID)
	}

	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var results []*entity.SendMessage

	for rows.Next() {
		var (
			messageID        int
			businessIDStr    sql.NullString
			message          sql.NullString
			aiResponse       sql.NullString
			chatIDValue      int64
			platform         sql.NullString
			replyToMessageID sql.NullInt64
			createdAt        sql.NullTime
		)

		if err := rows.Scan(
			&messageID, &businessIDStr, &message, &aiResponse,
			&chatIDValue, &platform, &replyToMessageID, &createdAt,
		); err != nil {
			return nil, p.db.Error(err)
		}

		resp := &entity.SendMessage{
			MessageId: messageID,
			Chatid:    chatIDValue,
			From:      "",
		}

		if businessIDStr.Valid {
			resp.BusinessId = businessIDStr.String
		}
		if platform.Valid {
			resp.Platform = platform.String
		}
		if createdAt.Valid {
			resp.Timestamp = createdAt.Time.Format(time.RFC3339)
		}
		if replyToMessageID.Valid {
			resp.ReplyToMessageID = int(replyToMessageID.Int64)
		}

		if message.Valid && message.String != "" {
			resp.Message = message.String
			resp.From = "user"
		}
		if aiResponse.Valid && aiResponse.String != "" {
			resp.AIResponse = aiResponse.String
			resp.From = "ai"
		}

		results = append(results, resp)
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

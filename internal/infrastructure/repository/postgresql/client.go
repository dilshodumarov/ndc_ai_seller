package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sugurta/internal/entity"

	"github.com/jackc/pgx"
)

func (r *userRepo) ListClients(ctx context.Context, filter entity.ClientFilter) (*entity.ListClients, error) {
	query := `
		SELECT guid, platform_id, first_name, phone,location,is_block, created_at,
		       user_name, from_chanel, order_status, goal
		FROM client
		WHERE 1=1
	`

	args := []interface{}{}
	argIdx := 1

	// Filtrlash
	if filter.Name != "" {
		query += fmt.Sprintf(" AND first_name ILIKE $%d", argIdx)
		args = append(args, "%"+filter.Name+"%")
		argIdx++
	}
	if filter.Phone != "" {
		query += fmt.Sprintf(" AND phone ILIKE $%d", argIdx)
		args = append(args, "%"+filter.Phone+"%")
		argIdx++
	}
	if filter.From != "" {
		query += fmt.Sprintf(" AND from_chanel ILIKE $%d", argIdx)
		args = append(args, "%"+filter.From+"%")
		argIdx++
	}
	if filter.Goal != "" {
		query += fmt.Sprintf(" AND goal ILIKE $%d", argIdx)
		args = append(args, "%"+filter.Goal+"%")
		argIdx++
	}
	if filter.OrderStatus != "" {
		query += fmt.Sprintf(" AND order_status ILIKE $%d", argIdx)
		args = append(args, "%"+filter.OrderStatus+"%")
		argIdx++
	}

	// Pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()

	var clients []entity.Client
	for rows.Next() {
		var c entity.Client
		var userName, from, goal, orderStatus sql.NullString

		err := rows.Scan(
			&c.ID,
			&c.PlatformID,
			&c.FirstName,

			&c.Phone,
			&c.Location,
			&c.IsBlock,
			&c.CreatedAt,

			&userName,
			&from,
			&orderStatus,
			&goal,
		)
		if err != nil {
			return nil, r.db.Error(err)
		}

		if userName.Valid {
			c.UserName = userName.String
		}
		if from.Valid {
			c.From = from.String
		}
		if orderStatus.Valid {
			c.OrderStatus = orderStatus.String
		}
		if goal.Valid {
			c.Goal = goal.String
		}

		clients = append(clients, c)
	}

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM client
		WHERE 1=1
	`
	argsCount := []interface{}{}
	countArgIdx := 1

	if filter.Name != "" {
		countQuery += fmt.Sprintf(" AND first_name ILIKE $%d", countArgIdx)
		argsCount = append(argsCount, "%"+filter.Name+"%")
		countArgIdx++
	}
	if filter.Phone != "" {
		countQuery += fmt.Sprintf(" AND phone ILIKE $%d", countArgIdx)
		argsCount = append(argsCount, "%"+filter.Phone+"%")
		countArgIdx++
	}
	if filter.From != "" {
		countQuery += fmt.Sprintf(" AND from_chanel ILIKE $%d", countArgIdx)
		argsCount = append(argsCount, "%"+filter.From+"%")
		countArgIdx++
	}
	if filter.Goal != "" {
		countQuery += fmt.Sprintf(" AND goal ILIKE $%d", countArgIdx)
		argsCount = append(argsCount, "%"+filter.Goal+"%")
		countArgIdx++
	}
	if filter.OrderStatus != "" {
		countQuery += fmt.Sprintf(" AND order_status ILIKE $%d", countArgIdx)
		argsCount = append(argsCount, "%"+filter.OrderStatus+"%")
		countArgIdx++
	}

	var totalCount int
	err = r.db.QueryRow(ctx, countQuery, argsCount...).Scan(&totalCount)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &entity.ListClients{
		Clients: clients,
		Count:   totalCount,
		Page:    page,
		Limit:   limit,
	}, nil
}

func (r *userRepo) GetClientByID(ctx context.Context, id string) (*entity.Client, error) {
	query := `
		SELECT guid, platform_id, first_name, phone, location,is_block,created_at, 
		       user_name, from_chanel, order_status, goal
		FROM client
		WHERE guid = $1
	`

	var c entity.Client
	var userName, from, orderStatus, goal sql.NullString

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.PlatformID,
		&c.FirstName,

		&c.Phone,
		&c.Location,
		&c.IsBlock,
		&c.CreatedAt,

		&userName,
		&from,
		&orderStatus,
		&goal,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("client with ID %s not found", id)
		}
		return nil, r.db.Error(err)
	}

	if userName.Valid {
		c.UserName = userName.String
	}
	if from.Valid {
		c.From = from.String
	}
	if orderStatus.Valid {
		c.OrderStatus = orderStatus.String
	}
	if goal.Valid {
		c.Goal = goal.String
	}

	return &c, nil
}

func (r *userRepo) BlockUser(ctx context.Context, req entity.BlockUser) error {
	query := `
		UPDATE client
		SET is_block = $1
		WHERE bussnes_id = $2 AND platform_id = $3
	`

	_, err := r.db.Exec(ctx, query, req.Block, req.BusinessID, req.PlatformId)
	if err != nil {
		return fmt.Errorf("userRepo - BlockUser - Exec: %w", err)
	}

	return nil
}

func (r *userRepo) PauzChat(ctx context.Context, req entity.PauzeChat) error {
	query := `
		UPDATE client
		SET is_pauze = $1
		WHERE bussnes_id = $2 AND platform_id = $3 and from_chanel=$4
	`
	_, err := r.db.Exec(ctx, query, req.Pauze, req.BusinessID, req.PlatformId, req.Type)
	if err != nil {
		return fmt.Errorf("userRepo - BlockUser - Exec: %w", err)
	}

	return nil
}

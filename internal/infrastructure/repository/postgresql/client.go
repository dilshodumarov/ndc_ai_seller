package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sugurta/internal/entity"
	"time"

	"github.com/jackc/pgx"
)

func (r *userRepo) ListClients(ctx context.Context, filter entity.ClientFilter) (*entity.ListClients, error) {
	fmt.Println(filter)
	query := `
		SELECT guid, platform_id, client_id,first_name, phone,location_text,location,is_block, created_at,
		       user_name, from_chanel, order_status, goal
		FROM client
		WHERE bussnes_id=$1
	`

	args := []interface{}{}
	argIdx := 1
	args = append(args, filter.BussinesID)
	argIdx++
	// Filtrlash
	if filter.ClientId > 0 {
		query += fmt.Sprintf(" AND client_id = $%d", argIdx)
		args = append(args, filter.ClientId)
		argIdx++
	}
	if filter.Search != "" {
		searchParam := "%" + filter.Search + "%"
		query += fmt.Sprintf(` AND (
			first_name ILIKE $%d OR 
			phone ILIKE $%d OR 
			user_name ILIKE $%d OR 
			from_chanel ILIKE $%d OR 
			order_status ILIKE $%d OR 
			location ILIKE $%d OR 
			goal ILIKE $%d
		)`, argIdx, argIdx, argIdx, argIdx, argIdx, argIdx, argIdx)

		args = append(args, searchParam)
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

		var platformID, firstName sql.NullString
		var phone, location, location_text sql.NullString
		var userName, from, goal, orderStatus sql.NullString
		var createdAt sql.NullTime
		var isBlock sql.NullBool
		var clientId sql.NullInt64
		var id sql.NullString

		err := rows.Scan(
			&id,
			&platformID,
			&clientId,
			&firstName,
			&phone,
			&location_text,
			&location,
			&isBlock,
			&createdAt,
			&userName,
			&from,
			&orderStatus,
			&goal,
		)
		if err != nil {
			return nil, r.db.Error(err)
		}

		if id.Valid {
			c.ID = id.String
		}
		if platformID.Valid {
			c.PlatformID = platformID.String
		}
		if clientId.Valid {
			c.ClientId = int(clientId.Int64)
		}
		if firstName.Valid {
			c.FirstName = firstName.String
		}
		if phone.Valid {
			c.Phone = phone.String
		}
		if location.Valid {
			c.Location = location.String
		}
		if location_text.Valid {
			c.LocationText = location_text.String
		}
		if isBlock.Valid {
			c.IsBlock = isBlock.Bool
		}
		if createdAt.Valid {
			c.CreatedAt = createdAt.Time
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
		WHERE bussnes_id=$1
	`
	argsCount := []interface{}{}
	countArgIdx := 1
	argsCount = append(argsCount, filter.BussinesID)
	countArgIdx++
	if filter.Search != "" {
		searchParam := "%" + filter.Search + "%"
		countQuery += fmt.Sprintf(` AND (
			first_name ILIKE $%d OR 
			phone ILIKE $%d OR 
			user_name ILIKE $%d OR 
			from_chanel ILIKE $%d OR 
			order_status ILIKE $%d OR 
			location ILIKE $%d OR 
			goal ILIKE $%d
		)`, countArgIdx, countArgIdx, countArgIdx, countArgIdx, countArgIdx, countArgIdx, countArgIdx)

		argsCount = append(argsCount, searchParam)
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
		Items: clients,
		Count: totalCount,
		Page:  page,
		Limit: limit,
	}, nil
}

func (r *userRepo) GetClientByID(ctx context.Context, id string) (*entity.Client, error) {
	query := `
		SELECT guid, platform_id, client_id,first_name, phone, location_text,location,is_block,created_at, 
		       user_name, from_chanel, order_status, goal
		FROM client
		WHERE guid = $1
	`

	var c entity.Client
	var (
		userName, from, orderStatus, goal sql.NullString
		phone, location, location_text    sql.NullString
	)

	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.PlatformID,
		&c.ClientId,
		&c.FirstName,

		&phone,
		&location_text,
		&location,
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

	// Tekshirishlar
	if phone.Valid {
		c.Phone = phone.String
	}
	if location.Valid {
		c.Location = location.String
	}
	if location_text.Valid {
		c.LocationText = location_text.String
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

func (r *userRepo) BlockUser(ctx context.Context, req entity.UpdateUser) error {
	setParts := []string{}
	args := []interface{}{}
	argID := 1

	if req.FirstName != "" {
		setParts = append(setParts, fmt.Sprintf("first_name = $%d", argID))
		args = append(args, req.FirstName)
		argID++
	}
	if req.Phone != "" {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argID))
		args = append(args, req.Phone)
		argID++
	}
	if req.UserName != "" {
		setParts = append(setParts, fmt.Sprintf("user_name = $%d", argID))
		args = append(args, req.UserName)
		argID++
	}
	if req.OrderStatus != "" {
		setParts = append(setParts, fmt.Sprintf("order_status = $%d", argID))
		args = append(args, req.OrderStatus)
		argID++
	}
	if req.Location != "" {
		setParts = append(setParts, fmt.Sprintf("location = $%d", argID))
		args = append(args, req.Location)
		argID++
	}
	if req.LocationText != "" {
		setParts = append(setParts, fmt.Sprintf("location_text = $%d", argID))
		args = append(args, req.LocationText)
		argID++
	}
	if req.Goal != "" {
		setParts = append(setParts, fmt.Sprintf("goal = $%d", argID))
		args = append(args, req.Goal)
		argID++
	}
	if req.Block != nil {
		setParts = append(setParts, fmt.Sprintf("is_block = $%d", argID))
		args = append(args, *req.Block)
		argID++
	}

	// always update updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now())
	argID++

	// WHERE conditions
	args = append(args, req.BusinessID, req.Id)

	query := fmt.Sprintf(`
		UPDATE client
		SET %s
		WHERE bussnes_id = $%d AND guid = $%d
	`, strings.Join(setParts, ", "), argID, argID+1)

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("userRepo - BlockUser - Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("userRepo - BlockUser: no rows updated")
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

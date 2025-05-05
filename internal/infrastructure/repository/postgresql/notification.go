package postgresql

import (
	"context"
	"fmt"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
	"time"
)

const (
	notificationTableName = "notifications"
)

type notificationRepo struct {
	db *postgres.Postgres
}

func NewNotificationRepo(db *postgres.Postgres) *notificationRepo {
	return &notificationRepo{
		db: db,
	}
}


func (r *notificationRepo) Create(ctx context.Context, req *entity.CreateNotificationRequest) (string, error) {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, title, message, type)
		VALUES ($1, $2, $3, $4)
		RETURNING guid`, notificationTableName)

	var id string
	err := r.db.QueryRow(ctx, query,
		req.UserID, req.Title, req.Message, req.Type).Scan(&id)

	if err != nil {
		return "", r.db.Error(err)
	}
	return id, nil
}


func (r *notificationRepo) Get(ctx context.Context, id string) (*entity.GetNotification, error) {
	query := fmt.Sprintf(`SELECT guid, user_id, title, message, type, is_read, created_at, read_at 
		FROM %s WHERE guid = $1`, notificationTableName)

	var res entity.GetNotification
	var readAt *time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(
		&res.GUID,
		&res.UserID,
		&res.Title,
		&res.Message,
		&res.Type,
		&res.IsRead,
		&res.CreatedAt,
		&readAt,
	)
	if err != nil {
		return nil, r.db.Error(err)
	}

	if readAt != nil {
		res.ReadAt = *readAt
	}

	return &res, nil
}


func (r *notificationRepo) Update(ctx context.Context, req *entity.UpdateNotificationRequest) error {
	query := fmt.Sprintf(`UPDATE %s 
		SET title = $1, message = $2, type = $3, is_read = $4, read_at = $5 
		WHERE guid = $6`, notificationTableName)

	_, err := r.db.Exec(ctx, query,
		req.Title,
		req.Message,
		req.Type,
		req.IsRead,
		req.ReadAt,
		req.GUID,
	)

	if err != nil {
		return r.db.Error(err)
	}
	return nil
}


func (r *notificationRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE guid = $1`, notificationTableName)

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.db.Error(err)
	}
	return nil
}


func (r *notificationRepo) List(ctx context.Context, req entity.ListNotificationRequest) (*entity.ListNotification, error) {
	query := fmt.Sprintf(`
		SELECT guid, user_id, title, message, type, is_read, created_at, read_at 
		FROM %s 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`, notificationTableName)

	var args []interface{}
	args = append(args, req.UserID)


	limit := req.Limit
	if limit <= 0 {
		limit = 10 
	}
	offset := (req.Page - 1) * limit
	if req.Page <= 0 {
		offset = 0
	}

	query += " LIMIT $2 OFFSET $3"
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()

	var notifications []entity.GetNotification

	for rows.Next() {
		var n entity.GetNotification
		var readAt *time.Time

		err := rows.Scan(
			&n.GUID,
			&n.UserID,
			&n.Title,
			&n.Message,
			&n.Type,
			&n.IsRead,
			&n.CreatedAt,
			&readAt,
		)
		if err != nil {
			return nil, r.db.Error(err)
		}

		if readAt != nil {
			n.ReadAt = *readAt
		}

		notifications = append(notifications, n)
	}

	return &entity.ListNotification{
		Notification: notifications,
		Count:  len(notifications),
	}, nil
}



func (r *notificationRepo) MarkAsRead(ctx context.Context, id string) error {
	query := fmt.Sprintf(`UPDATE %s SET is_read = true, read_at = now() WHERE guid = $1`, notificationTableName)
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.db.Error(err)
	}
	return nil
}

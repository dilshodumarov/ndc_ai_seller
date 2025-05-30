package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
)

type databaseRepo struct {
	db *postgres.Postgres
}

func NewDatabaseRepo(db *postgres.Postgres) *databaseRepo {
	return &databaseRepo{db: db}
}

func (r *databaseRepo) Create(ctx context.Context, req *entity.CreateDatabaseRequest) (string, error) {
	query := `
		INSERT INTO "database" (name, description, tokens)
		VALUES ($1, $2, $3)
		RETURNING guid
	`
	var guid string
	err := r.db.QueryRow(ctx, query, req.Name, req.Description, req.Tokens).Scan(&guid)
	if err != nil {
		return "", fmt.Errorf("create database: %w", err)
	}
	return guid, nil
}

func (r *databaseRepo) GetByID(ctx context.Context, guid string) (*entity.Database, error) {
	query := `
		SELECT guid, name, description, tokens, created_at, updated_at
		FROM "database"
		WHERE guid = $1
	`
	var (
		name        sql.NullString
		description sql.NullString
		tokens      sql.NullInt32
		db          entity.Database
	)
	err := r.db.QueryRow(ctx, query, guid).Scan(
		&db.Guid,
		&name,
		&description,
		&tokens,
		&db.CreatedAt,
		&db.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("get database by id: %w", err)
	}

	if name.Valid {
		db.Name = &name.String
	}
	if description.Valid {
		db.Description = &description.String
	}
	if tokens.Valid {
		val := int(tokens.Int32)
		db.Tokens = &val
	}

	return &db, nil
}

func (r *databaseRepo) Update(ctx context.Context, req *entity.UpdateDatabaseRequest) error {
	query := `
		UPDATE "database"
		SET name = $1, description = $2, tokens = $3, updated_at = CURRENT_TIMESTAMP
		WHERE guid = $4
	`
	_, err := r.db.Exec(ctx, query, req.Name, req.Description, req.Tokens, req.Guid)
	if err != nil {
		return fmt.Errorf("update database: %w", err)
	}
	return nil
}

func (r *databaseRepo) Delete(ctx context.Context, guid string) error {
	query := `DELETE FROM "database" WHERE guid = $1`
	_, err := r.db.Exec(ctx, query, guid)
	if err != nil {
		return fmt.Errorf("delete database: %w", err)
	}
	return nil
}

func (r *databaseRepo) List(ctx context.Context, filter *entity.Filter) (*entity.Databaselist, error) {
	var (
		args     []interface{}
		argID    = 1
		where    = "WHERE 1=1"
		limit    = 10
		offset   = 0
	)

	// Search filter
	if filter.Search != "" {
		where += fmt.Sprintf(` AND (name ILIKE $%d OR description ILIKE $%d)`, argID, argID+1)
		searchTerm := "%" + filter.Search + "%"
		args = append(args, searchTerm, searchTerm)
		argID += 2
	}

	// Limit and Page
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	offset = (filter.Page - 1) * limit

	// 🟡 1. Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "database" %s`, where)

	var count int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&count); err != nil {
		return nil, fmt.Errorf("count databases: %w", err)
	}

	// 🟢 2. Data query with LIMIT + OFFSET
	dataQuery := fmt.Sprintf(`
		SELECT guid, name, description, tokens, created_at, updated_at
		FROM "database"
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argID, argID+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("list databases: %w", err)
	}
	defer rows.Close()

	results := &entity.Databaselist{
		Page:  filter.Page,
		Limit: limit,
		Count: count,
	}

	for rows.Next() {
		var (
			name        sql.NullString
			description sql.NullString
			tokens      sql.NullInt32
			db          entity.Database
		)

		err := rows.Scan(
			&db.Guid,
			&name,
			&description,
			&tokens,
			&db.CreatedAt,
			&db.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan database row: %w", err)
		}

		if name.Valid {
			db.Name = &name.String
		}
		if description.Valid {
			db.Description = &description.String
		}
		if tokens.Valid {
			val := int(tokens.Int32)
			db.Tokens = &val
		}

		results.Items = append(results.Items, db)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

package postgresql

import (
	"context"
	"fmt"
	"time"

	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
)

const (
	clientTypeTableName = "client_type"
)

type clientTypeRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewClientTypeRepo(db *postgres.Postgres) *clientTypeRepo {
	return &clientTypeRepo{
		tableName: clientTypeTableName,
		db:        db,
	}
}

// CreateClientType -.
func (p *clientTypeRepo) Create(ctx context.Context, clientType *entity.CreateClientTypeRequest) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, created_at, updated_at) VALUES ($1, $2, $3)`, p.tableName)

	now := time.Now()
	_, err := p.db.Exec(ctx, query, clientType.Name, now, now)
	if err != nil {
		return p.db.Error(err)
	}

	return nil
}

// GetClientType -.
func (p *clientTypeRepo) Get(ctx context.Context, id string) (*entity.ClientType, error) {
	var clientType entity.ClientType

	query := fmt.Sprintf(`SELECT id, name, created_at, updated_at FROM %s WHERE id = $1`, p.tableName)

	err := p.db.QueryRow(ctx, query, id).Scan(
		&clientType.ID,
		&clientType.Name,
		&clientType.CreatedAt,
		&clientType.UpdatedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return &clientType, nil
}

// List -.
func (p *clientTypeRepo) List(ctx context.Context, limit, page int) ([]*entity.ClientType, error) {
	var clientTypes []*entity.ClientType

	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT id, name, created_at, updated_at FROM %s ORDER BY created_at DESC LIMIT $1 OFFSET $2`, p.tableName)

	rows, err := p.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ct entity.ClientType
		err = rows.Scan(&ct.ID, &ct.Name, &ct.CreatedAt, &ct.UpdatedAt)
		if err != nil {
			return nil, p.db.Error(err)
		}
		clientTypes = append(clientTypes, &ct)
	}

	return clientTypes, nil
}

// Update -.
func (p *clientTypeRepo) Update(ctx context.Context, clientType *entity.UpdateClientType) error {
	query := fmt.Sprintf(`UPDATE %s SET name = $1, updated_at = $2 WHERE id = $3`, p.tableName)

	result, err := p.db.Exec(ctx, query, clientType.Name, time.Now().UTC(), clientType.ID)
	if err != nil {
		return p.db.Error(err)
	}

	if result.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no rows updated"))
	}

	return nil
}

// DeleteClientType -.
func (p *clientTypeRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, p.tableName)

	result, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return p.db.Error(err)
	}

	if result.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no rows deleted"))
	}

	return nil
}

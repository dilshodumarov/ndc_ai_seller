package postgresql

import (
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
)

const (
	roleTableName      = "role"
	roleServiceName    = "roleService"
	roleSpanRepoPrefix = "roleRepo"
)

// roleRepo -.
type roleRepo struct {
	tableName string
	db        *postgres.Postgres
}

// NewroleRepo -.
func NewRoleRepo(db *postgres.Postgres) *roleRepo {
	return &roleRepo{
		tableName: userTableName,
		db:        db,
	}
}

func (r *roleRepo) roleSelectQueryPrefix() squirrel.SelectBuilder {
	return r.db.Sq.Builder.
		Select(
			"id",
			"name",
			"client_type_id",
			"created_at",
			"updated_at",
		).From(r.tableName)
}

// CreateRole -.
func (r *roleRepo) Create(ctx context.Context, role *entity.CreateRoleRequest) error {

	data := map[string]any{
		"name":           role.Name,
		"client_type_id": role.ClientTypeId,
		"created_at":     time.Now(),
		"updated_at":     time.Now(),
	}

	query, args, err := r.db.Sq.Builder.Insert(r.tableName).SetMap(data).ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "create"))
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return r.db.Error(err)
	}

	return nil
}

// GetRole -.
func (r *roleRepo) Get(ctx context.Context, params map[string]string) (*entity.Role, error) {
	var role entity.Role

	// Build SELECT with JOIN
	queryBuilder := r.db.Builder.
		Select(
			"role.id",
			"role.name",
			"role.client_type_id",
			"role.created_at",
			"role.updated_at",
			"client_type.id AS client_type_id",
			"client_type.name AS client_type_name",
		).
		From("role").
		Join("client_type ON client_type.id = role.client_type_id")

	// Dynamically add WHERE conditions from params
	for key, value := range params {
		if key == "id" || key == "name" {
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(key, value))
		}
	}

	// Build query and args
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s get", r.tableName))
	}

	// Execute and scan
	row := r.db.QueryRow(ctx, query, args...)
	err = row.Scan(
		&role.ID,
		&role.Name,
		&role.ClientTypeId,
		&role.CreatedAt,
		&role.UpdatedAt,
		&role.ClientTypeData.ID,
		&role.ClientTypeData.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("roleRepo - Get - row.Scan: %w", err)
	}

	return &role, nil
}
func (r *roleRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.RoleListResponse, error) {
	var (
		roles []entity.RoleResponse
	)
	queryBuilder := r.roleSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	for key, value := range filter {
		if key == "id" || key == "name" || key == "client_type_id" {
			queryBuilder = queryBuilder.Where(r.db.Sq.Equal(key, value))
			continue
		}
	}

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "list"))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()
	roles = make([]entity.RoleResponse, 0)
	for rows.Next() {
		var role entity.RoleResponse
		if err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.ClientTypeId,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, r.db.Error(err)
		}
		roles = append(roles, role)
	}

	countBuilder := r.db.Sq.Builder.
		Select("COUNT(*)").
		From(r.tableName)

	// Apply same filters
	for key, value := range filter {
		if key == "id" || key == "name" || key == "client_type_id" {
			countBuilder = countBuilder.Where(r.db.Sq.Equal(key, value))
		}
	}

	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "count"))
	}

	var totalCount uint64
	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, r.db.Error(err)
	}

	return &entity.RoleListResponse{
		Items: roles,
		Count: totalCount,
	}, nil
}

// UpdateRole -.
func (r *roleRepo) Update(ctx context.Context, role *entity.UpdateRoleRequest) error {

	clauses := map[string]any{
		"name":           role.Name,
		"client_type_id": role.ClientTypeId,
		"updated_at":     time.Now(),
	}
	sqlStr, args, err := r.db.Sq.Builder.
		Update(r.tableName).
		SetMap(clauses).
		Where(r.db.Sq.Equal("id", role.ID)).
		ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.tableName+" update")
	}

	commandTag, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return r.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil

}

// DeleteRole -.
func (r *roleRepo) Delete(ctx context.Context, id string) error {
	sqlStr, args, err := r.db.Sq.Builder.Delete(r.tableName).Where(r.db.Sq.Equal("id", id)).ToSql()
	if err != nil {
		return r.db.ErrSQLBuild(err, r.tableName+" delete")
	}

	commandTag, err := r.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return r.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return r.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

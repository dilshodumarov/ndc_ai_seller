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
	userTableName      = "user"
	userServiceName    = "userService"
	userSpanRepoPrefix = "userRepo"
)

type userRepo struct {
	tableName string
	db        *postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *userRepo {
	return &userRepo{
		tableName: userTableName,
		db:        db,
	}
}
func (p *userRepo) userSelectQueryPrefix() squirrel.SelectBuilder {
	return p.db.Builder.
		Select(
			"'guid'",
			"'first_name'",
			"'last_name'",
			"'email'",
			"'phone_number'",
			"'password'",
			"'role_id'",
			"'is_active'",
			"'created_at'",
			"'updated_at'",
		).From(p.tableName)
}

// CreateUser -.
func (p *userRepo) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	data := map[string]interface{}{
		"first_name":   u.FirstName,
		"last_name":    u.LastName,
		"email":        u.Email,
		"phone_number": u.PhoneNumber,
		"password":     u.Password,
		"role_id":      u.RoleID,
		"is_active":    u.IsActive,
	}

	// Build the INSERT query with RETURNING *
	queryBuilder := p.db.Builder.
		Insert(p.tableName).
		SetMap(data).
		Suffix("RETURNING id, first_name, last_name, email, phone_number, password, role_id, is_active, created_at, updated_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}

	// Prepare a variable to hold the returned user
	var insertedUser entity.User

	// QueryRow to fetch the returned data
	row := p.db.QueryRow(ctx, query, args...)
	err = row.Scan(
		&insertedUser.ID,
		&insertedUser.FirstName,
		&insertedUser.LastName,
		&insertedUser.Email,
		&insertedUser.PhoneNumber,
		&insertedUser.Password,
		&insertedUser.RoleID,
		&insertedUser.IsActive,
		&insertedUser.CreatedAt,
		&insertedUser.UpdatedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return &insertedUser, nil
}

func (p *userRepo) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	var (
		user entity.User
	)

	fmt.Println("params: ", params)

	// ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Get")
	// defer span.End()

	queryBuilder := p.userSelectQueryPrefix()

	for key, value := range params {
		switch key {
		case "id":
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal("'guid'", value))
		case "email":
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal("'email'", value))
		case "phone_number":
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal("'phone_number'", value))
		}
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}

	fmt.Println("sql: ", query)
	if err = p.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		fmt.Println("errr: ", p.db.Error(err))
		return nil, p.db.Error(err)
	}

	return &user, nil
}

func (p *userRepo) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error) {
	// ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"List")
	// defer span.End()

	var (
		users []*entity.User
	)
	queryBuilder := p.userSelectQueryPrefix()

	if limit != 0 {
		queryBuilder = queryBuilder.Limit(limit).Offset(offset)
	}

	for key, value := range filter {
		if key == "is_active" || key == "role_id" {
			queryBuilder = queryBuilder.Where(p.db.Sq.Equal(key, value))
			continue
		}
		if key == "created_at" {
			queryBuilder = queryBuilder.Where("created_at=?", value)
			continue
		}
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	users = make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		if err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PhoneNumber,
			&user.Password,
			&user.RoleID,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		users = append(users, &user)
	}

	return users, nil
}

func (p *userRepo) Update(ctx context.Context, user *entity.User) error {
	// ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Update")
	// defer span.End()

	clauses := map[string]any{
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"password":     user.Password,
		"role_id":      user.RoleID,
		"is_active":    user.IsActive,
		"updated_at":   user.UpdatedAt,
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("id", user.ID)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

func (p *userRepo) UpdatePassword(ctx context.Context, user *entity.UpdatePasswordRequest) error {
	clauses := map[string]any{
		"password":   user.Password,
		"updated_at": time.Now().UTC(),
	}
	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).
		SetMap(clauses).
		Where(p.db.Sq.Equal("email", user.Email)).
		ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" update")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

// DeleteUser -.
func (p *userRepo) Delete(ctx context.Context, guid string) error {
	// ctx, span := otlp.Start(ctx, userServiceName, userSpanRepoPrefix+"Delete")
	// defer span.End()

	sqlStr, args, err := p.db.Sq.Builder.Delete(p.tableName).Where(p.db.Sq.Equal("id", guid)).ToSql()
	if err != nil {
		return p.db.ErrSQLBuild(err, p.tableName+" delete")
	}

	commandTag, err := p.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return p.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return p.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}

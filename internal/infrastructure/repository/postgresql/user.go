package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sugurta/internal/entity"
	"sugurta/internal/pkg/postgres"
	"time"

	"github.com/jackc/pgx"
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

// CreateUser -.
func (p *userRepo) Create(ctx context.Context, u *entity.User) (*entity.User, error) {
	// Avval role name bo‘yicha role guid ni olish
	var roleID string
	err := p.db.QueryRow(ctx, `SELECT guid FROM "role" WHERE name = $1`, u.RoleData.Name).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, p.db.Error(err)
	}

	// User yaratish uchun query
	query := `
		INSERT INTO "user" 
			(first_name, last_name, email, phone_number, password, role_id, is_active)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING 
			guid, first_name, last_name, email, phone_number, password, role_id, is_active, created_at, updated_at;
	`

	var insertedUser entity.User

	err = p.db.QueryRow(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.PhoneNumber,
		u.Password,
		roleID,
		u.IsActive,
	).Scan(
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
		user                         entity.User
		businessID                   sql.NullString
		roleName                     sql.NullString
		clientTypeID                 sql.NullString
		roleCreatedAt, roleUpdatedAt sql.NullTime
	)

	var whereClause string
	var args []interface{}
	i := 1

	for key, value := range params {
		if whereClause != "" {
			whereClause += " AND "
		}
		switch key {
		case "id":
			whereClause += fmt.Sprintf("u.guid = $%d", i)
			args = append(args, value)
		case "email":
			whereClause += fmt.Sprintf("u.email = $%d", i)
			args = append(args, value)
		case "phone_number":
			whereClause += fmt.Sprintf("u.phone_number = $%d", i)
			args = append(args, value)
		}
		i++
	}

	query := fmt.Sprintf(`
		SELECT 
			u.guid, u.first_name, u.last_name, u.email, u.phone_number,
			u.password, u.role_id, u.is_active, u.created_at, u.updated_at,
			b.guid AS business_id,
			r.name AS role_name, r.client_type_id, r.created_at, r.updated_at
		FROM "user" u
		LEFT JOIN business b ON b.owner_id = u.guid AND b.deleted_at IS NULL
		LEFT JOIN role r ON r.guid = u.role_id
		WHERE %s
	`, whereClause)

	fmt.Println("query: ", query)

	err := p.db.QueryRow(ctx, query, args...).Scan(
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
		&businessID,
		&roleName,
		&clientTypeID,
		&roleCreatedAt,
		&roleUpdatedAt,
	)
	if err != nil {
		fmt.Println("errr: ", p.db.Error(err))
		return nil, p.db.Error(err)
	}

	if businessID.Valid {
		user.BusinessID = businessID.String
	}

	user.RoleData = entity.Role{
		ID:           user.RoleID,
		Name:         roleName.String,
		ClientTypeId: clientTypeID.String,
		CreatedAt:    roleCreatedAt.Time,
		UpdatedAt:    roleUpdatedAt.Time,
	}

	return &user, nil
}

func (p *userRepo) List(ctx context.Context, filter entity.UserFilter) ([]*entity.User, error) {
	var users []*entity.User

	// Asosiy so'rov
	query := `
		SELECT guid, first_name, last_name, email, phone_number, password, role_id, is_active, created_at, updated_at
		FROM "user"
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	// Filtrlash parametrlariga qo'shish
	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *filter.IsActive)
		argIdx++
	}

	if filter.RoleID != "" {
		query += fmt.Sprintf(" AND role_id = $%d", argIdx)
		args = append(args, filter.RoleID)
		argIdx++
	}

	if filter.CreatedAt != "" {
		query += fmt.Sprintf(" AND created_at = $%d", argIdx)
		args = append(args, filter.CreatedAt)
		argIdx++
	}

	// Pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	}
	page := filter.Offset
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	// So'rovni bajarish
	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	// Natijalarni yig'ish
	users = make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
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

func (p *userRepo) GetByIDs(ctx context.Context, id string) ([]*entity.User, error) {

	query := `
		SELECT 
			guid, first_name, last_name, email, phone_number, password, role_id, is_active, created_at, updated_at
		FROM "user"
		WHERE guid =$1
	`

	rows, err := p.db.Query(ctx, query, id)
	if err != nil {
		return nil, p.db.Error(err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
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
	query := `
		UPDATE "user"
		SET password = $1, updated_at = $2
		WHERE email = $3
	`

	result, err := p.db.Pool.Exec(ctx, query, user.Password, time.Now().UTC(), user.Email)
	if err != nil {
		return p.db.Error(err)
	}

	if result.RowsAffected() == 0 {
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

func (r *userRepo) ListClients(ctx context.Context, filter entity.ClientFilter) (*entity.ListClients, error) {
	// Asosiy so'rov
	query := `
		SELECT guid, platform_id, first_name, chat_id, bussnes_id, phone, created_at, updated_at
		FROM client
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	// Filtrlash parametrlariga qo'shish
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
		err := rows.Scan(
			&c.ID,
			&c.PlatformID,
			&c.FirstName,
			&c.ChatID,
			&c.BusinessID,
			&c.Phone,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, r.db.Error(err)
		}
		clients = append(clients, c)
	}
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


	var totalCount int
	err = r.db.QueryRow(ctx, countQuery, argsCount...).Scan(&totalCount)
	if err != nil {
		return nil, r.db.Error(err)
	}

	// Clients va umumiy sonni qaytarish
	return &entity.ListClients{
		Clients: clients,
		Count:   totalCount,
		Page:    page,
		Limit:   limit,
	}, nil
}

func (r *userRepo) GetClientByID(ctx context.Context, id string) (*entity.Client, error) {
	query := `
		SELECT guid, platform_id, first_name, chat_id, bussnes_id, phone, created_at, updated_at
		FROM client
		WHERE guid = $1
	`

	var c entity.Client
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.PlatformID,
		&c.FirstName,
		&c.ChatID,
		&c.BusinessID,
		&c.Phone,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("client with ID %s not found", id)
		}
		return nil, r.db.Error(err)
	}

	return &c, nil
}

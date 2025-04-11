package persistent

// import (
// 	"sugurta/internal/pkg/postgres"
// 	"context"
// 	"fmt"

// 	"sugurta/internal/entity"
// )

// type BusinessRepo struct {
// 	*postgres.Postgres
// }

// // NewBusinessRepo -.
// func NewBusinessRepo(pg *postgres.Postgres) *BusinessRepo {
// 	return &BusinessRepo{pg}
// }

// // CreateBusiness creates a new business record.
// func (r *BusinessRepo) CreateBusiness(ctx context.Context, b entity.CreateBusinessRequest) error {
// 	query := `
// 		INSERT INTO business (owner_id, name, description)
// 		VALUES ($1, $2, $3)
// 	`
// 	_, err := r.Pool.Exec(ctx, query, b.OwnerID, b.Name, b.Description)
// 	if err != nil {
// 		return fmt.Errorf("BusinessRepo - CreateBusiness - r.Pool.Exec: %w", err)
// 	}
// 	return nil
// }

// // GetBusiness retrieves a business record by ID.
// func (r *BusinessRepo) GetBusiness(ctx context.Context, id string) (*entity.Business, error) {
// 	var b entity.Business

// 	query := `
// 		SELECT id, owner_id, name, description, created_at, updated_at
// 		FROM business
// 		WHERE id = $1
// 	`
// 	row := r.Pool.QueryRow(ctx, query, id)
// 	err := row.Scan(&b.ID, &b.OwnerID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &b, nil
// }

// // UpdateBusiness updates an existing business record.
// func (r *BusinessRepo) UpdateBusiness(ctx context.Context, b entity.UpdateBusinessRequest) error {
// 	query := `
// 		UPDATE business
// 		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
// 		WHERE id = $3 and owner_id = $4
// 	`
// 	_, err := r.Pool.Exec(ctx, query, b.Name, b.Description, b.ID, b.OwnerID)
// 	return err
// }

// // DeleteBusiness deletes a business record by ID.
// func (r *BusinessRepo) DeleteBusiness(ctx context.Context, id, owner_id string) error {
// 	query := `
// 		DELETE FROM business
// 		WHERE id = $1 and owner_id = $2
// 	`
// 	_, err := r.Pool.Exec(ctx, query, id, owner_id)
// 	return err
// }

// // GetAllBusinesses retrieves all business records with pagination.
// func (r *BusinessRepo) GetAllBusinesses(ctx context.Context, req entity.GetAllBusinessesRequest) (*entity.GetAllBusinessesResponse, error) {
// 	var businesses entity.GetAllBusinessesResponse

// 	offset := req.Limit * req.Page

// 	query := `
// 		SELECT id, owner_id, name, description, created_at, updated_at
// 		FROM business
// 		WHERE owner_id = $1
// 		ORDER BY created_at DESC
// 		LIMIT $2 OFFSET $3
// 	`
// 	rows, err := r.Pool.Query(ctx, query, req.OwnerID, req.Limit, offset)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		b := entity.Business{}
// 		err := rows.Scan(&b.ID, &b.OwnerID, &b.Name, &b.Description, &b.CreatedAt, &b.UpdatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		businesses.Itmes = append(businesses.Itmes, b)
// 	}

// 	countQuery := `
// 		SELECT COUNT(1)
// 		FROM business
// 		WHERE owner_id = $1
// 	`
// 	var total int
// 	err = r.Pool.QueryRow(ctx, countQuery, req.OwnerID).Scan(&total)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get total count: %w", err)
// 	}

// 	businesses.Total = total

// 	return &businesses, nil
// }

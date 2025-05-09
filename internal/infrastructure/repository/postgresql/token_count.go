package postgresql

import (
	"context"
	"fmt"
	"sugurta/internal/entity"
)

func (r *integrationRepo) GetTokenUsageList(ctx context.Context, req *entity.IntegrationListRequest) (*entity.IntegrationListResponse, error) {
	query := `
		SELECT 
			guid,
			source_type,
			used_for,
			request_tokens,
			response_tokens,
			total_tokens,
			created_at
		FROM client_token_usage
		WHERE business_id = $1
	`
	var args []interface{}
	args = append(args, req.BusinessID)
	argID := 2

	// Filter: source_type
	if req.SourceType != "" {
		query += fmt.Sprintf(" AND source_type = $%d", argID)
		args = append(args, req.SourceType)
		argID++
	}

	// Filter: from_date
	if !req.FromDate.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argID)
		args = append(args, req.FromDate)
		argID++
	}

	// Filter: to_date
	if !req.ToDate.IsZero() {
		query += fmt.Sprintf(" AND created_at <= $%d", argID)
		args = append(args, req.ToDate)
		argID++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.Error(err)
	}
	defer rows.Close()

	var (
		usages      []entity.TokenUsage
		totalTokens int
	)

	for rows.Next() {
		var usage entity.TokenUsage
		err := rows.Scan(
			&usage.ID,
			&usage.SourceType,
			&usage.UsedFor,
			&usage.RequestTokens,
			&usage.ResponseTokens,
			&usage.TotalTokens,
			&usage.CreatedAt,
		)
		if err != nil {
			return nil, r.db.Error(err)
		}
		totalTokens += usage.TotalTokens
		usages = append(usages, usage)
	}

	return &entity.IntegrationListResponse{
		Usages:      usages,
		TotalTokens: totalTokens,
	}, nil
}

package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/korshunvladislav/splix/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	groupID *int,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    SELECT u.id, u.version, u.full_name, u.phone_number, u.email, u.created_at, u.updated_at
    FROM splix.users u
    %s
    ORDER BY u.id
    LIMIT $1
    OFFSET $2;
    `

	args := []any{limit, offset}

	if groupID != nil {
		query = fmt.Sprintf(query, `
        JOIN splix.group_members gm ON gm.user_id = u.id
        WHERE gm.group_id = $3
        `)
		args = append(args, groupID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
			&userModel.Email,
			&userModel.CreatedAt,
			&userModel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}

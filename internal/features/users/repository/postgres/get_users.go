package users_postgres_repository

import (
    "context"
    "fmt"
    "strings"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
    ctx context.Context,
    groupId *int,
    limit *int,
    offset *int,
) ([]domain.User, error) {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := strings.Builder{}

    query.WriteString(`
    SELECT u.id, u.version, u.full_name, u.phone_number, u.email, u.created_at, u.updated_at
    FROM splix.users u
    `)

    args := make([]interface{}, 0)
    argCounter := 1

    if groupId != nil {
        query.WriteString(` JOIN splix.group_members gm ON gm.user_id = u.id `)
        query.WriteString(fmt.Sprintf(" WHERE gm.group_id = $%d ", argCounter))
        args = append(args, *groupId)
        argCounter++
    }

    query.WriteString(` ORDER BY u.id ASC `)

    if limit != nil {
        query.WriteString(fmt.Sprintf(" LIMIT $%d ", argCounter))
        args = append(args, *limit)
        argCounter++
    }

    if offset != nil {
        query.WriteString(fmt.Sprintf(" OFFSET $%d ", argCounter))
        args = append(args, *offset)
    }

    query.WriteString(";")

    rows, err := r.pool.Query(ctx, query.String(), args...)
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

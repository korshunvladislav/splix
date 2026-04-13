package users_postgres_repository

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/korshunvladislav/splix/internal/core/domain"
    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func (r *UsersRepository) GetUser(
    ctx context.Context,
    id int,
) (domain.User, error) {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    SELECT id, version, full_name, phone_number, email, created_at, updated_at
    FROM splix.users
    WHERE id = $1
    `

    row := r.pool.QueryRow(ctx, query, id)

    var userModel UserModel

    err := row.Scan(
        &userModel.ID,
        &userModel.Version,
        &userModel.FullName,
        &userModel.PhoneNumber,
        &userModel.Email,
        &userModel.CreatedAt,
        &userModel.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return domain.User{}, fmt.Errorf(
                "user with id='%d': %w",
                id,
                core_errors.ErrNotFound,
            )
        }

        return domain.User{}, fmt.Errorf("scan error: %w", err)
    }

    userDomain := domain.NewUser(
        userModel.ID,
        userModel.Version,
        userModel.FullName,
        userModel.PhoneNumber,
        userModel.Email,
        userModel.CreatedAt,
        userModel.UpdatedAt,
    )

    return userDomain, nil
}

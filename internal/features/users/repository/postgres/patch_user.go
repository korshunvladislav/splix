package users_postgres_repository

import (
    "context"
    "errors"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/korshunvladislav/splix/internal/core/domain"
    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func (r *UsersRepository) PatchUser(
    ctx context.Context,
    id int,
    user domain.User,
) (domain.User, error) {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    UPDATE splix.users
    SET
        full_name = $1,
        phone_number = $2,
        email = $3,
        version = version + 1
    WHERE id = $4 AND version = $5
    RETURNING 
        id,
        version,
        full_name,
        phone_number,
        email,
        created_at,
        updated_at
    `

    row := r.pool.QueryRow(
        ctx,
        query,
        user.FullName,
        user.PhoneNumber,
        user.Email,
        id,
        user.Version,
    )

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
                "user with id='%d' concurrently accessed: %w",
                id,
                core_errors.ErrConflict,
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

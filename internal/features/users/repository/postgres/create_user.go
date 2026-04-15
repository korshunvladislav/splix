package users_postgres_repository

import (
    "context"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
    ctx context.Context,
    user domain.User,
) (domain.User, error) {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    INSERT INTO splix.users (full_name, phone_number, email)
    VALUES ($1, $2, $3)
    RETURNING id, version, full_name, phone_number, email, created_at, updated_at;
    `

    row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, user.Email)

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

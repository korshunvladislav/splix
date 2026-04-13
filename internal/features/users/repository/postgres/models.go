package users_postgres_repository

import (
    "time"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

type UserModel struct {
    ID          int
    Version     int
    FullName    string
    PhoneNumber *string
    Email       string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func userDomainsFromModels(users []UserModel) []domain.User {
    userDomains := make([]domain.User, len(users))

    for i, user := range users {
        userDomains[i] = domain.NewUser(
            user.ID,
            user.Version,
            user.FullName,
            user.PhoneNumber,
            user.Email,
            user.CreatedAt,
            user.UpdatedAt,
        )
    }
    return userDomains
}

package groups_postgres_repository

import (
    "context"
    "errors"
    "fmt"

    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
    core_postgres_pool "github.com/korshunvladislav/splix/internal/core/repository/postgres/pool"
)

func (r *GroupsRepository) AddMember(
    ctx context.Context,
    groupID int,
    userID int,
) error {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    INSERT INTO splix.group_members (group_id, user_id)
    VALUES ($1, $2)
    `

    _, err := r.pool.Exec(ctx, query, groupID, userID)
    if err != nil {
        switch {
        case errors.Is(err, core_postgres_pool.ErrViolatesForeignKey):
            return fmt.Errorf(
                "group with id='%d' or user with id='%d': %w",
                groupID,
                userID,
                core_errors.ErrNotFound,
            )

        case errors.Is(err, core_postgres_pool.ErrViolatesUnique):
            return fmt.Errorf(
                "user with id='%d' is already a member of group with id='%d': %w",
                userID,
                groupID,
                core_errors.ErrConflict,
            )

        default:
            return fmt.Errorf("insert group member: %w", err)
        }
    }

    return nil
}

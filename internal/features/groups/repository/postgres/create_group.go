package groups_postgres_repository

import (
    "context"
    "errors"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
    core_postgres_pool "github.com/korshunvladislav/splix/internal/core/repository/postgres/pool"
)

func (r *GroupsRepository) CreateGroup(
    ctx context.Context,
    group domain.Group,
) (domain.Group, error) {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    WITH new_group AS (
        INSERT INTO splix.groups (
            title, description, currency, creator_user_id
        )
        VALUES ($1, $2, $3, $4)
        RETURNING *
    ),
    new_member AS (
        INSERT INTO splix.group_members (group_id, user_id)
        SELECT id, creator_user_id
        FROM new_group
    )
    SELECT id, version, title, description, currency,
        creator_user_id, created_at, updated_at
    FROM new_group;
    `

    row := r.pool.QueryRow(
        ctx,
        query,
        group.Title,
        group.Description,
        group.Currency,
        group.CreatorUserID,
    )

    var groupModel GroupModel
    err := row.Scan(
        &groupModel.ID,
        &groupModel.Version,
        &groupModel.Title,
        &groupModel.Description,
        &groupModel.Currency,
        &groupModel.CreatorUserID,
        &groupModel.CreatedAt,
        &groupModel.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
            return domain.Group{}, fmt.Errorf(
                "%v: user with id='%d': %w",
                err,
                group.CreatorUserID,
                core_errors.ErrNotFound,
            )
        }

        return domain.Group{}, fmt.Errorf("scan error: %w", err)
    }

    groupDomain := groupDomainFromModel(groupModel)

    return groupDomain, nil
}

package groups_postgres_repository

import (
    "context"
    "errors"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
    core_postgres_pool "github.com/korshunvladislav/splix/internal/core/repository/postgres/pool"
)

func (r *GroupsRepository) GetGroup(
    ctx context.Context,
    id int,
) (domain.Group, error) {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    SELECT id, version, title, description, currency, creator_user_id, created_at, updated_at
    FROM splix.groups
    WHERE id = $1;
    `

    row := r.pool.QueryRow(ctx, query, id)

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
        if errors.Is(err, core_postgres_pool.ErrNoRows) {
            return domain.Group{}, fmt.Errorf(
                "group with id='%d': %w",
                id,
                core_errors.ErrNotFound,
            )
        }

        return domain.Group{}, fmt.Errorf("scan error: %w", err)
    }

    groupDomain := groupDomainFromModel(groupModel)

    return groupDomain, nil
}

package groups_postgres_repository

import (
    "context"
    "fmt"

    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func (r *GroupsRepository) DeleteGroup(
    ctx context.Context,
    id int,
) error {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    DELETE FROM splix.groups
    WHERE id = $1;
    `

    cmdTag, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("exec query: %w", err)
    }
    if cmdTag.RowsAffected() == 0 {
        return fmt.Errorf(
            "group with id='%d': %w",
            id,
            core_errors.ErrNotFound,
        )
    }

    return nil
}

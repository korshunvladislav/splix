package groups_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func (r *GroupsRepository) DeleteMember(
	ctx context.Context,
	groupID int,
	userID int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM splix.group_members
	WHERE group_id = $1
	AND user_id = $2;
	`

	cmdTag, err := r.pool.Exec(ctx, query, groupID, userID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"membership for user with id='%d' in group with id='%d': %w",
			userID,
			groupID,
			core_errors.ErrNotFound,
		)
	}

	return nil
}

package groups_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/korshunvladislav/splix/internal/core/domain"
	core_errors "github.com/korshunvladislav/splix/internal/core/errors"
	core_postgres_pool "github.com/korshunvladislav/splix/internal/core/repository/postgres/pool"
)

func (r *GroupsRepository) PatchGroup(
	ctx context.Context,
	id int,
	group domain.Group,
) (domain.Group, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    UPDATE splix.groups
    SET
        title = $1,
        description = $2,
        currency = $3,
        version = version + 1
    WHERE id = $4 AND version = $5
    RETURNING
        id,
        version,
        title,
        description,
        currency,
        creator_user_id,
        created_at,
        updated_at;
    `

	row := r.pool.QueryRow(
		ctx,
		query,
		group.Title,
		group.Description,
		group.Currency,
		id,
		group.Version,
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Group{}, fmt.Errorf(
				"group with id='%d' concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}

		return domain.Group{}, fmt.Errorf("scan error: %w", err)
	}

	groupDomain := groupDomainFromModel(groupModel)

	return groupDomain, nil
}

package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/korshunvladislav/splix/internal/core/errors"
	core_postgres_pool "github.com/korshunvladislav/splix/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
    DELETE FROM splix.users
    WHERE id = $1;
    `

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return fmt.Errorf(
				"cannot delete user with id='%d': referenced by other records: %w",
				id,
				core_errors.ErrConflict,
			)
		}

		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}

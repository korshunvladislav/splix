package groups_postgres_repository

import core_postgres_pool "github.com/korshunvladislav/splix/internal/core/repository/postgres/pool"

type GroupsRepository struct {
    pool core_postgres_pool.Pool
}

func NewGroupsRepository(
    pool core_postgres_pool.Pool,
) *GroupsRepository {
    return &GroupsRepository{
        pool: pool,
    }
}

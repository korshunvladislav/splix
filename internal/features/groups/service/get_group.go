package groups_service

import (
    "context"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

func (s *GroupsService) GetGroup(
    ctx context.Context,
    id int,
) (domain.Group, error) {
    group, err := s.groupsRepository.GetGroup(ctx, id)
    if err != nil {
        return domain.Group{}, fmt.Errorf("get group from repository: %w", err)
    }

    return group, nil
}

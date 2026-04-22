package groups_service

import (
    "context"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

func (s *GroupsService) CreateGroup(
    ctx context.Context,
    group domain.Group,
) (domain.Group, error) {
    if err := group.Validate(); err != nil {
        return domain.Group{}, fmt.Errorf("validate group domain: %w", err)
    }

    group, err := s.groupsRepository.CreateGroup(ctx, group)
    if err != nil {
        return domain.Group{}, fmt.Errorf("create group: %w", err)
    }

    return group, nil
}

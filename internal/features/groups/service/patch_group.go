package groups_service

import (
    "context"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

func (s *GroupsService) PatchGroup(
    ctx context.Context,
    id int,
    patch domain.GroupPatch,
) (domain.Group, error) {
    group, err := s.groupsRepository.GetGroup(ctx, id)
    if err != nil {
        return domain.Group{}, fmt.Errorf("get group: %w", err)
    }

    if err := group.ApplyPatch(patch); err != nil {
        return domain.Group{}, fmt.Errorf("apply group patch: %w", err)
    }

    patchedGroup, err := s.groupsRepository.PatchGroup(ctx, id, group)
    if err != nil {
        return domain.Group{}, fmt.Errorf("patch group: %w", err)
    }

    return patchedGroup, nil
}

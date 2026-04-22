package groups_service

import (
    "context"
    "fmt"
)

func (s *GroupsService) DeleteGroup(
    ctx context.Context,
    id int,
) error {
    if err := s.groupsRepository.DeleteGroup(ctx, id); err != nil {
        return fmt.Errorf("delete group from repository: %w", err)
    }

    return nil
}

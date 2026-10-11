package groups_service

import (
    "context"
    "fmt"

    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func (s *GroupsService) AddMember(
    ctx context.Context,
    groupID int,
    userID int,
) error {
    if groupID <= 0 {
        return fmt.Errorf(
            "groupID must be positive: %w",
            core_errors.ErrInvalidArgument,
        )
    }

    if userID <= 0 {
        return fmt.Errorf(
            "userID must be positive: %w",
            core_errors.ErrInvalidArgument,
        )
    }

    err := s.groupsRepository.AddMember(ctx, groupID, userID)
    if err != nil {
        return fmt.Errorf("add member to group: %w", err)
    }

    return nil
}

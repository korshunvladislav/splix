package groups_service

import (
    "context"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func (s *GroupsService) GetGroups(
    ctx context.Context,
    userID *int,
    limit *int,
    offset *int,
) ([]domain.Group, error) {
    if limit != nil && *limit < 0 {
        return nil, fmt.Errorf(
            "'limit' must be a non-negative integer: %w",
            core_errors.ErrInvalidArgument,
        )
    }

    if offset != nil && *offset < 0 {
        return nil, fmt.Errorf(
            "'offset' must be a non-negative integer: %w",
            core_errors.ErrInvalidArgument,
        )
    }

    groups, err := s.groupsRepository.GetGroups(ctx, userID, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("get groups from repository: %w", err)
    }

    return groups, nil
}

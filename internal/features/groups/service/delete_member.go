package groups_service

import (
	"context"
	"fmt"

	core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func (s *GroupsService) DeleteMember(
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

	groupDomain, err := s.groupsRepository.GetGroup(ctx, groupID)
	if err != nil {
		return fmt.Errorf("get group: %w", err)
	}
	if userID == groupDomain.CreatorUserID {
		return fmt.Errorf(
			"cannot remove creator of group with id='%d': %w",
			groupID,
			core_errors.ErrConflict,
		)
	}

	if err = s.groupsRepository.DeleteMember(ctx, groupID, userID); err != nil {
		return fmt.Errorf("delete member from groups repository: %w", err)
	}

	return nil
}

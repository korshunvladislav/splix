package groups_service

import (
    "context"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

type GroupsService struct {
    groupsRepository GroupsRepository
}

type GroupsRepository interface {
    CreateGroup(
        ctx context.Context,
        group domain.Group,
    ) (domain.Group, error)

    GetGroups(
        ctx context.Context,
        userID *int,
        limit *int,
        offset *int,
    ) ([]domain.Group, error)

    GetGroup(
        ctx context.Context,
        id int,
    ) (domain.Group, error)

    DeleteGroup(
        ctx context.Context,
        id int,
    ) error

    PatchGroup(
        ctx context.Context,
        id int,
        group domain.Group,
    ) (domain.Group, error)
}

func NewGroupsService(
    groupsRepository GroupsRepository,
) *GroupsService {
    return &GroupsService{
        groupsRepository: groupsRepository,
    }
}

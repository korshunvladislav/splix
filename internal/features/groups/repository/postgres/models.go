package groups_postgres_repository

import (
    "time"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

type GroupModel struct {
    ID            int
    Version       int
    Title         string
    Description   *string
    Currency      string
    CreatorUserID int
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

func groupDomainFromModel(groupModel GroupModel) domain.Group {
    return domain.NewGroup(
        groupModel.ID,
        groupModel.Version,
        groupModel.Title,
        groupModel.Description,
        groupModel.Currency,
        groupModel.CreatorUserID,
        groupModel.CreatedAt,
        groupModel.UpdatedAt,
    )
}

func groupDomainsFromModels(groups []GroupModel) []domain.Group {
    groupDomains := make([]domain.Group, len(groups))

    for i, group := range groups {
        groupDomains[i] = groupDomainFromModel(group)
    }

    return groupDomains
}

package groups_transport_http

import (
    "time"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

type GroupDTOResponse struct {
    ID            int       `json:"id"`
    Version       int       `json:"version"`
    Title         string    `json:"title"`
    Description   *string   `json:"description"`
    Currency      string    `json:"currency"`
    CreatorUserID int       `json:"creator_user_id"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

func groupDTOFromDomain(group domain.Group) GroupDTOResponse {
    return GroupDTOResponse{
        ID:            group.ID,
        Version:       group.Version,
        Title:         group.Title,
        Description:   group.Description,
        Currency:      group.Currency,
        CreatorUserID: group.CreatorUserID,
        CreatedAt:     group.CreatedAt,
        UpdatedAt:     group.UpdatedAt,
    }
}

func groupsDTOFromDomains(groups []domain.Group) []GroupDTOResponse {
    groupsDTO := make([]GroupDTOResponse, len(groups))
    for i, group := range groups {
        groupsDTO[i] = groupDTOFromDomain(group)
    }

    return groupsDTO
}

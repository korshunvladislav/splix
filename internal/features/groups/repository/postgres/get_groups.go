package groups_postgres_repository

import (
    "context"
    "fmt"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

func (r *GroupsRepository) GetGroups(
    ctx context.Context,
    userID *int,
    limit *int,
    offset *int,
) ([]domain.Group, error) {
    ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
    defer cancel()

    query := `
    SELECT g.id, g.version, g.title, g.description, g.currency, g.creator_user_id, g.created_at, g.updated_at
    FROM splix.groups g
    %s
    ORDER BY g.id
    LIMIT $1
    OFFSET $2;
    `

    args := []any{limit, offset}

    if userID != nil {
        query = fmt.Sprintf(query, `
        JOIN splix.group_members gm ON gm.group_id = g.id
        WHERE gm.user_id = $3
        `)
        args = append(args, userID)
    } else {
        query = fmt.Sprintf(query, "")
    }

    rows, err := r.pool.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("select groups: %w", err)
    }
    defer rows.Close()

    var groupModels []GroupModel

    for rows.Next() {
        var groupModel GroupModel

        err := rows.Scan(
            &groupModel.ID,
            &groupModel.Version,
            &groupModel.Title,
            &groupModel.Description,
            &groupModel.Currency,
            &groupModel.CreatorUserID,
            &groupModel.CreatedAt,
            &groupModel.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("scan groups: %w", err)
        }

        groupModels = append(groupModels, groupModel)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("next rows: %w", err)
    }

    groupDomains := groupDomainsFromModels(groupModels)

    return groupDomains, nil
}

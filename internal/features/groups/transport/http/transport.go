package groups_transport_http

import (
	"context"
	"net/http"

	"github.com/korshunvladislav/splix/internal/core/domain"
	core_http_server "github.com/korshunvladislav/splix/internal/core/transport/http/server"
)

type GroupsHTTPHandler struct {
	groupsService GroupsService
}

type GroupsService interface {
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
		patch domain.GroupPatch,
	) (domain.Group, error)

	AddMember(
		ctx context.Context,
		groupID int,
		userID int,
	) error

	DeleteMember(
		ctx context.Context,
		groupID int,
		userID int,
	) error
}

func NewGroupsHTTPHandler(
	groupsService GroupsService,
) *GroupsHTTPHandler {
	return &GroupsHTTPHandler{
		groupsService: groupsService,
	}
}

func (h *GroupsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/groups",
			Handler: h.CreateGroup,
		},
		{
			Method:  http.MethodGet,
			Path:    "/groups",
			Handler: h.GetGroups,
		},
		{
			Method:  http.MethodGet,
			Path:    "/groups/{id}",
			Handler: h.GetGroup,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/groups/{id}",
			Handler: h.DeleteGroup,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/groups/{id}",
			Handler: h.PatchGroup,
		},
		{
			Method:  http.MethodPost,
			Path:    "/groups/{group_id}/members",
			Handler: h.AddMember,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/groups/{group_id}/members/{user_id}",
			Handler: h.DeleteMember,
		},
	}
}

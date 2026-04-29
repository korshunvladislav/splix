package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/korshunvladislav/splix/internal/core/logger"
	core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
	core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	groupID, limit, offset, err := getGroupIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'group_id'/'limit'/'offset' query params",
		)

		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, groupID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'users'",
		)

		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getGroupIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		groupIDQueryParamKey = "group_id"
		limitQueryParamKey   = "limit"
		offsetQueryParamKey  = "offset"
	)

	groupID, err := core_http_request.GetIntQueryParam(r, groupIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'group_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return groupID, limit, offset, nil
}

package users_transport_http

import (
    "fmt"
    "net/http"

    core_logger "github.com/korshunvladislav/splix/internal/core/logger"
    core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
    core_http_utils "github.com/korshunvladislav/splix/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    log := core_logger.FromContext(ctx)
    responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

    log.Debug("invoke GetUsers handler")

    groupId, limit, offset, err := getGroupIdLimitOffsetQueryParams(r)
    if err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to get 'group_id'/'limit'/'offset' query params",
        )

        return
    }

    userDomains, err := h.usersService.GetUsers(ctx, groupId, limit, offset)
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

func getGroupIdLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
    groupId, err := core_http_utils.GetIntQueryParam(r, "group_id")
    if err != nil {
        return nil, nil, nil, fmt.Errorf("get 'group_id' query param: %w", err)
    }

    limit, err := core_http_utils.GetIntQueryParam(r, "limit")
    if err != nil {
        return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
    }

    offset, err := core_http_utils.GetIntQueryParam(r, "offset")
    if err != nil {
        return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
    }

    return groupId, limit, offset, nil
}

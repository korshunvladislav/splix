package groups_transport_http

import (
    "fmt"
    "net/http"

    core_logger "github.com/korshunvladislav/splix/internal/core/logger"
    core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
    core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
)

type GetGroupsResponse []GroupDTOResponse

func (h *GroupsHTTPHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    log := core_logger.FromContext(ctx)
    responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

    userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
    if err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to get 'user_id'/'limit'/'offset' query params",
        )

        return
    }

    groupDomains, err := h.groupsService.GetGroups(ctx, userID, limit, offset)
    if err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to get 'groups'",
        )

        return
    }

    response := GetGroupsResponse(groupsDTOFromDomains(groupDomains))

    responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
    const (
        userIDQueryParamKey = "user_id"
        limitQueryParamKey  = "limit"
        offsetQueryParamKey = "offset"
    )

    userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
    }

    limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
    }

    offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
    }

    return userID, limit, offset, nil
}

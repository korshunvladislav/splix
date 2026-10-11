package groups_transport_http

import (
    "net/http"

    core_logger "github.com/korshunvladislav/splix/internal/core/logger"
    core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
    core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
)

type AddMemberRequest struct {
    UserID int `json:"user_id" validate:"required,gte=1"`
}

func (h *GroupsHTTPHandler) AddMember(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    log := core_logger.FromContext(ctx)
    responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

    groupID, err := core_http_request.GetIntPathValue(r, "group_id")
    if err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to get groupID path value",
        )

        return
    }

    var request AddMemberRequest
    if err = core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to decode and validate HTTP request",
        )

        return
    }

    err = h.groupsService.AddMember(ctx, groupID, request.UserID)
    if err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to add member to group",
        )

        return
    }

    responseHandler.EmptyResponse(http.StatusCreated)
}

package groups_transport_http

import (
    "net/http"

    core_logger "github.com/korshunvladislav/splix/internal/core/logger"
    core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
    core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
)

func (h *GroupsHTTPHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    log := core_logger.FromContext(ctx)
    responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

    groupID, err := core_http_request.GetIntPathValue(r, "id")
    if err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to get groupID path value",
        )

        return
    }

    if err := h.groupsService.DeleteGroup(ctx, groupID); err != nil {
        responseHandler.ErrorResponse(
            err,
            "failed to delete group",
        )

        return
    }

    responseHandler.NoContentResponse()
}

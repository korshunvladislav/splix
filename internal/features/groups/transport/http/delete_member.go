package groups_transport_http

import (
	"net/http"

	core_logger "github.com/korshunvladislav/splix/internal/core/logger"
	core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
	core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
)

func (h *GroupsHTTPHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
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

	userID, err := core_http_request.GetIntPathValue(r, "user_id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	if err = h.groupsService.DeleteMember(ctx, groupID, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete member",
		)

		return
	}

	responseHandler.NoContentResponse()
}

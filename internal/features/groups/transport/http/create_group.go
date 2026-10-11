package groups_transport_http

import (
	"net/http"

	"github.com/korshunvladislav/splix/internal/core/domain"
	core_logger "github.com/korshunvladislav/splix/internal/core/logger"
	core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
	core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
)

type CreateGroupRequest struct {
	Title         string  `json:"title" validate:"required,min=1,max=100"`
	Description   *string `json:"description" validate:"omitempty,min=1,max=500"`
	Currency      string  `json:"currency" validate:"omitempty,iso4217"`
	CreatorUserID int     `json:"creator_user_id" validate:"required"`
}

type CreateGroupResponse GroupDTOResponse

func (h *GroupsHTTPHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateGroupRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	groupDomain := domainFromDTO(request)

	groupDomain, err := h.groupsService.CreateGroup(ctx, groupDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create group",
		)

		return
	}

	response := CreateGroupResponse(groupDTOFromDomain(groupDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateGroupRequest) domain.Group {
	currency := dto.Currency
	if currency == "" {
		currency = "RUB"
	}

	return domain.NewGroupUninitialized(dto.Title, dto.Description, currency, dto.CreatorUserID)
}

package groups_transport_http

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/korshunvladislav/splix/internal/core/domain"
	core_logger "github.com/korshunvladislav/splix/internal/core/logger"
	core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
	core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
	core_http_types "github.com/korshunvladislav/splix/internal/core/transport/http/types"
)

type PatchGroupRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Currency    core_http_types.Nullable[string] `json:"currency"`
}

func (r *PatchGroupRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` length must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 500 {
				return fmt.Errorf("`Description` length must be between 1 and 500 symbols")
			}
		}
	}

	if r.Currency.Set {
		if r.Currency.Value == nil {
			return fmt.Errorf("`Currency` can't be NULL")
		}

		currencyLen := len([]rune(*r.Currency.Value))
		if currencyLen != 3 {
			return fmt.Errorf("`Currency` length must be 3 symbols")
		}

		if matched, _ := regexp.MatchString(`^[A-Z]{3}$`, *r.Currency.Value); !matched {
			return fmt.Errorf("`Currency` must be 3 uppercase letters (ISO 4217)")
		}
	}

	return nil
}

type PatchGroupResponse GroupDTOResponse

func (h *GroupsHTTPHandler) PatchGroup(w http.ResponseWriter, r *http.Request) {
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

	var request PatchGroupRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	groupPatch := groupPatchFromRequest(request)

	groupDomain, err := h.groupsService.PatchGroup(ctx, groupID, groupPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch group",
		)

		return
	}

	response := PatchGroupResponse(groupDTOFromDomain(groupDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func groupPatchFromRequest(request PatchGroupRequest) domain.GroupPatch {
	return domain.NewGroupPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Currency.ToDomain(),
	)
}

package expenses_transport_http

import (
    "net/http"

    "github.com/korshunvladislav/splix/internal/core/domain"
    core_logger "github.com/korshunvladislav/splix/internal/core/logger"
    core_http_request "github.com/korshunvladislav/splix/internal/core/transport/http/request"
    core_http_response "github.com/korshunvladislav/splix/internal/core/transport/http/response"
    expenses_service "github.com/korshunvladislav/splix/internal/features/expenses/service"
    "github.com/shopspring/decimal"
)

type CreateExpenseRequest struct {
    Description *string         `json:"description" validate:"omitempty,min=1,max=200"`
    Date        string          `json:"date" validate:"required,datetime=2006-01-02"`
    Amount      decimal.Decimal `json:"amount" validate:"required"`
    SplitType   string          `json:"split_type" validate:"required,oneof=EQUAL PERCENTAGE CUSTOM"`
    PayerUserID int             `json:"payer_user_id" validate:"required"`
    Splits      []SplitRequest  `json:"splits" validate:"required,dive"`
}

type SplitRequest struct {
    UserID     int              `json:"user_id" validate:"required"`
    Share      *decimal.Decimal `json:"share,omitempty"`
    Percentage *decimal.Decimal `json:"percentage,omitempty"`
}

type CreateExpenseResponse ExpenseDTOResponse

func (h *ExpensesHTTPHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
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

    var request CreateExpenseRequest
    if err = core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
        responseHandler.ErrorResponse(err, "failed to decode and validate request")

        return
    }

    expenseDomain, splits, err := domainFromDTO(request, groupID)

    expenseDomain, err = h.expensesService.CreateExpense(ctx, expenseDomain, splits)
    if err != nil {
        responseHandler.ErrorResponse(err, "failed to create expense")

        return
    }

    response := CreateExpenseResponse(expenseDTOFromDomain(expenseDomain))

    responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateExpenseRequest, groupID int) (domain.Expense, []expenses_service.SplitInput, error) {
    date, err := core_http_request.ParseDate(dto.Date)
    if err != nil {
        return domain.Expense{}, nil, err
    }

    expense := domain.NewExpenseUninitialized(
        dto.Description,
        date,
        dto.Amount,
        domain.SplitType(dto.SplitType),
        groupID,
        dto.PayerUserID,
        nil,
    )

    if len(dto.Splits) == 0 {
        return expense, nil, nil
    }

    splitInputs := make([]expenses_service.SplitInput, 0, len(dto.Splits))
    for _, split := range dto.Splits {
        splitInputs = append(splitInputs, expenses_service.SplitInput{
            UserID:     split.UserID,
            Share:      split.Share,
            Percentage: split.Percentage,
        })
    }

    return expense, splitInputs, nil
}

package expenses_transport_http

import (
    "context"

    "github.com/korshunvladislav/splix/internal/core/domain"
    core_http_server "github.com/korshunvladislav/splix/internal/core/transport/http/server"
    expenses_service "github.com/korshunvladislav/splix/internal/features/expenses/service"
)

type ExpensesHTTPHandler struct {
    expensesService ExpensesService
}

type ExpensesService interface {
    CreateExpense(
        ctx context.Context,
        expense domain.Expense,
        splits []expenses_service.SplitInput,
    ) (domain.Expense, error)
}

func NewExpensesHTTPHandler(
    expensesService ExpensesService,
) *ExpensesHTTPHandler {
    return &ExpensesHTTPHandler{
        expensesService: expensesService,
    }
}

func (h *ExpensesHTTPHandler) Routes() []core_http_server.Route {
    return []core_http_server.Route{}
}

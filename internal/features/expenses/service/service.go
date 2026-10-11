package expenses_service

import (
    "context"

    "github.com/korshunvladislav/splix/internal/core/domain"
)

type ExpensesService struct {
    expensesRepository ExpensesRepository
}

type ExpensesRepository interface {
    CreateExpense(
        ctx context.Context,
        expense domain.Expense,
    ) (domain.Expense, error)
}

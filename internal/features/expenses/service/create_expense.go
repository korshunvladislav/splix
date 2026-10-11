package expenses_service

import (
    "context"

    "github.com/korshunvladislav/splix/internal/core/domain"
    "github.com/shopspring/decimal"
)

type SplitInput struct {
    UserID     int
    Share      *decimal.Decimal
    Percentage *decimal.Decimal
}

type SplitInputs struct {
    splits      []SplitInput
    totalAmount decimal.Decimal
}

func (s *ExpensesService) CreateExpense(
    ctx context.Context,
    expense domain.Expense,
    splits []SplitInput,
) (domain.Expense, error) {

}

func (sp *SplitInputs) Validate(totalAmount decimal.Decimal) error {

}

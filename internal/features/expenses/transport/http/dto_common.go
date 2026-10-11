package expenses_transport_http

import (
    "time"

    "github.com/korshunvladislav/splix/internal/core/domain"
    "github.com/shopspring/decimal"
)

type ExpenseDTOResponse struct {
    ID          int                   `json:"id"`
    Version     int                   `json:"version"`
    Description *string               `json:"description"`
    Date        time.Time             `json:"date"`
    Amount      decimal.Decimal       `json:"amount"`
    SplitType   domain.SplitType      `json:"split_type"`
    GroupID     int                   `json:"group_id"`
    PayerUserID int                   `json:"payer_user_id"`
    Splits      []domain.ExpenseSplit `json:"splits"`
    CreatedAt   time.Time             `json:"created_at"`
    UpdatedAt   time.Time             `json:"updated_at"`
}

func expenseDTOFromDomain(expense domain.Expense) ExpenseDTOResponse {
    return ExpenseDTOResponse{
        ID:          expense.ID,
        Version:     expense.Version,
        Description: expense.Description,
        Date:        expense.Date,
        Amount:      expense.Amount,
        SplitType:   expense.SplitType,
        GroupID:     expense.GroupID,
        PayerUserID: expense.PayerUserID,
        Splits:      expense.Splits,
        CreatedAt:   expense.CreatedAt,
        UpdatedAt:   expense.UpdatedAt,
    }
}

package domain

import (
    "fmt"
    "time"

    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
    "github.com/shopspring/decimal"
)

type SplitType string

const (
    SplitTypeEqual      SplitType = "EQUAL"
    SplitTypePercentage SplitType = "PERCENTAGE"
    SplitTypeCustom     SplitType = "CUSTOM"
)

func (st SplitType) Validate() error {
    switch st {
    case SplitTypeEqual, SplitTypePercentage, SplitTypeCustom:
        return nil
    }
    return fmt.Errorf("invalid split_type: %s", st)
}

type ExpenseSplit struct {
    UserID     int
    Share      decimal.Decimal
    Percentage *decimal.Decimal
}

type Expense struct {
    ID      int
    Version int

    Description *string
    Date        time.Time
    Amount      decimal.Decimal
    SplitType   SplitType
    GroupID     int
    PayerUserID int
    Splits      []ExpenseSplit

    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewExpense(
    id int,
    version int,
    description *string,
    date time.Time,
    amount decimal.Decimal,
    splitType SplitType,
    groupID int,
    payerUserID int,
    splits []ExpenseSplit,
    createdAt time.Time,
    updatedAt time.Time,
) Expense {
    return Expense{
        ID:          id,
        Version:     version,
        Description: description,
        Date:        date,
        Amount:      amount,
        SplitType:   splitType,
        GroupID:     groupID,
        PayerUserID: payerUserID,
        Splits:      splits,
        CreatedAt:   createdAt,
        UpdatedAt:   updatedAt,
    }
}

func NewExpenseUninitialized(
    description *string,
    date time.Time,
    amount decimal.Decimal,
    splitType SplitType,
    groupID int,
    payerUserID int,
    splits []ExpenseSplit,
) Expense {
    return NewExpense(
        UnitializedID,
        UnitializedVersion,
        description,
        date,
        amount,
        splitType,
        groupID,
        payerUserID,
        splits,
        time.Time{},
        time.Time{},
    )
}

func (e *Expense) Validate() error {
    if e.Description != nil {
        descriptionLength := len([]rune(*e.Description))
        if descriptionLength < 1 || descriptionLength > 200 {
            return fmt.Errorf(
                "invalid `Description` len: %d: %w",
                descriptionLength,
                core_errors.ErrInvalidArgument,
            )
        }
    }

    today := time.Now().Truncate(24 * time.Hour)
    if e.Date.After(today) {
        return fmt.Errorf(
            "`Date` cannot be in the future: %s: %w",
            e.Date,
            core_errors.ErrInvalidArgument,
        )
    }

    if e.Amount.LessThanOrEqual(decimal.Zero) {
        return fmt.Errorf(
            "`Amount` must be greater than zero: %s: %w",
            e.Amount.String(),
            core_errors.ErrInvalidArgument,
        )
    }

    if err := e.SplitType.Validate(); err != nil {
        return err
    }

    if e.SplitType == SplitTypeCustom {
        if len(e.Splits) == 0 {
            return fmt.Errorf(
                "splits required for CUSTOM split type: %w",
                core_errors.ErrInvalidArgument,
            )
        }

        total := decimal.Zero
        for _, split := range e.Splits {
            if split.Share.LessThanOrEqual(decimal.Zero) {
                return fmt.Errorf(
                    "each share must be positive: %s: %w",
                    split.Share.String(),
                    core_errors.ErrInvalidArgument,
                )
            }
            total = total.Add(split.Share)
        }
        if !total.Equal(e.Amount) {
            return fmt.Errorf(
                "total shares %s != amount %s: %w",
                total.String(),
                e.Amount.String(),
                core_errors.ErrInvalidArgument,
            )
        }
    }

    return nil
}

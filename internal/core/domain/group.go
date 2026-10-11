package domain

import (
	"fmt"
	"regexp"
	"time"

	core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

type Group struct {
	ID      int
	Version int

	Title         string
	Description   *string
	Currency      string
	CreatorUserID int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewGroup(
	id int,
	version int,
	title string,
	description *string,
	currency string,
	creatorUserID int,
	createdAt time.Time,
	updatedAt time.Time,
) Group {
	return Group{
		ID:            id,
		Version:       version,
		Title:         title,
		Description:   description,
		Currency:      currency,
		CreatorUserID: creatorUserID,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func NewGroupUninitialized(
	title string,
	description *string,
	currency string,
	creatorUserID int,
) Group {
	return NewGroup(
		UnitializedID,
		UnitializedVersion,
		title,
		description,
		currency,
		creatorUserID,
		time.Time{},
		time.Time{},
	)
}

func (g *Group) Validate() error {
	titleLength := len([]rune(g.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if g.Description != nil {
		descriptionLength := len([]rune(*g.Description))
		if descriptionLength < 1 || descriptionLength > 500 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	re := regexp.MustCompile(`^[A-Z]{3}$`)
	if !re.MatchString(g.Currency) {
		return fmt.Errorf(
			"invalid `Currency`: %s: %w",
			g.Currency,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type GroupPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Currency    Nullable[string]
}

func NewGroupPatch(
	title Nullable[string],
	description Nullable[string],
	currency Nullable[string],
) GroupPatch {
	return GroupPatch{
		Title:       title,
		Description: description,
		Currency:    currency,
	}
}

func (p *GroupPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"`Title` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Currency.Set && p.Currency.Value == nil {
		return fmt.Errorf(
			"`Currency` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (g *Group) ApplyPatch(patch GroupPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate group patch: %w", err)
	}

	tmp := *g

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Currency.Set {
		tmp.Currency = *patch.Currency.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched group: %w", err)
	}

	*g = tmp

	return nil
}

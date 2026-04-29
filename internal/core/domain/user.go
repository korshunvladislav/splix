package domain

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"

	core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

type User struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
	Email       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
	email string,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Email:       email,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
	email string,
) User {
	return NewUser(
		UnitializedID,
		UnitializedVersion,
		fullName,
		phoneNumber,
		email,
		time.Time{},
		time.Time{},
	)
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLength := len([]rune(*u.PhoneNumber))
		if phoneNumberLength < 10 || phoneNumberLength > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLength,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	emailLength := len([]rune(u.Email))
	if emailLength < 6 || emailLength > 255 {
		return fmt.Errorf(
			"invalid `Email` len: %d: %w",
			emailLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf(
			"invalid `Email` format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
	Email       Nullable[string]
}

func NewUserPatch(
	fullName Nullable[string],
	phoneNumber Nullable[string],
	email Nullable[string],
) UserPatch {
	return UserPatch{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if p.Email.Set && p.Email.Value == nil {
		return fmt.Errorf(
			"`Email` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if patch.Email.Set {
		tmp.Email = *patch.Email.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}

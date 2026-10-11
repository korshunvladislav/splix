package core_http_request

import (
    "fmt"
    "time"

    core_errors "github.com/korshunvladislav/splix/internal/core/errors"
)

func ParseDate(value string) (time.Time, error) {
    t, err := time.Parse(time.DateOnly, value)
    if err != nil {
        return time.Time{}, fmt.Errorf(
            "value='%s' not a valid date: %v: %w",
            value,
            err,
            core_errors.ErrInvalidArgument,
        )
    }

    return t, nil
}

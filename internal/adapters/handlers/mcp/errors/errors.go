package errors

import (
	"errors"
	"fmt"

	domainErrors "github.com/gllm-dev/molinets/internal/core/domain/errors"
)

var codes = []struct {
	err  error
	code string
}{
	{domainErrors.ErrEmptyGrinderQuery, "INVALID_ARGUMENT"},
	{domainErrors.ErrGrinderNotFound, "GRINDER_NOT_FOUND"},
	{domainErrors.ErrAmbiguousGrinder, "AMBIGUOUS_GRINDER"},
	{domainErrors.ErrInvalidGrinder, "INVALID_ARGUMENT"},
	{domainErrors.ErrEmptySetting, "INVALID_ARGUMENT"},
	{domainErrors.ErrInvalidSetting, "INVALID_SETTING"},
	{domainErrors.ErrSettingOutOfRange, "SETTING_OUT_OF_RANGE"},
	{domainErrors.ErrMethodNotFound, "METHOD_NOT_FOUND"},
	{domainErrors.ErrMethodUnreachable, "METHOD_UNREACHABLE"},
	{domainErrors.ErrInvalidMicrons, "INVALID_ARGUMENT"},
}

func Code(err error) string {
	for _, c := range codes {
		if errors.Is(err, c.err) {
			return c.code
		}
	}
	return "INTERNAL"
}

func ToolError(err error) error {
	code := Code(err)
	if code == "INTERNAL" {
		return errors.New("INTERNAL: unexpected error")
	}
	return fmt.Errorf("%s: %w", code, err)
}

package errors

import "errors"

var (
	ErrEmptyGrinderQuery = errors.New("empty grinder name")
	ErrGrinderNotFound   = errors.New("grinder not found")
	ErrAmbiguousGrinder  = errors.New("grinder name is ambiguous")
	ErrInvalidGrinder    = errors.New("invalid grinder definition")
)

var (
	ErrEmptySetting      = errors.New("empty setting")
	ErrInvalidSetting    = errors.New("invalid setting")
	ErrSettingOutOfRange = errors.New("setting outside the usable range")
)

var (
	ErrMethodNotFound    = errors.New("brew method not found")
	ErrMethodUnreachable = errors.New("grinder cannot reach this brew method")
)

var (
	ErrInvalidMicrons = errors.New("invalid particle size")
)

package format

import (
	"errors"
	"strings"
)

// for unit testing broken format results
type BrokenFormat struct {
}

func (helper BrokenFormat) Default() string {
	return ""
}

func (helper BrokenFormat) Marshal(value string) ([]byte, error) {
	if strings.ToLower(value) == "common" {
		return []byte{}, ErrCommon
	}
	if strings.ToLower(value) == "marshal" {
		return []byte{}, ErrMarshal
	}
	if strings.ToLower(value) == "error" {
		return []byte{}, ErrBasic
	}
	return []byte(value), nil
}

func (helper BrokenFormat) Unmarshal(value []byte) (string, error) {
	if strings.ToLower(string(value)) == "common" {
		return "", ErrCommon
	}
	if strings.ToLower(string(value)) == "unmarshal" {
		return "", ErrUnmarshal
	}
	if strings.ToLower(string(value)) == "error" {
		return "", ErrBasic
	}
	return string(value), nil
}

func Broken() Format[string] {
	return BrokenFormat{}
}

var (
	ErrUnmarshal = errors.New("format unmarshal error")
	ErrMarshal   = errors.New("format marshal error")
	ErrCommon    = errors.New("format common error")
	ErrBasic     = errors.New("format basic error")
)

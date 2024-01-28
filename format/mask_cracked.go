package format

import (
	"errors"
	"strings"
)

// for unit testing broken format masking
type CrackedMask struct {
}

func (helper CrackedMask) Mask(value []byte) ([]byte, error) {
	if strings.ToLower(string(value)) == "mask" {
		return []byte{}, ErrMask
	}
	if strings.ToLower(string(value)) == "error" {
		return []byte{}, ErrBasic
	}
	return reverse(value), nil
}

func (helper CrackedMask) Unmask(value []byte) ([]byte, error) {
	reversed := reverse(value)
	if strings.ToLower(string(reversed)) == "unmask" {
		return []byte{}, ErrUnmask
	}
	if strings.ToLower(string(reversed)) == "error" {
		return []byte{}, ErrBasic
	}
	return reversed, nil
}

func Cracked() Mask {
	return CrackedMask{}
}

var (
	ErrMask   = errors.New("error masking")
	ErrUnmask = errors.New("error unmasking")
)

func reverse(s []byte) []byte {
	reversed := make([]byte, len(s))
	slen := len(s)
	for i, val := range s {
		reversed[slen-i-1] = val
	}
	return reversed
}

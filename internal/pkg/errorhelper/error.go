package errorhelper

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	BadRequest = 400
	NotFound   = 404
	Internal   = 500
)

const (
	CacheMissError = "cache miss"
)

type Error struct {
	StatusCode int

	Err error
}

func (r *Error) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}

func WrapWithCode(err error, code int, msg string) *Error {
	return &Error{
		StatusCode: code,
		Err:        errors.Wrap(err, msg),
	}
}

func NewWithCode(code int, msg string) *Error {
	return &Error{
		StatusCode: code,
		Err:        errors.New(msg),
	}
}

func NewCacheMissError() error {
	return errors.New(CacheMissError)
}

func IfCacheMiss(err error) bool {
	return err.Error() == CacheMissError
}

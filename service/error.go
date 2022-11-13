package service

import "fmt"

type userError struct {
	message string
}

func (e *userError) Error() string {
	return e.message
}

func isUserError(err error) bool {
	_, ok := err.(*userError)
	return ok
}

func newUserError(msg string, args ...any) *userError {
	return &userError{message: fmt.Sprintf(msg, args...)}
}

package server

import (
	"errors"
)

var ErrBusy = errors.New("Server busy")
var ErrAborted = errors.New("Aborted")
var ErrNotImplemented = errors.New("Not implemented")
var ErrNotFound = errors.New("Not found")
var ErrCSRF = errors.New("CSRF token doesn't match")
var ErrBadRequest = errors.New("Bad request")
var ErrPermissionDenied = errors.New("Permission denied")

var ErrAlreadyLoggedIn = errors.New("Already logged in")
var ErrUserNotFound = errors.New("User not found")
var ErrPasswordIncorrect = errors.New("Password incorrect")
var ErrDuplicateUserName = errors.New("Duplicate user name")

var ErrNotLoggedIn = errors.New("Not logged in")

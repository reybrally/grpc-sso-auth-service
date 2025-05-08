package storage

import "errors"

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrAppNotFound = errors.New("app not found")

package customerror

import "errors"

const (
	GORMDuplicatedKeyNotAllowed = "duplicated key not allowed"
)

var (
	ErrorEmailAlreadyExist = errors.New("email already registered")
)

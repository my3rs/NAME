package dict

import "errors"

type Dict string

// TokenType represents different types of JWT tokens
type TokenType int

const (
	TypeAccessToken TokenType = iota
	TypeRefreshToken
)

const (
	StatusSuccess Dict = "success"
	StatusFailed  Dict = "failed"
)

var (
	ErrInvalidParameters    = errors.New("invalid parameters")
	ErrWrongContentType     = errors.New("wrong content type")
	ErrContentAlreadyExists = errors.New("content already exists")
	ErrContentNotExists     = errors.New("content doesn't exists")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotExists        = errors.New("user doesn't exists")
	ErrWrongPassword        = errors.New("wrong password")
	ErrEmptyContent         = errors.New("empty content")
	ErrNotFound             = errors.New("item not found")
)

func (d Dict) String() string {
	return string(d)
}

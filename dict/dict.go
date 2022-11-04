package dict

import "errors"

type Dict string

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
)

func (d Dict) String() string {
	return string(d)
}

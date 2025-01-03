package customerror

import "fmt"

// 定义错误类型
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeDatabase     ErrorType = "DATABASE"
	ErrorTypeUnexpected   ErrorType = "UNEXPECTED"
	ErrorTypeJWTExpired   ErrorType = "JWT_EXPIRED"
	ErrorTypeJWTInvalid   ErrorType = "JWT_INVALID"
	ErrorTypeJWTMalformed ErrorType = "JWT_MALFORMED"
)

// CustomError 自定义错误结构
type CustomError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// 错误构造函数
func NewValidationError(message string) error {
	return &CustomError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

func NewNotFoundError(message string) error {
	return &CustomError{
		Type:    ErrorTypeNotFound,
		Message: message,
	}
}

func NewDatabaseError(message string, err error) error {
	return &CustomError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Err:     err,
	}
}

func NewJWTExpiredError() error {
	return &CustomError{
		Type:    ErrorTypeJWTExpired,
		Message: "token has expired",
	}
}

func NewJWTInvalidError() error {
	return &CustomError{
		Type:    ErrorTypeJWTInvalid,
		Message: "token is invalid",
	}
}

func NewJWTMalformedError() error {
	return &CustomError{
		Type:    ErrorTypeJWTMalformed,
		Message: "token is malformed",
	}
}

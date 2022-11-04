package domain

type Error struct {
	CodeHTTP int    `json:"-"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(httpCode int, code, message string) error {
	return &Error{
		CodeHTTP: httpCode,
		Code:     code,
		Message:  message,
	}
}

package serverr

const (
	AuthError = "AUTH_ERROR"
)

type Error struct {
	Description string `json:"-"`
	ErrType     string `json:"type"`
	HttpStatus  int    `json:"-"`
}

func (e *Error) Error() string {
	return e.ErrType + ": " + e.Description
}

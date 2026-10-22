package apperror

type ClientError struct {
	Status  int
	Code    string
	Message string
}

func (e *ClientError) Error() string {
	return e.Code + ": " + e.Message
}

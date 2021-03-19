package handlers

type AuthError struct {
	reason string
}

func (err *AuthError) Error() string {
	return err.reason
}

type InternalError struct {
	reason string
}

func (err *InternalError) Error() string {
	return err.reason
}

type NotFoundError struct {
}

func (err *NotFoundError) Error() string {
	return "Not found"
}
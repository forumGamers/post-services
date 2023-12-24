package errors

const (
	Forbidden    = "Forbidden"
	Unauthorized = "Unauthorized"
	InvalidToken = "Invalid Token"
)

type forbiddenError struct {
	msg        string
	StatusCode int
}

type internalServerError struct {
	msg        string
	StatusCode int
}

type unauthorizedError struct {
	msg        string
	StatusCode int
}

type dataNotFoundError struct {
	msg        string
	StatusCode int
}

type invalidObjectId struct{}
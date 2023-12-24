package errors

func (err *forbiddenError) Error() string {
	return err.msg
}

func (err *internalServerError) Error() string {
	return err.msg
}

func (err *unauthorizedError) Error() string {
	return err.msg
}

func (err *dataNotFoundError) Error() string {
	return err.msg
}

func (err *conflictError) Error() string {
	return err.msg
}

func (err *invalidObjectId) Error() string {
	return "Invalid ObjectId"
}

func (err *badGatewayError) Error() string {
	return err.msg
}

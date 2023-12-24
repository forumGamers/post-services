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

func (err *invalidObjectId) Error() string {
	return "Invalid ObjectId"
}

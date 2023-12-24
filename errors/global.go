package errors

func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func NewError(msg string, code int) error {
	switch code {
	case 403:
		return &forbiddenError{msg, code}
	case 401:
		return &unauthorizedError{msg, code}
	case 404:
		return &dataNotFoundError{msg, code}
	default:
		return &internalServerError{msg, 500}
	}
}

func GetErrorMsg(err error) (string, int) {
	switch e := err.(type) {
	case *forbiddenError:
		return e.msg, e.StatusCode
	case *unauthorizedError:
		return e.msg, e.StatusCode
	case *dataNotFoundError:
		return e.msg, e.StatusCode
	case *invalidObjectId:
		return e.Error(), 400
	default:
		return "Internal Server Error", 500
	}
}

func NewInvalidObjectIdError() error {
	return &invalidObjectId{}
}

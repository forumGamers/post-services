package helper

import "fmt"

var (
	ErrInvalidObjectId = fmt.Errorf("Invalid ObjectID")
	Forbidden = fmt.Errorf("Forbidden")
	InvalidToken = fmt.Errorf("Invalid Token")
	NotFount = fmt.Errorf("Data not found")
	InternalServer = fmt.Errorf("Internal Server Error")
	InvalidChiper = fmt.Errorf("invalid ciphertext block size")
	BadGateway = fmt.Errorf("Bad Gateway")
	Conflict = fmt.Errorf("Conflicts")
)

func PanicIfError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
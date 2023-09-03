package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	h "github.com/post-services/helper"
)

type WebResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type HttpError struct {
	Error   error
	Code    int
	Message string
}

type InputError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getStatusMessage(status int) string {
	statusMessages := map[int]string{
		100: "Continue",
		101: "Switching Protocols",
		102: "Processing",
		103: "Early Hints",
		200: "OK",
		201: "Created",
		202: "Accepted",
		203: "Non Authoritative Info",
		204: "No Content",
		205: "Reset Content",
		206: "Partial Content",
		207: "Multi Status",
		208: "Already Reported",
		226: "IM Used",
		300: "Multiple Choices",
		301: "Moved Permanently",
		302: "Found",
		303: "See Other",
		304: "Not Modified",
		305: "Use Proxy",
		307: "Temporary Redirect",
		308: "Permanent Redirect",
		400: "Bad Request",
		401: "Unauthorized",
		402: "Payment Required",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		406: "Not Acceptable",
		407: "Proxy Auth Required",
		408: "Request Timeout",
		409: "Conflict",
		410: "Gone",
		411: "Length Required",
		412: "Precondition Failed",
		413: "Request Entity Too Large",
		414: "Request URI Too Long",
		415: "Unsupported Media Type",
		416: "Request Range Not Satisfiable",
		417: "Expectation Failed",
		418: "Teapot",
		421: "Misdirected Request",
		422: "Unprocessable Entity",
		423: "Locked",
		424: "Failed Dependency",
		425: "Too Early",
		426: "Upgrade Required",
		429: "Too Many Requests",
		431: "Request Header Fields Too Large",
		451: "Unavailable For Legal Reasons",
		500: "Internal Server Error",
		501: "Not Implemented",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Timeout",
		505: "HTTP Version Not Supported",
		506: "Variant Also Negotiates",
		507: "Insufficient Storage",
		508: "Loop Detected",
		510: "Not Extended",
		511: "Network Authentication Required",
	}

	return statusMessages[status]
}

func getErrorMsg(err error) (string, int) {
	switch err.Error() {
	case "Data not found":
		return err.Error(), 404
	case h.ErrInvalidObjectId.Error():
		return h.ErrInvalidObjectId.Error(), 400
	case h.Forbidden.Error():
		return h.Forbidden.Error(), 403
	case h.Conflict.Error():
		return h.Conflict.Error(), 409
	case h.BadGateway.Error():
		return h.BadGateway.Error(), 502
	case h.InvalidChiper.Error():
		return h.InvalidChiper.Error(), 500
	case h.InvalidToken.Error():
		return h.InvalidToken.Error(), 401
	case h.AccessDenied.Error():
		return h.AccessDenied.Error(), 401
	default:
		return "Internal Server Error", 500
	}
}

func WriteResponse(c *gin.Context, response WebResponse) {
	response.Status = getStatusMessage(response.Code)
	c.JSON(response.Code, response)
}

func AbortHttp(c *gin.Context, err error) {
	msg, code := getErrorMsg(err)

	c.AbortWithStatusJSON(code, WebResponse{
		Status:  getStatusMessage(code),
		Code:    code,
		Message: msg,
	})
}

func CustomMsgAbortHttp(c *gin.Context, message string, code int) {
	c.AbortWithStatusJSON(code, WebResponse{
		Status:  getStatusMessage(code),
		Code:    code,
		Message: message,
	})
}

func MsgTag(f validator.FieldError) string {
	switch f.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "oneof":
		switch true {
		case f.Field() == "Privacy":
			return "input must be one of Public, FriendOnly, Private"
		default:
			return "value must be one of the enum"
		}
	case "required_without":
		switch true {
		case f.Field() == "Text":
			return "text is required if file is empty"
		case f.Field() == "File":
			return "file is required if text is empty"
		default:
			return "this field is required if the other one is empty"
		}
	default:
		return f.Error()
	}
}

func HttpValidationErr(c *gin.Context, err error) {
	errMap := make(map[string]string)
	for _, val := range err.(validator.ValidationErrors) {
		errMap[val.Field()] = MsgTag(val)
	}

	c.AbortWithStatusJSON(400, WebResponse{
		Status: getStatusMessage(400),
		Code:   400,
		Data:   errMap,
	})
}

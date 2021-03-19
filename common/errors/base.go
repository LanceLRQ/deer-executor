package errors

import (
	"fmt"
)

type ApplicationError struct {
	error
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewApplicationError(code int, args ...interface{}) error {
	rawmsg, ok := ApplicationErrorMap[code]
	if ok {
		return &ApplicationError{Code: code, Msg: fmt.Sprintf(rawmsg, args)}
	}
	return &ApplicationError{Code: code, Msg: "Unknown error"}
}

func (e *ApplicationError) Error() string {
	//outmsg := utils.ObjectToJSONString(e)
	//if outmsg != "" {
	//    return outmsg
	//}
	return e.Msg
}

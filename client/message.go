package client

import (
	"fmt"
	"github.com/LanceLRQ/deer-common/utils"
)

// CliCommonMessage is a client comment message
type CliCommonMessage struct {
	// 是否错误
	Error bool `json:"error"`
	// 消息
	Message string `json:"message"`
	// 结果信息
	Data interface{} `json:"data"`
}

// Print print log to stdout
func (ccm CliCommonMessage) Print(formated bool) {
	fmt.Println(ccm.ToJSON(formated))
}

// ToJSON convert log record to json
func (ccm CliCommonMessage) ToJSON(formated bool) string {
	if formated {
		return utils.ObjectToJSONStringFormatted(ccm)
	}
	return utils.ObjectToJSONString(ccm)
}

// NewCliCommonMessage to create a client comment message
func NewCliCommonMessage(error bool, message string, data interface{}) CliCommonMessage {
	return CliCommonMessage{
		Error:   error,
		Message: message,
		Data:    data,
	}
}

// NewClientSuccessMessage to create a client success message
func NewClientSuccessMessage(data interface{}) CliCommonMessage {
	return CliCommonMessage{
		Error:   false,
		Message: "",
		Data:    data,
	}
}

// NewClientSuccessMessageText to create a new client message with text
func NewClientSuccessMessageText(message string) CliCommonMessage {
	return CliCommonMessage{
		Error:   false,
		Message: message,
		Data:    nil,
	}
}

// NewClientErrorMessage to create a client error message
func NewClientErrorMessage(err error, data interface{}) CliCommonMessage {
	return CliCommonMessage{
		Error:   true,
		Message: err.Error(),
		Data:    data,
	}
}

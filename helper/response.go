package helper

import "strings"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	// To return dynamic data
	Error interface{} `json:"error"`
	Data  interface{} `json:"data"`
}

type EmptyObject struct {
}

func SendSuccessResponse(success bool, message string, data interface{}) Response {
	res := Response{
		Success: true,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	return res
}

func SendErrorResponse(success bool, message string, err string, data interface{}) Response {
	splitedError := strings.Split(err, "\n")
	res := Response{
		Success: success,
		Message: message,
		Error:   splitedError,
		Data:    data,
	}
	return res
}

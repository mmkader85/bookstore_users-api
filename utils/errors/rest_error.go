package errors

import (
	"log"
	"net/http"
	"runtime"
)

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func BadRequestErr(message string) *RestErr {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("Error in %s\n [%s:%d] \n%s", runtime.FuncForPC(pc).Name(), fn, line, message)

	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func NotFoundErr(message string) *RestErr {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("Error in %s\n [%s:%d] \n%s", runtime.FuncForPC(pc).Name(), fn, line, message)

	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func InternalServerErr(message string) *RestErr {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("Error in %s\n [%s:%d] \n%s", runtime.FuncForPC(pc).Name(), fn, line, message)

	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}

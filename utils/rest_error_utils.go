package utils

import (
	"log"
	"net/http"
	"runtime"
)

var RestErrUtils restErrUtilsInterface = &restErrUtilsStruct{}

type restErrUtilsStruct struct{}

type restErrUtilsInterface interface {
	NotFoundErr(message string) *RestErr
	BadRequestErr(message string) *RestErr
	InternalServerErr(message string) *RestErr
}

type RestErr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func (restErrUtilsStruct) NotFoundErr(message string) *RestErr {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("Error in %s\n [%s:%d] \n%s", runtime.FuncForPC(pc).Name(), fn, line, message)

	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

func (restErrUtilsStruct) BadRequestErr(message string) *RestErr {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("Error in %s\n [%s:%d] \n%s", runtime.FuncForPC(pc).Name(), fn, line, message)

	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

func (restErrUtilsStruct) InternalServerErr(message string) *RestErr {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("Error in %s\n [%s:%d] \n%s", runtime.FuncForPC(pc).Name(), fn, line, message)

	return &RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
}

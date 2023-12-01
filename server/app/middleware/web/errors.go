package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorInterface struct {
	ctx          *gin.Context
	statusCode   int
	errorData    interface{}
	errorMessage string
	errorCode    string
}

func (e *ErrorInterface) GetStatusCode() int {
	return e.statusCode
}

func (e *ErrorInterface) GetErrorData() interface{} {
	return e.errorData
}

func (e *ErrorInterface) GetErrorMessage() string {
	return e.errorMessage
}

func (e *ErrorInterface) GetErrorCode() string {
	return e.errorCode
}

func NewHTTPBadRequestError(errorMessage string, errorData interface{}) *ErrorInterface {
	return NewError(http.StatusBadRequest, errorMessage, "BAD_REQUEST", errorData)
}

func NewHTTPNotFoundError(errorMessage string, errorData interface{}) *ErrorInterface {
	return NewError(http.StatusNotFound, errorMessage, "RESOURCE_NOT_FOUND", errorData)
}

func NewError(statusCode int, errorMessage string, errorCode string, errorData interface{}) *ErrorInterface {
	return &ErrorInterface{statusCode: statusCode, errorMessage: errorMessage, errorData: errorData, errorCode: errorCode}
}

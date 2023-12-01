/*
Date: 2023-05-11, Thu, 12:51
Version: v1.1
Copyright 2023, Akshay Singh Kanawat
Developer: Akshay Singh Kanawat
*/

package middleware

import (
	"github.com/gin-gonic/gin"
	"pgsync/server/app/middleware/web"
)

type HandlerFunc func(context *gin.Context) (*web.JSONResponse, *web.ErrorInterface)

func ServeEndpoint(nextHandler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Validate the request body
		data, responseErr := nextHandler(ctx)
		ctx.JSON(buildResponseCode(data, responseErr), buildResponseBody(data, responseErr))
	}
}

func buildResponseCode(data *web.JSONResponse, responseErr *web.ErrorInterface) int {
	if responseErr == nil {
		return data.GetStatusCode()
	}
	return responseErr.GetStatusCode()
}

func buildResponseBody(data *web.JSONResponse, responseErr *web.ErrorInterface) interface{} {
	if responseErr == nil {
		return data.GetResponseBody()
	}
	return buildErrorResponse(responseErr)
}

func buildErrorResponse(responseErr *web.ErrorInterface) interface{} {
	return map[string]interface{}{
		"errorCode":    responseErr.GetErrorCode(),
		"errorMessage": responseErr.GetErrorMessage(),
		"errorData":    responseErr.GetErrorData(),
	}
}

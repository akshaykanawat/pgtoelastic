package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type JSONResponse struct {
	ctx          *gin.Context
	statusCode   int
	responseBody map[string]interface{}
}

func (j *JSONResponse) GetStatusCode() int {
	return j.statusCode
}

func (j *JSONResponse) GetResponseBody() map[string]interface{} {
	return j.responseBody
}

func NewHTTPSuccessResponse(ctx *gin.Context, responseBody map[string]interface{}) *JSONResponse {
	return NewJSONResponse(ctx, http.StatusOK, responseBody)
}

func NewHTTPCreatedResponse(ctx *gin.Context, responseBody map[string]interface{}) *JSONResponse {
	return NewJSONResponse(ctx, http.StatusCreated, responseBody)
}

func NewJSONResponse(ctx *gin.Context, statusCode int, responseBody map[string]interface{}) *JSONResponse {
	return &JSONResponse{
		ctx:          ctx,
		statusCode:   statusCode,
		responseBody: responseBody,
	}
}

func (j *JSONResponse) JSON() {
	j.ctx.JSON(j.statusCode, j.responseBody)
}

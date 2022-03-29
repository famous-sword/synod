package render

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	aborted bool   `json:"-"`
}

func Fail() *Response {
	return &Response{}
}

func Success() *Response {
	return &Response{}
}

func (r *Response) WithStatus(status int) *Response {
	r.Status = status

	return r
}

func (r *Response) WithMessage(message string) *Response {
	r.Message = message

	return r
}

func (r *Response) With(data any) *Response {
	r.Data = data
	return r
}

func (r *Response) WithError(err error) *Response {
	r.Status = http.StatusBadRequest
	r.Message = err.Error()
	r.aborted = true

	return r
}

func (r *Response) To(ctx *gin.Context) {
	if r.aborted {
		ctx.AbortWithStatusJSON(r.Status, r)
	} else {
		ctx.JSON(r.Status, r)
	}
}

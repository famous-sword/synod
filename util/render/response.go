package render

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response a standard response of Synod
// and it provides multiple factory for use
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	aborted bool   `json:"-"`
}

// Fail used to bad request
func Fail() *Response {
	return &Response{}
}

// Success used to success request
func Success() *Response {
	return &Response{}
}

// OfError respond from an error quickly
func OfError(err error) *Response {
	r := Fail().WithError(err)

	return r
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

// To write response to context
func (r *Response) To(ctx *gin.Context) {
	if r.aborted {
		ctx.AbortWithStatusJSON(r.Status, r)
	} else {
		ctx.JSON(r.Status, r)
	}
}

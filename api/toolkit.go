package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

const (
	KeyHash   = "digest"
	KeySize   = "content-length"
	KeyOffset = "range"
)

func getHash(ctx *gin.Context) string {
	digest := ctx.GetHeader(KeyHash)

	if digest[:8] != "SHA-256=" {
		return ""
	}

	return digest[8:]
}

func getSize(ctx *gin.Context) int64 {
	size, _ := strconv.ParseInt(ctx.GetHeader(KeySize), 10, 63)

	return size
}

func getOffset(ctx *gin.Context) int64 {
	r := ctx.GetHeader(KeyOffset)

	if len(r) < 7 {
		return 0
	}

	if r[:6] != "bytes=" {
		return 0
	}

	pos := strings.Split(r[6:], "-")

	offset, _ := strconv.ParseInt(pos[0], 0, 64)

	return offset
}

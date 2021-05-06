package simplerr

import (
	"net/http"
)

// grpc codes from google.golang.org/grpc@v1.34.0/codes/codes.go
const (
	grpcInvalidArgument  int = 3
	grpcNotFound         int = 5
	grpcAlreadyExists    int = 6
	grpcPermissionDenied int = 7
	grpcInternal         int = 13
	grpcUnauthenticated  int = 16
)

type ErrCode interface {
	HTTP() int
	GRPC() int
	Int() int
}

type code struct{ http, grpc, code int }

func (c *code) HTTP() int { return c.http }
func (c *code) GRPC() int { return c.grpc }
func (c *code) Int() int  { return c.code }

func InvalidArgumentCode(c int) ErrCode {
	return &code{
		http: http.StatusBadRequest,
		grpc: grpcInvalidArgument,
		code: c,
	}
}

func NotFoundCode(c int) ErrCode {
	return &code{
		http: http.StatusNotFound,
		grpc: grpcNotFound,
		code: c,
	}
}

func AlreadyExistsCode(c int) ErrCode {
	return &code{
		http: http.StatusConflict,
		grpc: grpcAlreadyExists,
		code: c,
	}
}

func UnauthorizedCode(c int) ErrCode {
	return &code{
		http: http.StatusUnauthorized,
		grpc: grpcUnauthenticated,
		code: c,
	}
}

func ForbiddenCode(c int) ErrCode {
	return &code{
		http: http.StatusForbidden,
		grpc: grpcPermissionDenied,
		code: c,
	}
}

func InternalCode(c int) ErrCode {
	return &code{
		http: http.StatusInternalServerError,
		grpc: grpcInternal,
		code: c,
	}
}

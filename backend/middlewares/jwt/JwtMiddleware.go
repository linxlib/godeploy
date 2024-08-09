package jwt

import (
	"github.com/linxlib/fw"
)

type JwtMiddleware struct {
	*fw.MiddlewareCtl
}

func (j *JwtMiddleware) CloneAsMethod() fw.IMiddlewareMethod {
	//TODO implement me
	panic("implement me")
}

func (j *JwtMiddleware) HandlerMethod(next fw.HandlerFunc) fw.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (j *JwtMiddleware) CloneAsCtl() fw.IMiddlewareCtl {
	//TODO implement me
	panic("implement me")
}

func (j *JwtMiddleware) HandlerController(base string) []*fw.RouteItem {
	return []*fw.RouteItem{
		{
			Method: "POST",
			Path:   base + "/refresh_token",
			IsHide: false,
			H: func(c *fw.Context) {

			},
			Middleware: j,
		},
		{
			Method: "POST",
			Path:   base + "/refresh_token",
			IsHide: false,
			H: func(c *fw.Context) {

			},
			Middleware: j,
		},
	}
}

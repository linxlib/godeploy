package middlewares

import (
	"github.com/linxlib/fw"
	"os"
	"path/filepath"
)

type StaticMiddleware struct {
	*fw.MiddlewareGlobal
	dir string
}

func NewStaticMiddleware(dir string) *StaticMiddleware {
	return &StaticMiddleware{
		MiddlewareGlobal: fw.NewMiddlewareGlobal("StaticMiddleware"),
		dir:              dir,
	}
}

func (s *StaticMiddleware) CloneAsMethod() fw.IMiddlewareMethod {
	return s.CloneAsCtl()
}

func (s *StaticMiddleware) HandlerMethod(next fw.HandlerFunc) fw.HandlerFunc {
	return next
}

func (s *StaticMiddleware) CloneAsCtl() fw.IMiddlewareCtl {
	return NewStaticMiddleware(s.dir)
}

func (s *StaticMiddleware) HandlerController(base string) []*fw.RouteItem {
	return []*fw.RouteItem{&fw.RouteItem{
		Method: "GET",
		Path:   base + "/static/{path}",
		IsHide: false,
		H: func(ctx *fw.Context) {
			//var isThumbnail = ctx.GetFastContext().QueryArgs().GetUintOrZero("isThumbnail") == 1
			//var width = ctx.GetFastContext().QueryArgs().GetUintOrZero("width")
			//var height = ctx.GetFastContext().QueryArgs().GetUintOrZero("height")
			relativePath := ctx.Param("path")
			realPath := filepath.Join(s.dir, relativePath)
			if _, err := os.Stat(realPath); os.IsNotExist(err) {
				ctx.Status(404)
				return
			}
			ctx.File(realPath)
		},
		Middleware: s,
	}}
}

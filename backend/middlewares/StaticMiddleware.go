package middlewares

import (
	"fmt"
	"github.com/linxlib/fw"
	"os"
	"path/filepath"
)

type StaticMiddleware struct {
	*fw.MiddlewareGlobal
	options *StaticOptions
}

type StaticOptions struct {
	List []StaticOption `json:"list"`
}

type StaticOption struct {
	Path  string `yaml:"path"`
	Route string `yaml:"route"`
}

func NewStaticMiddleware() *StaticMiddleware {
	return &StaticMiddleware{
		MiddlewareGlobal: fw.NewMiddlewareGlobal("Static"),
	}
}
func (s *StaticMiddleware) DoInitOnce() {
	s.options = new(StaticOptions)
	s.LoadConfig("static", s.options)
}

func (s *StaticMiddleware) Router(ctx *fw.MiddlewareContext) []*fw.RouteItem {
	var results = make([]*fw.RouteItem, 0)
	for _, option := range s.options.List {
		result := &fw.RouteItem{
			Method: "GET",
			Path:   option.Route,
			IsHide: false,
			H: func(ctx *fw.Context) {
				relativePath := ctx.Param("path")
				realPath := filepath.Join(option.Path, relativePath)
				fmt.Printf("%+v\n", s.options)
				if _, err := os.Stat(realPath); os.IsNotExist(err) {
					//ctx.Status(404)
					ctx.String(404, realPath)
					return
				}
				ctx.File(realPath)
			},
			Middleware: s,
		}
		results = append(results, result)
	}
	return results

}

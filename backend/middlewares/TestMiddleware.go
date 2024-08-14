package middlewares

import "github.com/linxlib/fw"

type TestMiddleware struct {
	*fw.MiddlewareCtl
}

func NewTestMiddleware() *TestMiddleware {
	return &TestMiddleware{MiddlewareCtl: fw.NewMiddlewareCtl("Test", "Test")}
}

func (t *TestMiddleware) Router(ctx *fw.MiddlewareContext) []*fw.RouteItem {
	return []*fw.RouteItem{
		{
			Method: "GET",
			Path:   "/test",
			IsHide: false,
			H: func(context *fw.Context) {
				context.WriteString("Hello World")
			},
			Middleware: t,
		},
	}

}

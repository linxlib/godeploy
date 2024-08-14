package weblog

import (
	"embed"
	"github.com/linxlib/fw"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

//go:embed html/*
var content embed.FS

func NewWebLogMiddleware() *WebLogMiddleware {

	return &WebLogMiddleware{
		MiddlewareGlobal: fw.NewMiddlewareGlobal("weblog"),
	}
}

type WebLogMiddleware struct {
	*fw.MiddlewareGlobal
	service *Service
	Logger  *logrus.Logger `inject:""`
}

func (w *WebLogMiddleware) DoInitOnce() {
	svc, _ := New(logrus.New(), &Config{
		Root:              "log/",
		Bytes:             50000,
		Lines:             1000,
		MaxLineSize:       180,
		ListCache:         2,
		Poll:              false,
		Trace:             false,
		ClientBufferSize:  256,
		WSReadBufferSize:  1024,
		WSWriteBufferSize: 1024,
	})
	w.service = svc
	go w.service.Run()
}

func (w *WebLogMiddleware) Router(ctx *fw.MiddlewareContext) []*fw.RouteItem {
	return []*fw.RouteItem{&fw.RouteItem{
		Method:     "GET",
		Path:       "/ws/tail",
		Middleware: w,
		H:          w.service.ServeHTTP,
	},
		&fw.RouteItem{
			Method:     "GET",
			Path:       "/tail/",
			Middleware: w,
			H: func(context *fw.Context) {
				fasthttp.ServeFS(context.GetFastContext(), content, "/html/index.html")
			},
		},
		&fw.RouteItem{
			Method:     "GET",
			Path:       "/tail/{filepath:*}",
			Middleware: w,
			H: func(context *fw.Context) {
				path := context.GetFastContext().UserValue("filepath").(string)
				fasthttp.ServeFS(context.GetFastContext(), content, "/html/"+path)
			},
		},
	}
}

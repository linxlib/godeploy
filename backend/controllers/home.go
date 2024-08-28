package controllers

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/linxlib/conv"
	"github.com/linxlib/fw"
)

const html = `<!DOCTYPE html>
<html lang='en'>
<head>
    <meta charset='UTF-8'>
    <meta name='viewport' content='width=device-width, initial-scale=1.0'>
    <title>部署服务 {{.version}}</title>
</head>
<body>
<h1>{{.version}}</h1>
    <script>
        //window.location.href = "/deploy"
    </script>
</body>
</html>`

// HomeController
// @Controller
// @Route /
// @Test
type HomeController struct{}

// Index
// @GET /
func (h *HomeController) Index(ctx *fw.Context) {
	panic("implement me")
	ctx.HTMLPure(200, html, map[string]interface{}{
		"version": "1.0",
	})
}

// CheckFileExists
// @GET /check_file
func (h *HomeController) CheckFileExists(ctx *fw.Context) {
	a := fsutil.FileExist("/root/下载.png")
	ctx.String(200, conv.String(a))
}

// Websocket
// @Websocket
// @ANY /ws
func (h *HomeController) Websocket(bs []byte) {

}

// Home2Controller
// @Controller
// @Route /home2
// @Test
type Home2Controller struct{}

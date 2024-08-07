package controllers

import "github.com/linxlib/fw"

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
type HomeController struct{}

// Index
// @GET /
func (h *HomeController) Index(ctx *fw.Context) {
	panic("implement me")
	ctx.HTMLPure(200, html, map[string]interface{}{
		"version": "1.0",
	})
}

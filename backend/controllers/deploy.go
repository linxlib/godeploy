package controllers

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/controllers/models"
	"github.com/linxlib/godeploy/services/deploy"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"time"
)

// DeployController
// @Route /deploy
// @Controller
type DeployController struct {
}

// DeployFormRequest
// @Multipart
type DeployFormRequest struct {
	File      multipart.FileHeader `multipart:"file"` //需要部署的压缩包
	Hash      string               `multipart:"hash"`
	ServiceId int                  `multipart:"serviceId" validate:"required"` //服务ID
}

// Deploy
// @POST /
func (c *DeployController) Deploy(ctx *fw.Context, req *DeployFormRequest, db *gorm.DB) {
	service := new(models.Service)
	if err := db.First(&service, req.ServiceId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, map[string]interface{}{
				"code":    404,
				"message": "未找到服务",
				"data":    nil,
			})
		} else {
			ctx.JSON(500, map[string]interface{}{
				"code":    500,
				"message": err.Error(),
				"data":    nil,
			})
		}
		return
	}
	id := deploy.CreateNewDeploy(&req.File, service, db)
	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"data": id,
	})

}

// @Query
type SSEReq struct {
	ID int `query:"id" validate:"required"` //deploy接口返回的id
}

// SSE deploy processing
// @GET /sse
func (c *DeployController) SSE(ctx *fw.Context, req *SSEReq) {
	ctx.GetFastContext().SetContentType("text/event-stream")
	ctx.SetHeader("Cache-Control", "no-cache")
	ctx.SetHeader("Connection", "keep-alive")
	//ctx.SetHeader("Access-Control-Allow-Origin", "*")
	ctx.SetHeader("Transfer-Encoding", "chunked")
	eventId := req.ID
	if eventId == 0 {
		ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "event not found",
		})
		return
	}
	timer := time.NewTicker(time.Second * 1)
	ch, done := deploy.StartDeploy(eventId)

	ctx.Stream(func(w *bufio.Writer) {
		for {
			select {
			case e := <-ch:
				log.Println("write data:", e.Message)
				_, err := fmt.Fprintf(w, "data: %s\n\n", e.Message)
				if err != nil {
					log.Println("write data:", err)
					continue
				}
				w.Flush()
			case <-timer.C:
				//log.Println("write time:")
				//_, err := fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.DateTime))
				//if err != nil {
				//	log.Println("write time:", err)
				//	continue
				//}
				//w.Flush()
			case <-done:
				fmt.Fprintf(w, "data: %s\n\n", "done")
				return
			}
		}

	})

}

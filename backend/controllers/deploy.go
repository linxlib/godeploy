package controllers

import (
	"context"
	"errors"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/controllers/models"
	"github.com/saracen/fastzip"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
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
	File      multipart.FileHeader `multipart:"file"`
	Hash      string               `multipart:"hash"`
	ServiceId int                  `multipart:"serviceId" validate:"required"`
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
	f, err := req.File.Open()
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	defer f.Close()
	//TODO file hash check

	tmpEx := "./tmp"
	os.Mkdir(tmpEx, 0777)
	os.RemoveAll(tmpEx)
	reader, err := fastzip.NewExtractorFromReader(f, req.File.Size, tmpEx)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	err = reader.Extract(context.Background())
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if !service.OverwriteFrom(tmpEx) {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "文件覆盖失败",
			"data":    nil,
		})
		return
	}
	if !service.Start() {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "服务启动失败",
			"data":    nil,
		})
		return
	}
	service.LastDeployTime = time.Now()
	db.Save(service)
	ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "服务部署成功",
		"data":    nil,
	})
	return

}

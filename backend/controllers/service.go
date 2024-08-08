package controllers

import (
	"errors"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/base"
	"github.com/linxlib/godeploy/controllers/models"
	"gorm.io/gorm"
)

func NewServiceController(db *gorm.DB) *ServiceController {
	a := &ServiceController{
		CrudController: base.NewCrudController[*models.Service](db),
	}
	return a
}

// ServiceController
// @Route /service
// @Controller
// @Session
type ServiceController struct {
	*base.CrudController[*models.Service]
}

// StatusRequest
// @Query
type StatusRequest struct {
	ServiceId uint `query:"service_id"`
}

// Status
// @GET /status
func (s *ServiceController) Status(ctx *fw.Context, req *StatusRequest, db *gorm.DB) {
	var service = new(models.Service)
	if err := db.First(&service, req.ServiceId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, map[string]interface{}{
				"code":    404,
				"message": err.Error(),
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
	ctx.JSON(200, base.Resp(200, "", service.Status()))
}

// ActionRequest
// @Body
type ActionRequest struct {
	ServiceId  uint               `json:"service_id"`
	ActionType *models.ActionType `json:"action_type"`
}

// Action
// @POST /action
func (s *ServiceController) Action(ctx *fw.Context, req *ActionRequest, db *gorm.DB) {
	var service = new(models.Service)
	if err := db.First(&service, req.ServiceId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, map[string]interface{}{
				"code":    404,
				"message": err.Error(),
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
	switch *req.ActionType {
	case models.Reboot:
		service.Stop()
		service.Start()
		ctx.JSON(200, base.Resp(200, "", service.Status()))
	case models.Start:
		service.Start()
		ctx.JSON(200, base.Resp(200, "", service.Status()))
	case models.Stop:
		service.Stop()
		ctx.JSON(200, base.Resp(200, "", service.Status()))
	default:
		ctx.JSON(200, base.Resp(200, "", service.Status()))
	}

}

// DirTree
// @GET /dir_tree
func (s *ServiceController) DirTree(ctx *fw.Context) {

}

// IISServiceList
// @GET /IISServiceList
func (s *ServiceController) IISServiceList() {

}

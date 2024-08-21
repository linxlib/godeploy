package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/base"
	"github.com/linxlib/godeploy/controllers/models"
	"github.com/linxlib/godeploy/pkgs/dir_tree"
	"github.com/linxlib/godeploy/pkgs/golang-iis/iis"
	"github.com/linxlib/godeploy/pkgs/golang-iis/iis/websites"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"strings"
)

func NewServiceController() *ServiceController {
	a := &ServiceController{
		SimpleCrudController: base.NewSimpleCrudController[uint, *models.Service](),
	}
	return a
}

// ServiceController
// @Route /service
// @Controller
// @Session
// @Test
type ServiceController struct {
	*base.SimpleCrudController[uint, *models.Service]
}

// StatusRequest
// @Query
type StatusRequest struct {
	ServiceId uint `query:"service_id"`
}

// Status
// @GET /status
func (s *ServiceController) Status(ctx *fw.Context, req *StatusRequest) {

	var service = new(models.Service)
	if err := s.DB.First(&service, req.ServiceId).Error; err != nil {
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
func (s *ServiceController) Action(ctx *fw.Context, req *ActionRequest) {
	var service = new(models.Service)
	if err := s.DB.First(&service, req.ServiceId).Error; err != nil {
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

		ctx.JSON(200, base.Resp(200, fmt.Sprint(service.Restart()), service.Status()))
	case models.Start:
		ctx.JSON(200, base.Resp(200, fmt.Sprint(service.Start()), service.Status()))
	case models.Stop:
		ctx.JSON(200, base.Resp(200, fmt.Sprint(service.Stop()), service.Status()))
	default:
		ctx.JSON(200, base.Resp(200, "", service.Status()))
	}

}

// DirTreeRequest
// @Query
type DirTreeRequest struct {
	Home string `query:"home" validate:"required"`
}

// DirTree
// @GET /dir_tree
func (s *ServiceController) DirTree(ctx *fw.Context, req *DirTreeRequest) {
	tree, err := dir_tree.NewDirList(req.Home)
	if err != nil {
		ctx.JSON(200, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	ctx.JSON(fasthttp.StatusOK, base.Data(tree))
}

// IISServiceList
// @GET /IISServiceList
func (s *ServiceController) IISServiceList(ctx *fw.Context) {
	client, err := iis.NewClient()
	if err != nil {
		ctx.JSON(200, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	type Website1 struct {
		*websites.Website
		Binding string `json:"binding"`
	}
	wss, err := client.Websites.GetAll()
	if err != nil {
		ctx.JSON(200, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	results := make([]*Website1, 0)
	for _, website := range wss {
		result := &Website1{Website: website}
		bindings, err := client.Websites.GetBindings(website.Name)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{
				"code":    500,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		for i, binding := range bindings {
			if i != 0 {
				result.Binding += ","
			}
			result.Binding += fmt.Sprintf("%s:%d", binding.IPAddress, binding.Port)
		}

		results = append(results, result)
	}

	ctx.JSON(fasthttp.StatusOK, base.ListData(results, len(results), 0).Resp())
}

// FindUnitRequest
// @Query
type FindUnitRequest struct {
	Name string `query:"name" validate:"required"` //unit name with or without .service suffix
}

type FindUnitResponse struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	JobPath string `json:"job_path"`
	Desc    string `json:"desc"`
}

// FindUnit
// @GET /find_unit
func (s *ServiceController) FindUnit(ctx *fw.Context, query *FindUnitRequest) {
	conn, _ := dbus.NewSystemdConnectionContext(context.Background())
	defer conn.Close()
	if !strings.HasSuffix(query.Name, ".service") {
		query.Name += ".service"
	}
	names, err := conn.ListUnitsByNamesContext(context.Background(), []string{query.Name})
	if err != nil {
		return
	}
	results := make([]*FindUnitResponse, 0)
	for _, name := range names {
		results = append(results, &FindUnitResponse{
			Name:    name.Name,
			Path:    string(name.Path),
			JobPath: string(name.JobPath),
			Desc:    name.Description,
		})
	}

	ctx.JSON(fasthttp.StatusOK, base.ListData(results, len(results), 0).Resp())

}

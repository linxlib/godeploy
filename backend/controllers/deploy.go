package controllers

import "mime/multipart"

// DeployController
// @Route /deploy
type DeployController struct {
}

// DeployFormRequest
// @Multipart
type DeployFormRequest struct {
	File      multipart.FileHeader `form:"file"`
	Hash      string               `form:"hash"`
	Size      int64                `form:"size"`
	serviceId string               `form:"serviceId"`
}

// Deploy
// @POST /
func (c *DeployController) Deploy(req *DeployFormRequest) {

}

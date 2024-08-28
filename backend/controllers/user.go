package controllers

import (
	"errors"
	"github.com/fasthttp/session/v2"
	"github.com/linxlib/conv"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/base"
	models2 "github.com/linxlib/godeploy/base/models"
	"github.com/linxlib/godeploy/controllers/models"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"time"
)

func NewUserController() *UserController {
	a := &UserController{
		SimpleCrudController: base.NewSimpleCrudController[uint, *models.User](),
	}
	return a
}

// UserController 用户
// @Controller
// @Route /user
// @Session
type UserController struct {
	*base.SimpleCrudController[uint, *models.User]
}

// UserApproveRequest
// @Path
type UserApproveRequest struct {
	ID uint `path:"userid"` //用户id
}

// Approve 用户审核
// @POST /approve/{userid}
func (u *UserController) Approve(ctx *fw.Context, req *UserApproveRequest) {
	var user models.User
	if err := u.DB.First(&user, req.ID).Error; err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if user.Enabled {
		ctx.JSON(200, base.Resp[int, any](200, "用户已经审核", nil))
		return
	}
	user.Enabled = true
	if err := u.DB.Save(&user).Error; err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	ctx.JSON(200, base.Resp(200, "审核成功", user))
}

// SignInRequest
// @Body
type SignInRequest struct {
	Username string `json:"username" validate:"required"` //用户名
	Password string `json:"password" validate:"required"` //密码
}

// LoginToken
// @Resp
type LoginToken struct {
	UserName string    `json:"username"`
	Token    string    `json:"token"`
	UserID   uint      `json:"user_id"`
	IsAdmin  bool      `json:"is_admin"`
	Expired  time.Time `json:"expired"`
	Avatar   string    `json:"avatar"`
	Session  string    `json:"session"`
}

// SignIn 登录
// @POST /signIn
// @Ignore Session
func (u *UserController) SignIn(ctx *fw.Context, req *SignInRequest, s *session.Store) {
	var user models.User
	if err := u.DB.Where("Name = ?", req.Username).Or("Email=?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(404, base.Resp[int, any](404, "用户不存在", nil))
			return
		}
	}
	if !user.Enabled {
		ctx.JSON(403, base.Resp[int, any](403, "用户已停用, 请联系管理员操作", nil))
		return
	}
	if user.Password == req.Password {
		user.LastLoginTime = time.Now()
		user.LastLoginIp = ctx.RemoteIP()
		u.DB.Save(&user)
		u.SetAuthed(s, map[string]any{
			"user_id":   user.ID,
			"user_name": user.Name,
			"isAdmin":   user.IsAdmin,
		})

		ctx.JSON(200, base.Resp(200, "ok", conv.String(s.GetSessionID())))
		return
	} else {
		ctx.JSON(403, base.Resp[int, any](401, "登录失败, 用户名或密码不正确", nil))
	}

}

// SignOut 退出
// @POST /signOut
func (u *UserController) SignOut(ctx *fw.Context, store *session.Store) {
	u.ClearSession(store)
	ctx.JSON(200, base.Resp[int, any](200, "ok", nil))
}

// SignUpRequest
// @Body
type SignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"email"` //邮箱地址
}

// SignUp 注册
// @POST /signUp
// @Ignore Session
func (u *UserController) SignUp(ctx *fw.Context, req *SignUpRequest) {
	//panic("ss")
	var user = new(models.User)
	if err := u.DB.Where("name=?", req.Username).Or("email=?", req.Email).First(user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(500, map[string]interface{}{
				"code":    500,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
	}
	if _, ok := user.GetID(); ok {
		ctx.JSON(200, base.Resp[int, any](200, "用户已经存在", nil))
		return
	}
	var count int64
	u.DB.Model(&models.User{}).Count(&count)
	isAdmin := false
	// the first user is administrator
	if count == 0 {
		isAdmin = true
	}
	user1 := &models.User{
		BaseModel:   models2.NewBaseModel(0),
		Name:        req.Username,
		Email:       req.Email,
		Avatar:      "", //default avatar
		Password:    req.Password,
		IsAdmin:     isAdmin,
		Enabled:     isAdmin,
		LastLoginIp: "",
	}
	if err := u.DB.Create(&user1).Error; err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	ctx.JSON(200, base.Resp(200, "ok", ""))

}

// UploadFormRequest
// @Multipart
type UploadFormRequest struct {
	File *multipart.FileHeader `multipart:"file"`
}

// UploadAvatar 上传头像
// @POST /avatar/upload
func (u *UserController) UploadAvatar(ctx *fw.Context, req *UploadFormRequest) {

	filename := fw.UUID() + filepath.Ext(req.File.Filename)
	full := filepath.Join("./static/", filename)
	err := ctx.SaveUploadFile(req.File, full)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(200,
		base.Resp(200, "ok", map[string]interface{}{
			"path": "/static/" + filename,
			"id":   filename,
		}))

}

// Profile
// @GET /profile
func (u *UserController) Profile(ctx *fw.Context, store *session.Store) {
	if conv.Uint(u.CurrentUserID(store)) > 0 {
		user := &models.User{}
		u.DB.Model(new(models.User)).Where("id = ?", u.CurrentUserID(store)).First(user)

		ctx.JSON(200, base.Resp(200, "ok", user))
	}
}

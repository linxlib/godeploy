package base

import (
	"errors"
	"github.com/fasthttp/session/v2"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/base/models"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

func NewCrudController[T models.IBase[uint]](db *gorm.DB) *CrudController[T] {
	c := &CrudController[T]{
		db: db,
	}
	return c
}

// CrudController
// a base crud controller for models which has one primary field
type CrudController[T models.IBase[uint]] struct {
	db *gorm.DB
}

func (c *CrudController[T]) CheckSession(store *session.Store) bool {
	if store == nil {
		return false
	}
	return store.Get("user_id") != nil
}
func (c *CrudController[T]) SetAuthed(store *session.Store, data map[string]any) {
	for key, val := range data {
		store.Set(key, val)
	}
}
func (c *CrudController[T]) ClearSession(store *session.Store) {
	store.Flush()
}

// GetByID 根据ID获取
// @GET /
func (c *CrudController[T]) GetByID(ctx *fw.Context, q *IDQuery2) {
	v := new(T)
	err := c.db.First(v, q.ID).Error
	if err != nil {
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
	ctx.JSON(fasthttp.StatusOK, Resp(200, "ok", v))
}

// GetPageList 获取分页
// @GET /list
func (c *CrudController[T]) GetPageList(ctx *fw.Context, q *PageSize2, store *session.Store) {
	var vs []T
	var count int64
	d := c.db
	if q.Search != "" {
		d = d.Where("name like ?", "%"+q.Search+"%")
	}
	d1 := d.Model(new(T))
	d1.Count(&count)
	d1.Offset(q.Offset()).Limit(q.Size).Find(&vs)
	ctx.JSON(fasthttp.StatusOK, Resp(200, "ok", ListData(vs, count, q.Size)))
}

// InsertOrUpdate 插入或修改
// @POST /
func (c *CrudController[T]) InsertOrUpdate(ctx *fw.Context, body T, store *session.Store) {
	db := c.db
	m := body.CheckExistColumns()
	for col, value := range m {
		db = db.Where(col, value)
	}
	var ex int64
	db.Model(new(T)).Count(&ex)
	if ex > 0 {
		ctx.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "已存在",
			"data":    nil,
		})
		return
	}

	var err error
	if _, ok := body.GetID(); ok {
		err = c.db.Model(body).Updates(body).Error
	} else {
		err = c.db.Create(body).Error
	}
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "ok",
	})

}

// Delete
// @POST /delete
func (c *CrudController[T]) Delete(ctx *fw.Context, d *DeleteBody) {
	err := c.db.Delete(new(T), d.IDS).Error
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"code":    500,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"code":    200,
		"message": "ok",
	})
}

package base

import (
	"errors"
	"github.com/fasthttp/session/v2"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/base/models"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"math"
)

func NewSimpleCrudController[E models.PrimaryKey, T models.IBase[E]](db *gorm.DB) *SimpleCrudController[E, T] {
	c := &SimpleCrudController[E, T]{
		db: db,
	}
	return c
}

// SimpleCrudController
// a base crud controller for models which has one primary field
type SimpleCrudController[E models.PrimaryKey, T models.IBase[E]] struct {
	db *gorm.DB
}

func (c *SimpleCrudController[E, T]) CheckSession(store *session.Store) bool {
	if store == nil {
		return false
	}
	return store.Get("user_id") != nil
}
func (c *SimpleCrudController[E, T]) SetAuthed(store *session.Store, data map[string]any) {
	for key, val := range data {
		store.Set(key, val)
	}
}
func (c *SimpleCrudController[E, T]) ClearSession(store *session.Store) {
	store.Flush()
}

// IDQuery2
// @Query
type IDQuery2 struct {
	ID uint `query:"id"` //id
}

// GetByID 根据ID获取
// @GET /
func (c *SimpleCrudController[E, T]) GetByID(ctx *fw.Context, q *IDQuery2) {
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

type IData interface {
	any | ListDataBase[any]
}

type RespBase[E int | string, T IData] struct {
	Code    E      `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func RespInt[T IData](code int, message string, data T) RespBase[int, T] {
	return RespBase[int, T]{
		Code:    code,
		Message: message,
		Data:    data,
	}

}
func Resp[E int | string, T IData](code E, message string, data T) RespBase[E, T] {
	return RespBase[E, T]{
		Code:    code,
		Message: message,
		Data:    data,
	}

}
func RespString[T IData](code string, message string, data T) RespBase[string, T] {
	return RespBase[string, T]{
		Code:    code,
		Message: message,
		Data:    data,
	}

}

// PageSize2
// @Query
type PageSize2 struct {
	models.PageSizeBase
	Search string `query:"search" default:""` //搜索名称
}

type ListDataBase[T any] struct {
	List      []T   `json:"list"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"totalPage"`
}

func ListData[T any](list []T, total int64, size int) ListDataBase[T] {
	return ListDataBase[T]{
		List:      list,
		Total:     total,
		TotalPage: int(math.Ceil(float64(total) / float64(size))),
	}
}

// GetPageList 获取分页
// @GET /list
func (c *SimpleCrudController[E, T]) GetPageList(ctx *fw.Context, q *PageSize2, store *session.Store) {
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
func (c *SimpleCrudController[E, T]) InsertOrUpdate(ctx *fw.Context, body T, store *session.Store) {
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

// DeleteBody
// @Body
type DeleteBody struct {
	IDS []uint `json:"ids"`
}

// Delete 删除
// @POST /delete
func (c *SimpleCrudController[E, T]) Delete(ctx *fw.Context, d *DeleteBody) {
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

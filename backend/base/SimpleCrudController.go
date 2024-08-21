package base

import (
	"errors"
	"github.com/fasthttp/session/v2"
	"github.com/linxlib/fw"
	"github.com/linxlib/godeploy/base/models"
	"github.com/linxlib/inject"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"math"
)

// NewSimpleCrudController returns a new SimpleCrudController.
//
// It's a base crud controller for models which has one primary field.
func NewSimpleCrudController[E models.PrimaryKey, T models.IBase[E]]() *SimpleCrudController[E, T] {
	c := &SimpleCrudController[E, T]{}
	return c
}

// SimpleCrudController
// a base crud controller for models which has one primary field
type SimpleCrudController[E models.PrimaryKey, T models.IBase[E]] struct {
	DB *gorm.DB
}

func (c *SimpleCrudController[E, T]) Init(provider inject.Provider) {
	c.DB = &gorm.DB{}
	provider.Provide(c.DB)
}

// CheckSession checks if the session is authenticated.
//
// If the session is not authenticated, this method returns false.
// Otherwise, it returns true.
func (c *SimpleCrudController[E, T]) CheckSession(store *session.Store) bool {
	// If the session is nil, it's not authenticated.
	if store == nil {
		return false
	}
	// If the user_id key is not set in the session, it's not authenticated.
	return store.Get("user_id") != nil
}

// SetAuthed sets the provided data into the session store.
//
// Parameters:
// - store: the session store to set the data into.
// - data: a map of key-value pairs to set in the session store.
func (c *SimpleCrudController[E, T]) SetAuthed(store *session.Store, data map[string]any) {
	// Iterate over each key-value pair in the data map.
	for key, val := range data {
		// Set the key-value pair in the session store.
		store.Set(key, val)
	}
}

// ClearSession clears the session store.
//
// This method is used to clear the session after the user logs out.
func (c *SimpleCrudController[E, T]) ClearSession(store *session.Store) {
	// Clear the session store.
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
	err := c.DB.First(v, q.ID).Error
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
	any | ListDataBase[any, int64] | ListDataBase[any, int]
}

type RespBase[E int | string, T IData] struct {
	Code    E      `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// RespInt is a helper function that returns a RespBase[int, T] struct.
// It is used to create a response with an integer code and data.
//
// Parameters:
// - code: the integer code to be set in the RespBase struct.
// - message: the message to be set in the RespBase struct.
// - data: the data to be set in the RespBase struct.
//
// Returns:
// - RespBase[int, T]: the created RespBase struct.
func RespInt[T IData](code int, message string, data T) RespBase[int, T] {
	// Create a new RespBase struct with the provided parameters.
	return RespBase[int, T]{
		Code:    code,    // Set the code field to the provided code.
		Message: message, // Set the message field to the provided message.
		Data:    data,    // Set the data field to the provided data.
	}
}

// Resp is a helper function that returns a RespBase[E, T] struct.
// It is used to create a response with a generic code and data.
//
// Parameters:
// - code: the generic code to be set in the RespBase struct.
// - message: the message to be set in the RespBase struct.
// - data: the data to be set in the RespBase struct.
//
// Returns:
// - RespBase[E, T]: the created RespBase struct.
func Resp[E int | string, T IData](code E, message string, data T) RespBase[E, T] {
	// Create a new RespBase struct with the provided parameters.
	// The code, message, and data fields are set to the provided values.
	return RespBase[E, T]{
		Code:    code,    // Set the code field to the provided code.
		Message: message, // Set the message field to the provided message.
		Data:    data,    // Set the data field to the provided data.
	}
}

// RespString is a helper function that returns a RespBase[string, T] struct.
// It is used to create a response with a string code and data.
//
// Parameters:
// - code: the string code to be set in the RespBase struct.
// - message: the message to be set in the RespBase struct.
// - data: the data to be set in the RespBase struct.
//
// Returns:
// - RespBase[string, T]: the created RespBase struct.
func RespString[T IData](code string, message string, data T) RespBase[string, T] {
	// Create a new RespBase struct with the provided parameters.
	// The code, message, and data fields are set to the provided values.
	return RespBase[string, T]{
		Code:    code,    // Set the code field to the provided code.
		Message: message, // Set the message field to the provided message.
		Data:    data,    // Set the data field to the provided data.
	}
}

// PageSize2
// @Query
type PageSize2 struct {
	models.PageSizeBase
	Search string `query:"search" default:""` //搜索名称
}

type ListDataBase[T any, T1 int | int64] struct {
	List      []T `json:"list"`
	Total     T1  `json:"total"`
	TotalPage int `json:"totalPage"`
}

func (l ListDataBase[T, T1]) Resp() RespBase[int, ListDataBase[T, T1]] {
	return Resp(200, "ok", l)
}

// ListData returns a ListDataBase struct containing the provided list, total, and totalPage.
//
// Parameters:
// - list: the list of items to be included in the ListDataBase struct.
// - total: the total number of items to be included in the ListDataBase struct.
// - size: the number of items per page to be included in the ListDataBase struct.
//
// Returns:
// - ListDataBase[T]: the created ListDataBase struct.
func ListData[T any, T1 int | int64](list []T, total T1, size int) ListDataBase[T, T1] {
	// Calculate the total number of pages based on the total number of items and the number of items per page.
	// The math.Ceil function is used to round up the result to the nearest integer.
	totalPage := int(math.Ceil(float64(total) / float64(size)))

	// Create a new ListDataBase struct with the provided parameters.
	// The List, Total, and TotalPage fields are set to the provided values.
	return ListDataBase[T, T1]{
		List:      list,      // Set the List field to the provided list.
		Total:     total,     // Set the Total field to the provided total.
		TotalPage: totalPage, // Set the TotalPage field to the calculated totalPage.
	}
}

func Data[T any](data T) RespBase[int, T] {
	return Resp(200, "ok", data)
}

// GetPageList 获取分页
// @GET /list
func (c *SimpleCrudController[E, T]) GetPageList(ctx *fw.Context, q *PageSize2) {
	var vs []T
	var count int64
	d := c.DB
	if q.Search != "" {
		d = d.Where("name like ?", "%"+q.Search+"%")
	}
	d1 := d.Model(new(T))
	d1.Count(&count)
	d.Model(new(T)).Offset(q.Offset()).Limit(q.Size).Find(&vs)
	ctx.JSON(fasthttp.StatusOK, Resp(200, "ok", ListData(vs, count, q.Size)))
}

// InsertOrUpdate 插入或修改
// @POST /
func (c *SimpleCrudController[E, T]) InsertOrUpdate(ctx *fw.Context, body T) {
	db := c.DB
	m := body.CheckExistColumns()
	for col, value := range m {
		db = db.Where(col, value)
	}

	var err error
	if _, ok := body.GetID(); ok {
		err = c.DB.Model(body).Select("*").Omit("created_at", "updated_at", "deleted_at").Updates(body).Error
	} else {
		var ex int64
		db.Model(new(T)).Count(&ex)
		if ex > 0 {
			ctx.JSON(200, map[string]interface{}{
				"code":    400,
				"message": "已存在",
				"data":    nil,
			})
			return
		}
		err = c.DB.Create(body).Error
	}
	if err != nil {
		ctx.JSON(200, map[string]interface{}{
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
	err := c.DB.Delete(new(T), d.IDS).Error
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

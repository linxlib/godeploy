package models

import (
	"go/constant"
	"go/token"
	"gorm.io/gorm"
)

type PrimaryKey interface {
	comparable
}
type IBase[T PrimaryKey] interface {
	GetID() (T, bool)
	CheckExistColumns() map[string]any
}
type Base[T PrimaryKey] struct {
	ID T `gorm:"primaryKey" json:"id"`
}

func (b *Base[T]) CheckExistColumns() map[string]any {
	return map[string]any{}
}

func (b *Base[T]) GetID() (T, bool) {
	v := constant.Make(b.ID)
	switch v.Kind() {
	case constant.Int:
		if constant.Compare(v, token.EQL, constant.Make(0)) {
			return b.ID, true
		}
		return b.ID, false
	case constant.String:
		if constant.Compare(v, token.EQL, constant.MakeString("")) {
			return b.ID, true
		}
		return b.ID, false
	default:
		return b.ID, false
	}
}

var _ IBase[uint] = (*BaseModel)(nil)

type BaseModel struct {
	gorm.Model
}

func (b *BaseModel) CheckExistColumns() map[string]any {
	return map[string]any{}
}

func NewBaseModel(id uint) *BaseModel {
	return &BaseModel{
		Model: gorm.Model{
			ID: id,
		},
	}
}

func (b *BaseModel) GetID() (uint, bool) {

	if b == nil || b.ID == 0 {
		return 0, false
	}
	return b.ID, true
}

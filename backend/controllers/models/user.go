package models

import (
	"github.com/linxlib/godeploy/base/models"
	"time"
)

var _ models.IBase[uint] = (*User)(nil)

// @Body
type User struct {
	*models.BaseModel
	Name          string `gorm:"type:varchar(255);unique_index;not null" validate:"required"`
	Email         string `gorm:"type:varchar(255);unique_index;not null" validate:"email"`
	Avatar        string `gorm:"type:varchar(255)"`
	Password      string `gorm:"type:varchar(255);not null"`
	IsAdmin       bool   `gorm:"type:boolean;not null"`
	Enabled       bool   `gorm:"type:boolean;not null"`
	LastLoginTime time.Time
	LastLoginIp   string `gorm:"type:varchar(255)"`
}

func (u *User) GetID() (uint, bool) {

	if u.ID == 0 {
		return 0, false
	}
	return u.ID, true
}

func (u *User) CheckExistColumns() map[string]any {
	m := make(map[string]any)
	m["name"] = u.Name
	m["email"] = u.Email
	return m
}

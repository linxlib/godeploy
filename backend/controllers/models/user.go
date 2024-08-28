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
	Password      string `gorm:"type:varchar(255);not null" json:"-"`
	IsAdmin       bool   `gorm:"type:boolean;not null"`
	Enabled       bool   `gorm:"type:boolean;not null"`
	LastLoginTime time.Time
	LastLoginIp   string `gorm:"type:varchar(255)"`
}

// GetID returns the ID of the user and a boolean indicating if the ID is valid.
//
// If the ID is not valid, the second return value will be false.
// Returns the user's ID and a boolean indicating if the ID is valid.
func (u *User) GetID() (uint, bool) {
	if u.BaseModel == nil {
		return 0, false
	}
	// Check if the ID is valid
	if u.ID == 0 {
		// If the ID is not valid, return 0 and false
		return 0, false
	}
	// Return the ID and true
	return u.ID, true
}

func (u *User) CheckExistColumns() map[string]any {
	m := make(map[string]any)
	m["name"] = u.Name
	m["email"] = u.Email
	return m
}

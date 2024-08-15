package db

import (
	"errors"
	"github.com/glebarez/sqlite"
	"github.com/linxlib/config"
	"github.com/linxlib/godeploy/controllers/models"
	"gorm.io/gorm"
)

type DBMapper struct {
	Type        string `yaml:"type" default:"sqlite"`
	Dsn         string `yaml:"dsn" default:"deploy.db"`
	AutoMigrate bool   `yaml:"autoMigrate" default:"true"`
}

func (D *DBMapper) Init(config *config.Config) (any, error) {
	var err error
	var db *gorm.DB
	err = config.LoadWithKey("db", D)
	if err != nil {
		return nil, err
	}

	if D.Type == "sqlite" {
		db, err = gorm.Open(sqlite.Open("deploy.db"))
		if err != nil {
			return nil, errors.New("connect to db failed: " + err.Error())
		}
	} else {
		return nil, errors.New("unsupported database type: " + D.Type)
	}
	if D.AutoMigrate {
		err = db.AutoMigrate(&models.User{}, &models.Service{})
		if err != nil {
			return nil, errors.New("db migration failed:" + err.Error())
		}
	}
	return db, nil

}

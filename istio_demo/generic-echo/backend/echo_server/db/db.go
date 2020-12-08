package database

import (
	"errors"
	"fmt"
	"server/conf"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var instance *gorm.DB
var HARD_CODE_DELAY int

type Player struct {
	gorm.Model
	Name     string
	Account  string
	Password string
}

func InitDbConn(config conf.Database) error {
	if !config.External {
		return nil
	}

	conninfo := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=true", config.User, config.Passwd, config.Host, config.Db)
	HARD_CODE_DELAY = config.HardcodeDelay
	var err error
	instance, err = gorm.Open(mysql.Open(conninfo), &gorm.Config{}) //gorm.Open("mysql", conninfo)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect to db '%s', err: %s", conninfo, err.Error())
		return errors.New(msg)
	}

	return nil
}

func GetTestingSQLService() string {
	var player Player
	player.Name = "Guest"
	if instance != nil {
		instance.First(&player, "account = ?", "test123")
		time.Sleep(time.Duration(HARD_CODE_DELAY) * time.Second)
	}
	return player.Name
}

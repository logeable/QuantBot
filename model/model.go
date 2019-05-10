package model

import (
	"log"
	"time"

	"github.com/hprose/hprose-golang/io"
	"github.com/jinzhu/gorm"
	"github.com/phonegapX/QuantBot/config"

	// for db SQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	// DB Database
	DB *gorm.DB
)

func init() {
	io.Register((*User)(nil), "User", "json")
	io.Register((*Exchange)(nil), "Exchange", "json")
	io.Register((*Algorithm)(nil), "Algorithm", "json")
	io.Register((*Trader)(nil), "Trader", "json")
	io.Register((*Log)(nil), "Log", "json")
	var err error
	DB, err = gorm.Open(config.Config.Database.Driver, config.Config.Database.DSN)
	if err != nil {
		log.Fatalf("Connect database error: %v\n", err)
	}
	DB.AutoMigrate(&User{}, &Exchange{}, &Algorithm{}, &TraderExchange{}, &Trader{}, &Log{})
	users := []User{}
	DB.Find(&users)
	if len(users) == 0 {
		admin := User{
			Username: "admin",
			Password: "admin",
			Level:    99,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalln("Create admin error:", err)
		}
	}
	DB.LogMode(false)
	go ping()
}

func ping() {
	for {
		if err := DB.Exec("SELECT 1").Error; err != nil {
			log.Println("Database ping error:", err)
			if DB, err = gorm.Open(config.Config.Database.Driver, config.Config.Database.DSN); err != nil {
				log.Println("Retry connect to database error:", err)
			}
		}
		time.Sleep(time.Minute)
	}
}

// NewOrm ...
func NewOrm() (*gorm.DB, error) {
	return gorm.Open(config.Config.Database.Driver, config.Config.Database.DSN)
}

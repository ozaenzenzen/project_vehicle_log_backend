package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func SetupDB() *gorm.DB {
	USER := "root"
	PASS := ""
	// HOST := "localhost"
	HOST := "127.0.0.1"
	// PORT := "3306"
	PORT := "3306"
	DBNAME := "project_vehicle_log_backend" // staging
	// DBNAME := "project_vehicle_log_backend2" // development
	// URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open("mysql", URL)
	if err != nil {
		panic(err.Error())
	}
	return db

}

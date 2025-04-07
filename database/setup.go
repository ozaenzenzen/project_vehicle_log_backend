package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	configModel "project_vehicle_log_backend/data"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func DBConfig() (*string, *string, *string, *string, *string) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
		return nil, nil, nil, nil, nil
	}

	var config configModel.Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
		return nil, nil, nil, nil, nil
	}

	fmt.Println("Secret:", config.Secret)
	fmt.Println("User:", config.User)
	fmt.Println("Pass:", config.Pass)
	fmt.Println("Port:", config.Port)
	fmt.Println("DBName:", config.DBName)
	fmt.Println("Host:", config.Host)
	return &config.User, &config.Pass, &config.Port, &config.DBName, &config.Host
}

func SetupDB() *gorm.DB {
	user, pass, port, dbName, host := DBConfig()

	USER := *user
	PASS := *pass
	HOST := *host
	PORT := *port
	DBNAME := *dbName
	// DBNAME := "project_vehicle_log_backend" // staging
	// DBNAME := "project_vehicle_log_backend2" // development
	// URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open("mysql", URL)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func SetupDBOld() *gorm.DB {
	USER := "root"
	PASS := ""
	// HOST := "localhost"
	HOST := "127.0.0.1"
	// PORT := "3306"
	PORT := "3306"
	DBNAME := "project_vehicle_log_backend_development" // development
	// DBNAME := "project_vehicle_log_backend" // staging
	// DBNAME := "project_vehicle_log_backend2" // development
	// URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	db, err := gorm.Open("mysql", URL)
	if err != nil {
		panic(err.Error())
	}
	return db
}

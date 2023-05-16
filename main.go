package main

import (
	setup "project_vehicle_log_backend/database"
	routes "project_vehicle_log_backend/routes"
)

func main() {
	db := setup.SetupDB()
	sqlDB := db.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)
	// db.AutoMigrate(&account.AccountUserV2{}, &event.Event{}, &event.VolunteersForm{}, &notif.Notification{})

	r := routes.SetupRoutes(db)
	err := r.Run(":8080")

	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
	return

}

package main

import (
	setup "project_vehicle_log_backend/database"
	account "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"
	vehicle "project_vehicle_log_backend/models/vehicle"
	routes "project_vehicle_log_backend/routes"
)

func main() {
	db := setup.SetupDB()
	sqlDB := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	db.AutoMigrate(&account.AccountUserModel{}, &vehicle.VehicleModel{}, &notif.Notification{})

	r := routes.SetupRoutes(db)
	err := r.Run(":8080")

	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
	return

}

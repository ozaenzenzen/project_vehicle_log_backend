package controllers

import (
	"fmt"
	"log"
	"net/http"

	notif "project_vehicle_log_backend/models/notification"
	vehicle "project_vehicle_log_backend/models/vehicle"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

type CreateVehicleResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CreateVehicleReqeust struct {
	UserId         uint   `gorm:"not null" json:"user_id" validate:"required"`
	VehicleName    string `gorm:"not null" json:"vehicle_name" validate:"required"`
	VehicleImage   string `gorm:"not null" json:"vehicle_image" validate:"required"`
	Year           string `gorm:"not null" json:"year" validate:"required"`
	EngineCapacity string `gorm:"not null" json:"engine_capacity" validate:"required"`
	TankCapacity   string `gorm:"not null" json:"tank_capacity" validate:"required"`
	Color          string `gorm:"not null" json:"color" validate:"required"`
	MachineNumber  string `gorm:"not null" json:"machine_number" validate:"required"`
	ChassisNumber  string `gorm:"not null" json:"chassis_number" validate:"required"`
}

func CreateVehicle(c *gin.Context) {
	var vehicleInput CreateVehicleReqeust
	if err := c.ShouldBindJSON(&vehicleInput); err != nil {
		log.Println(fmt.Sprintf("error log: %s", err))
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  500,
			Message: "Create Vehicle Failed",
		})
		return
	}
	validate := validator.New()

	if err := validate.Struct(vehicleInput); err != nil {
		log.Println(fmt.Sprintf("error log2: %s", err))
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  400,
			Message: "Create Vehicle Failed2",
		})
		return
	}

	vehicleData := vehicle.VehicleModel{
		UserId:         vehicleInput.UserId,
		VehicleName:    vehicleInput.VehicleName,
		VehicleImage:   vehicleInput.VehicleImage,
		Year:           vehicleInput.Year,
		EngineCapacity: vehicleInput.EngineCapacity,
		TankCapacity:   vehicleInput.TankCapacity,
		Color:          vehicleInput.Color,
		MachineNumber:  vehicleInput.MachineNumber,
		ChassisNumber:  vehicleInput.ChassisNumber,
	}

	createVehicleResponse := CreateVehicleResponse{
		Status:  201,
		Message: "Vehicle created successfully",
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}
	result := db.Create(&vehicleData)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  vehicleInput.UserId,
		NotificationTitle:       "Add Vehicle",
		NotificationDescription: "Anda Telah Menambahkan Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	c.JSON(http.StatusCreated, createVehicleResponse)
}

package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	jwthelper "project_vehicle_log_backend/helper"
	account "project_vehicle_log_backend/models/account"
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

type CreateVehicleRequest struct {
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
	var vehicleInput CreateVehicleRequest
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

	db := c.MustGet("db").(*gorm.DB)

	//--------check id--------check id--------check id--------

	checkID := db.Table("account_user_models").Where("id = ?", vehicleInput.UserId).Find(&account.AccountUserModel{
		ID: vehicleInput.UserId,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

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

	if db.Error != nil {
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}
	// result := db.Create(&vehicleData)
	result := db.FirstOrCreate(&vehicleData, vehicle.VehicleModel{
		MachineNumber: vehicleInput.MachineNumber,
		ChassisNumber: vehicleInput.ChassisNumber,
	})

	if result.Value == nil && result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, CreateVehicleResponse{
			Status:  400,
			Message: "Record found",
		})
		return
	}

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

type EditVehicleResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type EditVehicleReqeust struct {
	// UserId         uint   `gorm:"not null" json:"user_id" validate:"required"`
	VehicleName    string `gorm:"not null" json:"vehicle_name" validate:"required"`
	VehicleImage   string `gorm:"not null" json:"vehicle_image" validate:"required"`
	Year           string `gorm:"not null" json:"year" validate:"required"`
	EngineCapacity string `gorm:"not null" json:"engine_capacity" validate:"required"`
	TankCapacity   string `gorm:"not null" json:"tank_capacity" validate:"required"`
	Color          string `gorm:"not null" json:"color" validate:"required"`
	MachineNumber  string `gorm:"not null" json:"machine_number" validate:"required"`
	ChassisNumber  string `gorm:"not null" json:"chassis_number" validate:"required"`
}

func EditVehicle(c *gin.Context) {
	headerid := c.Request.Header.Get("usd")

	if headerid == "" {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: "headers empty",
		})
		return
	}
	var editVehicleRequest EditVehicleReqeust

	if err := c.ShouldBindJSON(&editVehicleRequest); err != nil {
		log.Println(fmt.Sprintf("error log: %s", err))
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  500,
			Message: "Edit Vehicle Failed",
		})
		return
	}
	validate := validator.New()

	if err := validate.Struct(editVehicleRequest); err != nil {
		log.Println(fmt.Sprintf("error log2: %s", err))
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: "Edit Vehicle Failed2",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(headerid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", headerid).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	vehicleDataOutput := vehicle.VehicleModel{
		UserId:         iduint,
		VehicleName:    editVehicleRequest.VehicleName,
		VehicleImage:   editVehicleRequest.VehicleImage,
		Year:           editVehicleRequest.Year,
		EngineCapacity: editVehicleRequest.EngineCapacity,
		TankCapacity:   editVehicleRequest.TankCapacity,
		Color:          editVehicleRequest.Color,
		MachineNumber:  editVehicleRequest.MachineNumber,
		ChassisNumber:  editVehicleRequest.ChassisNumber,
	}

	editVehicleResponse := EditVehicleResponse{
		Status:  201,
		Message: "Vehicle update successfully",
	}

	if db.Error != nil {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}
	// result := db.Create(&vehicleDataOutput)
	result := db.Table("vehicle_models").Where("id = ?", headerid).Update(&vehicleDataOutput)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  iduint,
		NotificationTitle:       "Edit Vehicle",
		NotificationDescription: "Anda Telah Mengubah Data Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	c.JSON(http.StatusCreated, editVehicleResponse)
}

type VehicleData struct {
	Id             uint   `json:"id" gorm:"primary_key"`
	UserId         uint   `json:"user_id"`
	VehicleName    string `json:"vehicle_name"`
	VehicleImage   string `json:"vehicle_image"`
	Year           string `json:"year"`
	EngineCapacity string `json:"engine_capacity"`
	TankCapacity   string `json:"tank_capacity"`
	Color          string `json:"color"`
	MachineNumber  string `json:"machine_number"`
	ChassisNumber  string `json:"chassis_number"`
	// VehicleMeasurementLogModel []vehicle.VehicleMeasurementLogModel `json:"vehicle_measurement_log_models" gorm:"foreignKey:user_id;references:UserId"`
	VehicleMeasurementLogModels []VehicleMeasurementLogModel `json:"vehicle_measurement_log_models" gorm:"foreignKey:vehicle_id;references:Id"`
}

type VehicleMeasurementLogModel struct {
	Id                  uint      `json:"id" `
	UserId              uint      `json:"user_id"`
	VehicleId           uint      `json:"vehicle_id"`
	MeasurementTitle    string    `json:"measurement_title"`
	CurrentOdo          string    `json:"current_odo"`
	EstimateOdoChanging string    `json:"estimate_odo_changing"`
	AmountExpenses      string    `json:"amount_expenses"`
	CheckpointDate      string    `json:"checkpoint_date"`
	Notes               string    `json:"notes"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type GetAllVehicleDataResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	// Data    []vehicle.VehicleModel
	Data []VehicleData
}

func GetAllVehicleData(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, GetAllVehicleDataResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, GetAllVehicleDataResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)

	// var vehicleData []vehicle.VehicleModel
	var vehicleData []VehicleData

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, GetAllVehicleDataResponse{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", returnEmailsOrUid).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, GetAllVehicleDataResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	result := db.Preload("VehicleMeasurementLogModels", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, user_id, vehicle_id, measurement_title, current_odo, estimate_odo_changing, amount_expenses, checkpoint_date, notes, created_at, updated_at")
	}).Table("vehicle_models").Where("user_id = ?", returnEmailsOrUid).Find(&vehicleData)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, GetAllVehicleDataResponse{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	if result.Value == nil {
		log.Println(fmt.Sprintf("error log3: %s", result.Error))
		c.JSON(http.StatusNotFound, GetAllVehicleDataResponse{
			Status:  404,
			Message: "get all vehicle data Failed",
			Data:    []VehicleData{},
			// Data:    []vehicle.VehicleModel{},
		})
		return
	}

	c.JSON(http.StatusOK, GetAllVehicleDataResponse{
		Status:  200,
		Message: "get all vehicle data success",
		Data:    vehicleData,
	})
}

type CreateLogVehicleRequest struct {
	UserId              uint   `json:"user_id"`
	VehicleId           uint   `json:"vehicle_id"`
	MeasurementTitle    string `json:"measurement_title"`
	CurrentOdo          string `json:"current_odd"`
	EstimateOdoChanging string `json:"estimate_odo_changing"`
	AmountExpenses      string `json:"amount_expenses"`
	CheckpointDate      string `json:"checkpoint_date"`
	Notes               string `json:"notes"`
}
type CreateLogVehicleResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func CreateLogVehicle(c *gin.Context) {
	var createLogVehicle CreateLogVehicleRequest
	if err := c.ShouldBindJSON(&createLogVehicle); err != nil {
		log.Println(fmt.Sprintf("error log JoinEvent1: %s", err))
		c.JSON(http.StatusBadRequest, CreateLogVehicleResponse{
			Status:  500,
			Message: "Create Log Vehicle Failed1",
		})
		return
	}

	createLogVehicleResponse := CreateLogVehicleResponse{
		Status:  201,
		Message: "Create Log Vehicle successfully",
	}

	createLogVehicleData := vehicle.VehicleMeasurementLogModel{
		UserId:              createLogVehicle.UserId,
		VehicleId:           createLogVehicle.VehicleId,
		MeasurementTitle:    createLogVehicle.MeasurementTitle,
		CurrentOdo:          createLogVehicle.CurrentOdo,
		EstimateOdoChanging: createLogVehicle.EstimateOdoChanging,
		AmountExpenses:      createLogVehicle.AmountExpenses,
		CheckpointDate:      createLogVehicle.CheckpointDate,
		Notes:               createLogVehicle.Notes,
	}
	// Volunteers Status
	// 0: idle
	// 1: accepted
	// 2: rejected
	// 3: waiting response

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		c.JSON(http.StatusBadRequest, CreateLogVehicleResponse{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}

	// check := db.Where("user_id = ?", createLogVehicleData.UserId).Where("vehicle_id = ?", createLogVehicleData.VehicleId).First(&vehicle.VehicleMeasurementLogModel{})
	// log.Println(fmt.Sprintf("error check error: %s", check.Error))
	// log.Println(fmt.Sprintf("error check value: %s", check.Value))
	// if check.Error == nil {
	// 	c.JSON(http.StatusBadRequest, CreateLogVehicleResponse{
	// 		Status:  400,
	// 		Message: "Already Submitted",
	// 	})
	// 	return
	// }

	result := db.Create(&createLogVehicleData)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, CreateLogVehicleResponse{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}
	inputNotifModel := notif.Notification{
		UserId:                  createLogVehicle.UserId,
		NotificationTitle:       "Add Vehicle Log",
		NotificationDescription: "Anda Telah Menambahkan Log Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusBadRequest, CreateLogVehicleResponse{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	c.JSON(http.StatusCreated, createLogVehicleResponse)
}

type GetLogVehicleRequest struct {
	UserId              uint   `json:"user_id"`
	VehicleId           uint   `json:"vehicle_id"`
	MeasurementTitle    string `json:"measurement_title"`
	CurrentOdo          string `json:"current_odo"`
	EstimateOdoChanging string `json:"estimate_odo_changing"`
	AmountExpenses      string `json:"amount_expenses"`
	CheckpointDate      string `json:"checkpoint_date"`
	Notes               string `json:"notes"`
}
type GetLogVehicleResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []vehicle.VehicleMeasurementLogModel
}

func GetLogVehicle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, GetAllVehicleDataResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, GetAllVehicleDataResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)
	var logVehicleData []vehicle.VehicleMeasurementLogModel

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, GetLogVehicleResponse{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", returnEmailsOrUid).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, GetLogVehicleResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	result := db.Table("vehicle_measurement_log_models").Where("user_id = ?", returnEmailsOrUid).Find(&logVehicleData)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, GetLogVehicleResponse{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	if result.Value == nil {
		log.Println(fmt.Sprintf("error log3: %s", result.Error))
		c.JSON(http.StatusNotFound, GetLogVehicleResponse{
			Status:  404,
			Message: "get log vehicle data Failed",
			Data:    []vehicle.VehicleMeasurementLogModel{},
		})
		return
	}

	c.JSON(http.StatusOK, GetLogVehicleResponse{
		Status:  200,
		Message: "get log vehicle data success",
		Data:    logVehicleData,
	})
}

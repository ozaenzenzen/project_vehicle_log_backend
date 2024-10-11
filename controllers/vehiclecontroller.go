package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	req "project_vehicle_log_backend/data/vehicle/request"
	resp "project_vehicle_log_backend/data/vehicle/response"
	jwthelper "project_vehicle_log_backend/helper"
	account "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"
	vehicle "project_vehicle_log_backend/models/vehicle"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

func CreateVehicle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)
	var createVehicleReq req.CreateVehicleRequestModel
	if err := c.ShouldBindJSON(&createVehicleReq); err != nil {
		log.Println(fmt.Sprintf("error log: %s", err))
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  500,
			Message: "Create Vehicle Failed",
		})
		return
	}
	validate := validator.New()

	if err := validate.Struct(createVehicleReq); err != nil {
		log.Println(fmt.Sprintf("error log2: %s", err))
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  400,
			Message: "validate error, field required",
		})
		return
	}

	//--------check id--------check id--------check id--------
	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", iduint).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	vehicleData := vehicle.VehicleModel{
		UserId:         iduint,
		VehicleName:    createVehicleReq.VehicleName,
		VehicleImage:   createVehicleReq.VehicleImage,
		Year:           createVehicleReq.Year,
		EngineCapacity: createVehicleReq.EngineCapacity,
		TankCapacity:   createVehicleReq.TankCapacity,
		Color:          createVehicleReq.Color,
		MachineNumber:  createVehicleReq.MachineNumber,
		ChassisNumber:  createVehicleReq.ChassisNumber,
	}

	createVehicleResponse := resp.CreateVehicleResponseModel{
		Status:  201,
		Message: "Vehicle created successfully",
	}

	if db.Error != nil {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}
	// result := db.Create(&vehicleData)
	result := db.FirstOrCreate(&vehicleData, vehicle.VehicleModel{
		MachineNumber: createVehicleReq.MachineNumber,
		ChassisNumber: createVehicleReq.ChassisNumber,
	})

	if result.Value == nil && result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  400,
			Message: "Record found",
		})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  iduint,
		NotificationTitle:       "Add Vehicle",
		NotificationDescription: "Anda Telah Menambahkan Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusBadRequest, resp.CreateVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	c.JSON(http.StatusCreated, createVehicleResponse)
}

func EditVehicle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)

	var editVehicleRequest req.EditVehicleRequestModel

	if err := c.ShouldBindJSON(&editVehicleRequest); err != nil {
		log.Println(fmt.Sprintf("error log: %s", err))
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "data required",
		})
		return
	}
	validate := validator.New()

	if err := validate.Struct(editVehicleRequest); err != nil {
		log.Println(fmt.Sprintf("error log2: %s", err))
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
			Status:  400,
			Message: "validate error, field required",
		})
		return
	}

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
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

	editVehicleResponse := resp.EditVehicleResponseModel{
		Status:  201,
		Message: "Vehicle update successfully",
	}

	if db.Error != nil {
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}

	var vehicle vehicle.VehicleModel
	if err := db.Where("id = ?", editVehicleRequest.VehicleId).First(&vehicle).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println(fmt.Printf("Vehicle with ID %d not found.\n", editVehicleRequest.VehicleId))
			c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		} else {
			log.Println("Error while querying the database.")
			c.JSON(http.StatusInternalServerError, resp.EditVehicleResponseModel{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
	} else {
		log.Println(fmt.Printf("Vehicle with ID %d is associated with User ID %d.\n", editVehicleRequest.VehicleId, vehicle.UserId))
		if vehicle.UserId != iduint {
			log.Println(fmt.Printf("Tidak sama userid"))
			c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
				Status:  http.StatusBadRequest,
				Message: "User tidak valid dengan data kendaraan",
			})
			return
		} else {
			log.Println(fmt.Printf("userid sama"))
		}
	}

	result := db.Table("vehicle_models").Where("id = ?", editVehicleRequest.VehicleId).Where("user_id = ?", returnEmailsOrUid).Update(&vehicleDataOutput)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.EditVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	c.JSON(http.StatusCreated, editVehicleResponse)
}

func GetAllVehicleData(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)

	// var vehicleData []vehicle.VehicleModel
	var vehicleData []resp.VehicleDataModel

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModel{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	if result.Value == nil {
		log.Println(fmt.Sprintf("error log3: %s", result.Error))
		c.JSON(http.StatusNotFound, resp.GetAllVehicleDataResponseModel{
			Status:  404,
			Message: "get all vehicle data Failed",
			Data:    nil,
			// Data:    []vehicle.VehicleModel{},
		})
		return
	}

	c.JSON(http.StatusOK, resp.GetAllVehicleDataResponseModel{
		Status:  200,
		Message: "get all vehicle data success",
		Data:    &vehicleData,
	})
}

func CreateLogVehicle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)

	var createLogVehicle req.CreateLogVehicleRequestModel
	if err := c.ShouldBindJSON(&createLogVehicle); err != nil {
		log.Println(fmt.Sprintf("error log CreateLogVehicleResponse: %s", err))
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "data required",
		})
		return
	}
	validate := validator.New()

	if err := validate.Struct(&createLogVehicle); err != nil {
		log.Println(fmt.Sprintf("error log2: %s", err))
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
			Status:  400,
			Message: "validate error, field required",
		})
		return
	}

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	var checkVehicleRelateWithId vehicle.VehicleModel
	if err := db.Where("id = ?", createLogVehicle.VehicleId).First(&checkVehicleRelateWithId).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println(fmt.Printf("Vehicle with ID %d not found.\n", createLogVehicle.VehicleId))
			c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		} else {
			log.Println("Error while querying the database.")
			c.JSON(http.StatusInternalServerError, resp.CreateLogVehicleResponseModel{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
	} else {
		log.Println(fmt.Printf("Vehicle with ID %d is associated with User ID %d.\n", createLogVehicle.VehicleId, checkVehicleRelateWithId.UserId))
		if checkVehicleRelateWithId.UserId != iduint {
			log.Println(fmt.Printf("Tidak sama userid"))
			c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
				Status:  http.StatusBadRequest,
				Message: "User tidak valid dengan data kendaraan",
			})
			return
		} else {
			log.Println(fmt.Printf("userid sama"))
		}
	}

	createLogVehicleResponse := resp.CreateLogVehicleResponseModel{
		Status:  201,
		Message: "Create Log Vehicle successfully",
	}

	createLogVehicleData := vehicle.VehicleMeasurementLogModel{
		UserId:              iduint,
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

	if db.Error != nil {
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}
	inputNotifModel := notif.Notification{
		UserId:                  iduint,
		NotificationTitle:       "Add Vehicle Log",
		NotificationDescription: "Anda Telah Menambahkan Log Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusBadRequest, resp.CreateLogVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	c.JSON(http.StatusCreated, createLogVehicleResponse)
}

func GetLogVehicle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.GetLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, resp.GetLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)
	var logVehicleData []resp.GetLogVehicleDataModel

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.GetLogVehicleResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.GetLogVehicleResponseModel{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	result := db.Table("vehicle_measurement_log_models").Where("user_id = ?", returnEmailsOrUid).Find(&logVehicleData)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, resp.GetLogVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	if result.Value == nil {
		log.Println(fmt.Sprintf("error log3: %s", result.Error))
		c.JSON(http.StatusNotFound, resp.GetLogVehicleResponseModel{
			Status:  404,
			Message: "get log vehicle data Failed",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp.GetLogVehicleResponseModel{
		Status:  200,
		Message: "get log vehicle data success",
		Data:    &logVehicleData,
	})
}

func EditMeasurementLogVehicle(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == false {
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid token",
		})
		return
	}

	returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)

	var editMeasurementLogVehicle req.EditMeasurementLogVehicleRequestModel
	if err := c.ShouldBindJSON(&editMeasurementLogVehicle); err != nil {
		log.Println(fmt.Sprintf("error log EditMeasurementLogVehicleResponse: %s", err))
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  http.StatusBadRequest,
			Message: "data required",
		})
		return
	}
	validate := validator.New()

	if err := validate.Struct(&editMeasurementLogVehicle); err != nil {
		log.Println(fmt.Sprintf("error log2: %s", err))
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  400,
			Message: "validate error, field required",
		})
		return
	}

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
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
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	var checkVehicleRelateWithId vehicle.VehicleModel
	if err := db.Where("id = ?", editMeasurementLogVehicle.VehicleId).First(&checkVehicleRelateWithId).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println(fmt.Printf("Vehicle with ID %d not found.\n", editMeasurementLogVehicle.VehicleId))
			c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		} else {
			log.Println("Error while querying the database.")
			c.JSON(http.StatusInternalServerError, resp.EditMeasurementLogVehicleResponseModel{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
	} else {
		log.Println(fmt.Printf("Vehicle with ID %d is associated with User ID %d.\n", editMeasurementLogVehicle.VehicleId, checkVehicleRelateWithId.UserId))
		if checkVehicleRelateWithId.UserId != iduint {
			log.Println(fmt.Printf("Tidak sama userid"))
			c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
				Status:  http.StatusBadRequest,
				Message: "User tidak valid dengan data kendaraan",
			})
			return
		} else {
			log.Println(fmt.Printf("userid sama"))
		}
	}

	editMeasurementLogVehicleResponse := resp.EditMeasurementLogVehicleResponseModel{
		Status:  202,
		Message: "Create Log Vehicle successfully",
	}

	dataUpdated := vehicle.VehicleMeasurementLogModel{
		// UserId:              iduint,
		// VehicleId:           editMeasurementLogVehicle.VehicleId,
		MeasurementTitle:    editMeasurementLogVehicle.MeasurementTitle,
		CurrentOdo:          editMeasurementLogVehicle.CurrentOdo,
		EstimateOdoChanging: editMeasurementLogVehicle.EstimateOdoChanging,
		AmountExpenses:      editMeasurementLogVehicle.AmountExpenses,
		CheckpointDate:      editMeasurementLogVehicle.CheckpointDate,
		Notes:               editMeasurementLogVehicle.Notes,
	}

	if db.Error != nil {
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}

	result := db.Table("vehicle_measurement_log_models").Where("id = ?", editMeasurementLogVehicle.ID).Update(&dataUpdated)
	log.Println(fmt.Sprintf("RowsAffected log: %d", result.RowsAffected))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}
	inputNotifModel := notif.Notification{
		UserId:                  iduint,
		NotificationTitle:       "Edit Vehicle Log",
		NotificationDescription: "Anda Telah Mengubah Log Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusBadRequest, resp.EditMeasurementLogVehicleResponseModel{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	c.JSON(http.StatusAccepted, editMeasurementLogVehicleResponse)
}

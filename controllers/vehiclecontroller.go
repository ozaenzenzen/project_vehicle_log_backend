package controllers

import (
	"fmt"
	"log"
	"net/http"

	req "project_vehicle_log_backend/data/vehicle/request"
	resp "project_vehicle_log_backend/data/vehicle/response"
	helper "project_vehicle_log_backend/helper"
	notif "project_vehicle_log_backend/models/notification"
	vehicle "project_vehicle_log_backend/models/vehicle"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

func CreateVehicle(c *gin.Context) {
	baseResponse := resp.CreateVehicleResponseModel{}

	var createVehicleReq req.CreateVehicleRequestModel
	if err := c.ShouldBindJSON(&createVehicleReq); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Create Vehicle Failed"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	validate := validator.New()
	if err := validate.Struct(createVehicleReq); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "validate error, field required"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	db, userStamp, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	vehicleData := vehicle.VehicleModel{
		UserId:         userData.ID,
		UserStamp:      *userStamp,
		VehicleName:    createVehicleReq.VehicleName,
		VehicleImage:   createVehicleReq.VehicleImage,
		Year:           createVehicleReq.Year,
		EngineCapacity: createVehicleReq.EngineCapacity,
		TankCapacity:   createVehicleReq.TankCapacity,
		Color:          createVehicleReq.Color,
		MachineNumber:  createVehicleReq.MachineNumber,
		ChassisNumber:  createVehicleReq.ChassisNumber,
	}

	// result := db.Create(&vehicleData)
	result := db.Where("user_stamp = ?", *userStamp).
		FirstOrCreate(&vehicleData, vehicle.VehicleModel{
			MachineNumber: createVehicleReq.MachineNumber,
			ChassisNumber: createVehicleReq.ChassisNumber,
		})
	if result.Value == nil && result.RowsAffected == 0 {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Record found"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	if result.Error != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = result.Error.Error()
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  userData.ID,
		UserStamp:               *userStamp,
		NotificationTitle:       "Add Vehicle",
		NotificationDescription: "Anda Telah Menambahkan Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = result.Error.Error() + "Notif error"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = 201
	baseResponse.Message = "Vehicle created successfully"
	c.JSON(http.StatusCreated, baseResponse)
}

func EditVehicle(c *gin.Context) {
	baseResponse := resp.EditVehicleResponseModel{}

	var editVehicleRequest req.EditVehicleRequestModel
	if err := c.ShouldBindJSON(&editVehicleRequest); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "data required"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	validate := validator.New()
	if err := validate.Struct(editVehicleRequest); err != nil {
		baseResponse.Status = 400
		baseResponse.Message = "validate error, field required"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	db, userStamp, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	vehicleDataOutput := vehicle.VehicleModel{
		UserId:         userData.ID,
		UserStamp:      *userStamp,
		VehicleName:    editVehicleRequest.VehicleName,
		VehicleImage:   editVehicleRequest.VehicleImage,
		Year:           editVehicleRequest.Year,
		EngineCapacity: editVehicleRequest.EngineCapacity,
		TankCapacity:   editVehicleRequest.TankCapacity,
		Color:          editVehicleRequest.Color,
		MachineNumber:  editVehicleRequest.MachineNumber,
		ChassisNumber:  editVehicleRequest.ChassisNumber,
	}

	var vehicle vehicle.VehicleModel
	if err := db.Where("id = ?", editVehicleRequest.VehicleId).First(&vehicle).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println(fmt.Printf("Vehicle with ID %d not found.\n", editVehicleRequest.VehicleId))
			baseResponse.Status = http.StatusBadRequest
			baseResponse.Message = err.Error()
			c.JSON(http.StatusBadRequest, baseResponse)
			return
		} else {
			log.Println("Error while querying the database.")
			baseResponse.Status = http.StatusInternalServerError
			baseResponse.Message = err.Error()
			c.JSON(http.StatusInternalServerError, baseResponse)
			return
		}
	} else {
		log.Println(fmt.Printf("Vehicle with ID %d is associated with User ID %d.\n", editVehicleRequest.VehicleId, vehicle.UserId))
		if vehicle.UserId != userData.ID {
			log.Println(fmt.Printf("Tidak sama userid"))
			baseResponse.Status = http.StatusBadRequest
			baseResponse.Message = "User tidak valid dengan data kendaraan"
			c.JSON(http.StatusBadRequest, baseResponse)
			return
		} else {
			log.Println(fmt.Printf("userid sama"))
		}
	}

	result := db.Table("vehicle_models").
		Where("id = ?", editVehicleRequest.VehicleId).Where("user_id = ?", userData.ID).
		Update(&vehicleDataOutput)
	if result.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error()
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  userData.ID,
		UserStamp:               *userStamp,
		NotificationTitle:       "Edit Vehicle",
		NotificationDescription: "Anda Telah Mengubah Data Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error() + "Notif error"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = 201
	baseResponse.Message = "Vehicle update successfully"
	c.JSON(http.StatusCreated, baseResponse)
}

func GetAllVehicleData(c *gin.Context) {
	baseResponse := resp.GetAllVehicleDataResponseModel{}

	db, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	var vehicles []resp.GetAllVehicleDataModel

	// Fetch all vehicles first
	if err := db.Table("vehicle_models").
		Where("user_id = ?", userData.ID).
		Find(&vehicles).Error; err != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = "Failed Fetching1"
		baseResponse.Data = nil
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}

	// For each vehicle, fetch its measurement logs manually
	for i := range vehicles {
		var measurementLogs []resp.GetAllVehicleMeasurementLogDataModel
		if err := db.Table("vehicle_measurement_log_models").
			Where("vehicle_id = ?", vehicles[i].Id).
			Find(&measurementLogs).Error; err != nil {
			baseResponse.Status = http.StatusInternalServerError
			baseResponse.Message = "Failed Fetching2"
			baseResponse.Data = nil
			c.JSON(http.StatusInternalServerError, baseResponse)
			return
		}
		vehicles[i].VehicleMeasurementLogModels = &measurementLogs
	}

	baseResponse.Status = 200
	baseResponse.Message = "get all vehicle data success"
	baseResponse.Data = &vehicles
	c.JSON(http.StatusOK, baseResponse)
}

// func GetAllVehicleDataOld(c *gin.Context) {
// 	baseResponse := resp.GetAllVehicleDataResponseModel{}

// 	db, _, userData, errorResp := helper.CustomValidatorAC(c)
// 	if errorResp != nil {
// 		baseResponse.Status = errorResp.Status
// 		baseResponse.Message = errorResp.Message
// 		baseResponse.Data = nil
// 		c.JSON(errorResp.Status, baseResponse)
// 		return
// 	}

// 	//--------check id--------check id--------check id--------
// 	// var vehicles []vehicle.VehicleModel
// 	var vehicleData []resp.GetAllVehicleDataModel

// 	result := db.
// 		Preload("VehicleMeasurementLogModels", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("id, user_id, vehicle_id, measurement_title, current_odo, estimate_odo_changing, amount_expenses, checkpoint_date, notes, created_at, updated_at")
// 			// return db.Order("created_at DESC")
// 		}).
// 		Table("vehicle_models").
// 		Where("user_id = ?", userData.ID).
// 		Find(&vehicleData)

// 	if result.Error != nil {
// 		baseResponse.Status = 400
// 		baseResponse.Message = result.Error.Error()
// 		baseResponse.Data = nil
// 		c.JSON(http.StatusBadRequest, baseResponse)
// 		return
// 	}

// 	if result.Value == nil {
// 		fmt.Println("error log3: ", result.Error)
// 		baseResponse.Status = 404
// 		baseResponse.Message = "get all vehicle data Failed"
// 		baseResponse.Data = nil
// 		c.JSON(http.StatusNotFound, baseResponse)
// 		return
// 	}

// 	// Transform data into the desired response structure
// 	var vehicleDataList []resp.GetAllVehicleDataModel
// 	for _, vehicle := range vehicleDataList {
// 		var measurementLogs []resp.GetAllVehicleMeasurementLogDataModel
// 		for _, log := range *vehicle.VehicleMeasurementLogModels {
// 			measurementLogs = append(measurementLogs, resp.GetAllVehicleMeasurementLogDataModel{
// 				Id:                  log.Id,
// 				UserId:              log.UserId,
// 				VehicleId:           log.VehicleId,
// 				MeasurementTitle:    log.MeasurementTitle,
// 				CurrentOdo:          log.CurrentOdo,
// 				EstimateOdoChanging: log.EstimateOdoChanging,
// 				AmountExpenses:      log.AmountExpenses,
// 				CheckpointDate:      log.CheckpointDate,
// 				Notes:               log.Notes,
// 				CreatedAt:           log.CreatedAt,
// 				UpdatedAt:           log.UpdatedAt,
// 			})
// 		}

// 		vehicleDataList = append(vehicleDataList, resp.GetAllVehicleDataModel{
// 			Id:                          vehicle.Id,
// 			UserId:                      vehicle.UserId,
// 			VehicleName:                 vehicle.VehicleName,
// 			VehicleImage:                vehicle.VehicleImage,
// 			Year:                        vehicle.Year,
// 			EngineCapacity:              vehicle.EngineCapacity,
// 			TankCapacity:                vehicle.TankCapacity,
// 			Color:                       vehicle.Color,
// 			MachineNumber:               vehicle.MachineNumber,
// 			ChassisNumber:               vehicle.ChassisNumber,
// 			VehicleMeasurementLogModels: &measurementLogs,
// 		})
// 	}

// 	baseResponse.Status = 200
// 	baseResponse.Message = "get all vehicle data success"
// 	baseResponse.Data = &vehicleDataList
// 	c.JSON(http.StatusOK, baseResponse)
// }

func CreateLogVehicle(c *gin.Context) {
	baseResponse := resp.CreateLogVehicleResponseModel{}

	var createLogVehicle req.CreateLogVehicleRequestModel
	if err := c.ShouldBindJSON(&createLogVehicle); err != nil {
		fmt.Println("error log CreateLogVehicleResponse: ", err)
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "data required"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&createLogVehicle); err != nil {
		fmt.Println("error log2: ", err)
		baseResponse.Status = 400
		baseResponse.Message = "validate error, field required"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	db, userStamp, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	var checkVehicleRelateWithId vehicle.VehicleModel
	if err := db.Where("id = ?", createLogVehicle.VehicleId).First(&checkVehicleRelateWithId).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println(fmt.Printf("Vehicle with ID %d not found.\n", createLogVehicle.VehicleId))
			baseResponse.Status = http.StatusBadRequest
			baseResponse.Message = err.Error()
			c.JSON(http.StatusBadRequest, baseResponse)
			return
		} else {
			log.Println("Error while querying the database.")
			baseResponse.Status = http.StatusInternalServerError
			baseResponse.Message = err.Error()
			c.JSON(http.StatusInternalServerError, baseResponse)
			return
		}
	} else {
		log.Println(fmt.Printf("Vehicle with ID %d is associated with User ID %d.\n", createLogVehicle.VehicleId, checkVehicleRelateWithId.UserId))
		if checkVehicleRelateWithId.UserId != userData.ID {
			log.Println(fmt.Printf("Tidak sama userid"))
			baseResponse.Status = http.StatusBadRequest
			baseResponse.Message = "User tidak valid dengan data kendaraan"
			c.JSON(http.StatusBadRequest, baseResponse)
			return
		} else {
			log.Println(fmt.Printf("userid sama"))
		}
	}

	createLogVehicleData := vehicle.VehicleMeasurementLogModel{
		UserId:              userData.ID,
		UserStamp:           *userStamp,
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
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error()
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  userData.ID,
		UserStamp:               *userStamp,
		NotificationTitle:       "Add Vehicle Log",
		NotificationDescription: "Anda Telah Menambahkan Log Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error() + "Notif error"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = 201
	baseResponse.Message = "Create Log Vehicle successfully"
	c.JSON(http.StatusCreated, baseResponse)
}

func GetLogVehicle(c *gin.Context) {
	baseResponse := resp.GetLogVehicleResponseModel{}

	db, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	var logVehicleData []resp.GetLogVehicleDataModel
	result := db.Table("vehicle_measurement_log_models").
		Where("user_id = ?", userData.ID).
		Find(&logVehicleData)
	if result.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	if result.Value == nil {
		fmt.Println("error log3: ", result.Error)
		baseResponse.Status = 404
		baseResponse.Message = "get log vehicle data Failed"
		baseResponse.Data = nil
		c.JSON(http.StatusNotFound, baseResponse)
		return
	}

	baseResponse.Status = 200
	baseResponse.Message = "get log vehicle data success"
	baseResponse.Data = &logVehicleData
	c.JSON(http.StatusOK, baseResponse)
}

func EditMeasurementLogVehicle(c *gin.Context) {
	baseResponse := resp.EditMeasurementLogVehicleResponseModel{}

	var editMeasurementLogVehicle req.EditMeasurementLogVehicleRequestModel
	if err := c.ShouldBindJSON(&editMeasurementLogVehicle); err != nil {
		fmt.Println("error log EditMeasurementLogVehicleResponse: ", err)
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "data required"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&editMeasurementLogVehicle); err != nil {
		fmt.Println("error log2: ", err)
		baseResponse.Status = 400
		baseResponse.Message = "validate error, field required"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	db, userStamp, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	var checkVehicleRelateWithId vehicle.VehicleModel
	if err := db.Where("id = ?", editMeasurementLogVehicle.VehicleId).First(&checkVehicleRelateWithId).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Println(fmt.Printf("Vehicle with ID %d not found.\n", editMeasurementLogVehicle.VehicleId))
			baseResponse.Status = http.StatusBadRequest
			baseResponse.Message = err.Error()
			c.JSON(http.StatusBadRequest, baseResponse)
			return
		} else {
			log.Println("Error while querying the database.")
			baseResponse.Status = http.StatusInternalServerError
			baseResponse.Message = err.Error()
			c.JSON(http.StatusInternalServerError, baseResponse)
			return
		}
	} else {
		log.Println(fmt.Printf("Vehicle with ID %d is associated with User ID %d.\n", editMeasurementLogVehicle.VehicleId, checkVehicleRelateWithId.UserId))
		if checkVehicleRelateWithId.UserId != userData.ID {
			log.Println(fmt.Printf("Tidak sama userid"))
			baseResponse.Status = http.StatusBadRequest
			baseResponse.Message = "User tidak valid dengan data kendaraan"
			c.JSON(http.StatusBadRequest, baseResponse)
			return
		} else {
			log.Println(fmt.Printf("userid sama"))
		}
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

	result := db.Table("vehicle_measurement_log_models").
		Where("id = ?", editMeasurementLogVehicle.ID).
		Update(&dataUpdated)
	fmt.Println("RowsAffected log: ", result.RowsAffected)
	if result.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error()
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  userData.ID,
		UserStamp:               *userStamp,
		NotificationTitle:       "Edit Vehicle Log",
		NotificationDescription: "Anda Telah Mengubah Log Kendaraan",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error() + "Notif error"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = 202
	baseResponse.Message = "Edit Measurement Log Vehicle successfully"
	c.JSON(http.StatusAccepted, baseResponse)
}

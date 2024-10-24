package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	req "project_vehicle_log_backend/data/vehicle/request"
	resp "project_vehicle_log_backend/data/vehicle/response"
	helper "project_vehicle_log_backend/helper"
	vehicle "project_vehicle_log_backend/models/vehicle"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func CreateVehicle(c *gin.Context) {
	baseResponse := resp.CreateVehicleResponseModel{}

	var reqBody req.CreateVehicleRequestModel

	createVehicleReq, db, userStamp, userData, errorResp := helper.CustomValidatorWithRequestBody(c, reqBody)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	vehicleStamp := uuid.New().String()
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
		VehicleStamp:   vehicleStamp,
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

	respNotif := helper.InsertNotification(
		c,
		db,
		userData,
		"Add Vehicle",
		"Anda Telah Menambahkan Kendaraan",
	)
	if respNotif != nil {
		baseResponse.Status = respNotif.Status
		baseResponse.Message = respNotif.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	baseResponse.Status = 201
	baseResponse.Message = "Vehicle created successfully"
	c.JSON(http.StatusCreated, baseResponse)
}

func EditVehicle(c *gin.Context) {
	baseResponse := resp.EditVehicleResponseModel{}

	var reqBody req.EditVehicleRequestModel

	editVehicleRequest, db, userStamp, userData, errorResp := helper.CustomValidatorWithRequestBody(c, reqBody)
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

	respNotif := helper.InsertNotification(
		c,
		db,
		userData,
		"Edit Vehicle",
		"Anda Telah Mengubah Data Kendaraan",
	)
	if respNotif != nil {
		baseResponse.Status = respNotif.Status
		baseResponse.Message = respNotif.Message
		c.JSON(errorResp.Status, baseResponse)
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

func GetAllVehicleDataV3(c *gin.Context) {
	baseResponse := resp.GetAllVehicleDataResponseModelV3{}

	db, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	var vehicles []resp.GetAllVehicleDataModelV3

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
		var measurementLogs []resp.GetAllVehicleMeasurementLogDataModelV3
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

		// QUERY RAW
		query := "SELECT measurement_title FROM vehicle_measurement_log_models WHERE vehicle_id = ? GROUP BY measurement_title"
		measurementType := []string{}
		if errQuery := db.Raw(query, vehicles[i].Id).Pluck("measurement_title", &measurementType).Error; errQuery != nil {
			baseResponse.Status = http.StatusBadRequest
			baseResponse.Message = errQuery.Error()
			baseResponse.Data = nil
			c.JSON(http.StatusBadRequest, baseResponse)
			return
		}
		vehicles[i].MeasurementTitle = &measurementType
	}

	baseResponse.Status = 200
	baseResponse.Message = "get all vehicle data success"
	baseResponse.Data = &vehicles
	c.JSON(http.StatusOK, baseResponse)
}

func GetAllVehicleDataV2(c *gin.Context) {
	baseResponse := resp.GetAllVehicleDataResponseModelV2{}

	var reqData req.GetAllVehicleDataRequestModelV2
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModelV2{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap1",
			Data:    nil,
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqData); err != nil {
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModelV2{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap2",
			Data:    nil,
		})
		return
	}

	db, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	// resultData, errPagination := GetAllVehiclePagination(
	// resultData, errPagination := GetAllVehiclePaginationUsingRaw(
	resultData, errPagination := GetAllVehiclePaginationUsingRawV2(
		db,
		reqData.CurrentPage,
		reqData.Limit,
		reqData.SortOrder,
		&userData.UserStamp,
	)
	if errPagination != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = errPagination.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = 200
	baseResponse.Message = "get all vehicle data success"
	baseResponse.Data = resultData
	c.JSON(http.StatusOK, baseResponse)
}

func GetAllVehiclePagination(
	db *gorm.DB,
	currentPage int,
	limit int,
	sortOrder *string,
	userStamp *string,
) (*resp.GetAllVehicleDataModelV2, error) {
	// Handle limit with a maximum of 20
	if limit > 15 {
		limit = 15
	} else if limit <= 0 {
		limit = 10
	}

	// Ensure currentPage is at least 1
	if currentPage <= 0 {
		currentPage = 1
	}

	offset := (currentPage - 1) * limit
	var totalItems int64
	var vehicles []vehicle.VehicleModel
	var result []resp.DataGetAllVehicleV2

	// Count total items
	query := db.Model(&vehicle.VehicleModel{}).
		Where("vehicle_models.user_stamp = ?", userStamp)

	if sortOrder != nil && *sortOrder != "" {
		query = query.Order("created_at " + *sortOrder)
	} else {
		query = query.Order("created_at desc")
	}

	query.Count(&totalItems)

	// Pagination calculation
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}

	// If currentPage exceeds totalPages, return empty list
	if currentPage > totalPages {
		return &resp.GetAllVehicleDataModelV2{
			CurrentPage: currentPage,
			NextPage:    0,
			TotalPages:  totalPages,
			TotalItems:  int(totalItems),
			Data:        &[]resp.DataGetAllVehicleV2{},
		}, nil
	}

	// Fetch monitor data with pagination
	query = query.Limit(limit).Offset(offset).Find(&vehicles)
	if query.Error != nil {
		return nil, query.Error
	}

	// Fetch corresponding user data for each monitor and prepare custom response
	for _, vehicle := range vehicles {

		measurementType := []string{}
		// QUERY RAW
		query := "SELECT measurement_title FROM vehicle_measurement_log_models WHERE vehicle_id = ? GROUP BY measurement_title"
		if errQuery := db.Raw(query, vehicle.Id).Pluck("measurement_title", &measurementType).Error; errQuery != nil {
			return nil, errQuery
		}

		dataVehicle := resp.DataGetAllVehicleV2{
			Id:               vehicle.Id,
			UserId:           vehicle.UserId,
			UserStamp:        vehicle.UserStamp,
			VehicleName:      vehicle.VehicleName,
			VehicleImage:     vehicle.VehicleImage,
			Year:             vehicle.Year,
			EngineCapacity:   vehicle.EngineCapacity,
			TankCapacity:     vehicle.TankCapacity,
			Color:            vehicle.Color,
			MachineNumber:    vehicle.MachineNumber,
			ChassisNumber:    vehicle.ChassisNumber,
			MeasurementTitle: &measurementType,
		}

		result = append(result, dataVehicle)
	}

	// Calculate nextPage
	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	response2 := resp.GetAllVehicleDataModelV2{
		CurrentPage: currentPage,
		NextPage:    currentPage + 1,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
		Data:        &result,
	}

	return &response2, nil
}

func GetAllVehiclePaginationUsingRaw(
	db *gorm.DB,
	currentPage int,
	limit int,
	sortOrder *string,
	userStamp *string,
) (*resp.GetAllVehicleDataModelV2, error) {
	// Handle limit with a maximum of 20
	if limit > 15 {
		limit = 15
	} else if limit <= 0 {
		limit = 10
	}

	// Ensure currentPage is at least 1
	if currentPage <= 0 {
		currentPage = 1
	}

	offset := (currentPage - 1) * limit
	var totalItems int64
	// var vehicles []vehicle.VehicleModel
	var result []resp.DataGetAllVehicleV2

	// Count total items
	query := db.Model(&vehicle.VehicleModel{}).
		Where("vehicle_models.user_stamp = ?", userStamp)

	if sortOrder != nil && *sortOrder != "" {
		query = query.Order("created_at " + *sortOrder)
	} else {
		query = query.Order("created_at desc")
	}

	query.Count(&totalItems)

	// Pagination calculation
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}

	// If currentPage exceeds totalPages, return empty list
	if currentPage > totalPages {
		return &resp.GetAllVehicleDataModelV2{
			CurrentPage: currentPage,
			NextPage:    0,
			TotalPages:  totalPages,
			TotalItems:  int(totalItems),
			Data:        &[]resp.DataGetAllVehicleV2{},
		}, nil
	}

	// Fetch monitor data with pagination
	// query = query.Limit(limit).Offset(offset).Find(&vehicles)
	// Query vehicle data and join with measurement data
	query1 := `
		SELECT v.id, v.user_id, v.user_stamp, v.vehicle_name, v.vehicle_image, v.year, v.engine_capacity, v.tank_capacity, 
			   v.color, v.machine_number, v.chassis_number, GROUP_CONCAT(m.measurement_title) as measurment_title
		FROM vehicle_models v
		LEFT JOIN vehicle_measurement_log_models m ON v.id = m.vehicle_id
		GROUP BY v.id
		LIMIT ? OFFSET ?
		`

	// Query raw SQL
	db.Raw(query1, limit, offset).Scan(&result)
	if query.Error != nil {
		return nil, query.Error
	}

	// Fetch corresponding user data for each monitor and prepare custom response
	// for _, vehicle := range vehicles {

	// 	measurementType := []string{}
	// 	// QUERY RAW
	// 	query := "SELECT measurement_title FROM vehicle_measurement_log_models WHERE vehicle_id = ? GROUP BY measurement_title"
	// 	if errQuery := db.Raw(query, vehicle.Id).Pluck("measurement_title", &measurementType).Error; errQuery != nil {
	// 		return nil, errQuery
	// 	}

	// 	dataVehicle := resp.DataGetAllVehicleV2{
	// 		Id:               vehicle.Id,
	// 		UserId:           vehicle.UserId,
	// 		UserStamp:        vehicle.UserStamp,
	// 		VehicleName:      vehicle.VehicleName,
	// 		VehicleImage:     vehicle.VehicleImage,
	// 		Year:             vehicle.Year,
	// 		EngineCapacity:   vehicle.EngineCapacity,
	// 		TankCapacity:     vehicle.TankCapacity,
	// 		Color:            vehicle.Color,
	// 		MachineNumber:    vehicle.MachineNumber,
	// 		ChassisNumber:    vehicle.ChassisNumber,
	// 		MeasurementTitle: &measurementType,
	// 	}

	// 	result = append(result, dataVehicle)
	// }

	// Calculate nextPage
	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	response2 := resp.GetAllVehicleDataModelV2{
		CurrentPage: currentPage,
		NextPage:    currentPage + 1,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
		Data:        &result,
	}

	return &response2, nil
}

func GetAllVehiclePaginationUsingRawV2(
	db *gorm.DB,
	currentPage int,
	limit int,
	sortOrder *string,
	userStamp *string,
) (*resp.GetAllVehicleDataModelV2, error) {
	// Handle limit with a maximum of 15
	if limit > 15 {
		limit = 15
	} else if limit <= 0 {
		limit = 10
	}

	// Ensure currentPage is at least 1
	if currentPage <= 0 {
		currentPage = 1
	}

	offset := (currentPage - 1) * limit
	var totalItems int64
	var result []resp.DataGetAllVehicleV2

	// Query to count the total number of items
	query := db.Model(&vehicle.VehicleModel{}).
		Where("vehicle_models.user_stamp = ?", userStamp)

	// Sorting order
	if sortOrder != nil && *sortOrder != "" {
		query = query.Order("created_at " + *sortOrder)
	} else {
		query = query.Order("created_at desc")
	}

	// Get total count
	query.Count(&totalItems)

	// Pagination calculation
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}

	// If currentPage exceeds totalPages, return empty list
	if currentPage > totalPages {
		return &resp.GetAllVehicleDataModelV2{
			CurrentPage: currentPage,
			NextPage:    0,
			TotalPages:  totalPages,
			TotalItems:  int(totalItems),
			Data:        &[]resp.DataGetAllVehicleV2{},
		}, nil
	}

	// Raw SQL query with JOIN to fetch data
	// rawQuery := `
	// SELECT
	// 	v.id,
	// 	v.user_id,
	// 	v.user_stamp,
	// 	v.vehicle_name,
	// 	v.vehicle_image,
	// 	v.year,
	// 	v.engine_capacity,
	// 	v.tank_capacity,
	// 	v.color,
	// 	v.machine_number,
	// 	v.chassis_number,
	// 	(
	//         SELECT
	//             IFNULL(JSON_ARRAYAGG(measurements.measurement_title), '[]')
	//         FROM
	//         (
	//             SELECT DISTINCT m.measurement_title
	//             FROM vehicle_measurement_log_models m
	//             WHERE m.vehicle_id = v.id
	//         ) AS measurements
	//     ) AS measurement_title
	// FROM vehicle_models v
	// LEFT JOIN vehicle_measurement_log_models m ON m.vehicle_id = v.id
	// WHERE v.user_stamp = ?
	// GROUP BY v.id
	// ORDER BY v.created_at DESC
	// LIMIT ? OFFSET ?`

	// For mariadb 10.4 and below
	rawQuery := `
			SELECT
				v.id,
				v.user_id,
				v.user_stamp,
				v.vehicle_name,
				v.vehicle_image,
				v.year,
				v.engine_capacity,
				v.tank_capacity,
				v.color,
				v.machine_number,
				v.chassis_number,
				(
			        SELECT
			            IFNULL(CONCAT('[', GROUP_CONCAT(DISTINCT CONCAT('"', m.measurement_title, '"')), ']'),'[]')
			        FROM
			        (
			            SELECT DISTINCT m.measurement_title
			            FROM vehicle_measurement_log_models m
			            WHERE m.vehicle_id = id
			        ) AS measurements
			    ) AS measurement_title
			FROM vehicle_models v
			LEFT JOIN vehicle_measurement_log_models m ON m.vehicle_id = v.id
			WHERE v.user_stamp = ?
			GROUP BY id
			ORDER BY v.created_at DESC
			LIMIT ? OFFSET ?`

	// IFNULL(JSON_ARRAYAGG(measurements.measurement_title), '[]')
	// CONCAT('[', GROUP_CONCAT(CONCAT(measurements.measurement_title)), ']')
	// Execute the query
	if err := db.Raw(rawQuery, userStamp, limit, offset).Scan(&result).Error; err != nil {
		return nil, err
	}

	// Unmarshal MeasurementTitle JSON into a slice
	for i := range result {
		// Example with string
		var value2 any = result[i].MeasurementTitle
		bytes2, errConvert := convertToBytes(value2)
		if errConvert != nil {
			fmt.Println("Error:", errConvert)
			return nil, errConvert
		}

		var titles []string
		if err := json.Unmarshal([]byte(bytes2), &titles); err != nil {
			return nil, err
		}
		result[i].MeasurementTitle = titles // Set the unmarshaled result back to the struct
	}

	// Calculate nextPage
	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	// Prepare the final response
	response := resp.GetAllVehicleDataModelV2{
		CurrentPage: currentPage,
		NextPage:    nextPage,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
		Data:        &result,
	}

	return &response, nil
}

func convertToBytes(value any) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil // Directly return if it's already []byte
	case string:
		return []byte(v), nil // Convert string to []byte
	default:
		return nil, fmt.Errorf("unsupported type: %T", v) // Handle unsupported types
	}
}

func CreateLogVehicle(c *gin.Context) {
	baseResponse := resp.CreateLogVehicleResponseModel{}

	var reqBody req.CreateLogVehicleRequestModel

	createLogVehicle, db, userStamp, userData, errorResp := helper.CustomValidatorWithRequestBody(c, reqBody)
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

	logStamp := uuid.New().String()
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
		LogStamp:            logStamp,
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

	respNotif := helper.InsertNotification(
		c,
		db,
		userData,
		"Add Vehicle Log",
		"Anda Telah Menambahkan Log Kendaraan",
	)
	if respNotif != nil {
		baseResponse.Status = respNotif.Status
		baseResponse.Message = respNotif.Message
		c.JSON(errorResp.Status, baseResponse)
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

func GetListLogType(c *gin.Context) {
	baseResponse := resp.GetListLogTypeResponseModel{}

	db, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	// QUERY RAW
	var measurementType []resp.LogDataModel
	query := "SELECT measurement_title FROM vehicle_measurement_log_models WHERE user_stamp = ? GROUP BY measurement_title"
	if errQuery := db.Raw(query, userData.UserStamp).Scan(&measurementType).Error; errQuery != nil {
		c.JSON(
			http.StatusBadRequest, resp.GetListLogTypeResponseModel{
				Status:  http.StatusBadRequest,
				Message: errQuery.Error(),
				Data:    nil,
			},
		)
		return
	}
	c.JSON(http.StatusOK,
		resp.GetListLogTypeResponseModel{
			Status:  200,
			Message: "Get log type success",
			Data:    &measurementType,
		},
	)
}

func GetLogVehicleV2(c *gin.Context) {
	baseResponse := resp.GetLogVehicleDataResponseModelV2{}

	var reqData req.GetLogVehicleDataRequestModelV2
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModelV2{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap1",
			Data:    nil,
		})
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqData); err != nil {
		c.JSON(http.StatusBadRequest, resp.GetAllVehicleDataResponseModelV2{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap2",
			Data:    nil,
		})
		return
	}

	db, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	var resultDataAnalytics resp.DataAnalyticsVehicleV2
	db.Raw(`
    	SELECT 
    	    SUM(CAST(amount_expenses AS DECIMAL(10,2))) AS total_expenses, 
        	MAX(created_at) AS last_created_at,
        	IFNULL(CONCAT('[', GROUP_CONCAT(DISTINCT CONCAT('"', measurement_title, '"')), ']'),'[]') AS measurement_titles,
			AVG(CAST(amount_expenses AS DECIMAL(10,2))) AS avg_expenses_per_meas,
        	AVG(CAST(estimate_odo_changing AS DECIMAL(10,2))) AS avg_odo_change,
        	(SELECT AVG(DATEDIFF(current.checkpoint_date, previous.checkpoint_date))
				FROM vehicle_measurement_log_models AS current
				JOIN vehicle_measurement_log_models AS previous ON current.user_id = previous.user_id AND current.created_at > previous.created_at
				WHERE current.user_id = ?) AS avg_service_freq,
			(SELECT GROUP_CONCAT(measurement_title ORDER BY freq DESC SEPARATOR ', ')
				FROM (
					SELECT measurement_title, COUNT(*) AS freq
					FROM vehicle_measurement_log_models
					WHERE user_id = ?
					GROUP BY measurement_title
					ORDER BY freq DESC
					LIMIT 5
				) AS freq_titles) AS most_frequent_titles,
			(SELECT GROUP_CONCAT(CONCAT(measurement_title, ': ', total_expenses))
				FROM (
					SELECT measurement_title, SUM(CAST(amount_expenses AS DECIMAL(10,2))) AS total_expenses
					FROM vehicle_measurement_log_models
					WHERE user_id = ?
					GROUP BY measurement_title
				) AS breakdown) AS cost_breakdown
    	FROM vehicle_measurement_log_models 
    	WHERE user_id = ?
	`, userData.ID, userData.ID, userData.ID, userData.ID).Scan(&resultDataAnalytics)

	// (SELECT GROUP_CONCAT(measurement_title ORDER BY COUNT(*) DESC SEPARATOR ', ')
	// 			FROM vehicle_measurement_log_models
	// 			WHERE user_id = ?
	// 			GROUP BY measurement_title
	// 			ORDER BY COUNT(*) DESC
	// 			LIMIT 5) AS most_frequent_titles
	// (SELECT GROUP_CONCAT(CONCAT(measurement_title, ': ', SUM(CAST(amount_expenses AS DECIMAL(10,2)))))
	// 			FROM vehicle_measurement_log_models
	// 			WHERE user_id = ?
	// 			GROUP BY measurement_title) AS cost_breakdown

	// IFNULL(CONCAT('[', GROUP_CONCAT(DISTINCT CONCAT('"', m.measurement_title, '"')), ']'),'[]')
	// GROUP_CONCAT(DISTINCT measurement_title) AS measurement_titles

	// Unmarshal MeasurementTitle JSON into a slice
	var value2 any = resultDataAnalytics.MeasurementTitles
	bytes2, errConvert := convertToBytes(value2)
	if errConvert != nil {
		fmt.Println("Error:", errConvert)
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = errConvert.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	var titles []string
	if errUnmarshal := json.Unmarshal([]byte(bytes2), &titles); errUnmarshal != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = errUnmarshal.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}
	resultDataAnalytics.MeasurementTitles = titles

	resultData, errPagination := GetLogVehiclePagination(
		db,
		reqData.CurrentPage,
		reqData.Limit,
		reqData.VehicleID,
		reqData.MeasurementTitle,
		reqData.SortOrder,
		&userData.UserStamp,
	)
	if errPagination != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = errPagination.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	resultData.CollectionLogData = resultDataAnalytics

	baseResponse.Status = 200
	baseResponse.Message = "get log vehicle data success"
	baseResponse.Data = resultData
	c.JSON(http.StatusOK, baseResponse)
}

func GetLogVehiclePagination(
	db *gorm.DB,
	currentPage int,
	limit int,
	vehicleId *string,
	measurementTitle *string,
	sortOrder *string,
	userStamp *string,
) (*resp.GetLogVehicleDataModelV2, error) {
	// Handle limit with a maximum of 20
	if limit > 15 {
		limit = 15
	} else if limit <= 0 {
		limit = 10
	}

	// Ensure currentPage is at least 1
	if currentPage <= 0 {
		currentPage = 1
	}

	offset := (currentPage - 1) * limit
	var totalItems int64
	var vehicles []vehicle.VehicleMeasurementLogModel
	var result []resp.DataGetLogVehicleV2

	// Count total items
	query := db.Model(&vehicle.VehicleMeasurementLogModel{}).
		Where("vehicle_measurement_log_models.user_stamp = ?", userStamp)

	if sortOrder != nil && *sortOrder != "" {
		query = query.Order("created_at " + *sortOrder)
	} else {
		query = query.Order("created_at desc")
	}

	if vehicleId != nil && *vehicleId != "" {
		query = query.Where("vehicle_id = ?", *vehicleId)
	}

	if measurementTitle != nil && *measurementTitle != "" {
		query = query.Where("measurement_title = ?", *measurementTitle)
	}

	query.Count(&totalItems)

	// Pagination calculation
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}

	// If currentPage exceeds totalPages, return empty list
	if currentPage > totalPages {
		return &resp.GetLogVehicleDataModelV2{
			CurrentPage: currentPage,
			NextPage:    0,
			TotalPages:  totalPages,
			TotalItems:  int(totalItems),
			Data:        &[]resp.DataGetLogVehicleV2{},
		}, nil
	}

	// Fetch monitor data with pagination
	query = query.Limit(limit).Offset(offset).Find(&vehicles)
	if query.Error != nil {
		return nil, query.Error
	}

	// Fetch corresponding user data for each monitor and prepare custom response
	for _, vehicle := range vehicles {

		dataVehicle := resp.DataGetLogVehicleV2{
			Id:                  vehicle.Id,
			UserId:              vehicle.UserId,
			UserStamp:           vehicle.UserStamp,
			VehicleId:           vehicle.VehicleId,
			MeasurementTitle:    vehicle.MeasurementTitle,
			CurrentOdo:          vehicle.CurrentOdo,
			EstimateOdoChanging: vehicle.EstimateOdoChanging,
			AmountExpenses:      vehicle.AmountExpenses,
			CheckpointDate:      vehicle.CheckpointDate,
			Notes:               vehicle.Notes,
			CreatedAt:           vehicle.CreatedAt,
			UpdatedAt:           vehicle.UpdatedAt,
		}

		result = append(result, dataVehicle)
	}

	// Calculate nextPage
	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	response2 := resp.GetLogVehicleDataModelV2{
		CurrentPage: currentPage,
		NextPage:    currentPage + 1,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
		Data:        &result,
	}

	return &response2, nil
}

func EditMeasurementLogVehicle(c *gin.Context) {
	baseResponse := resp.EditMeasurementLogVehicleResponseModel{}

	var reqBody req.EditMeasurementLogVehicleRequestModel

	editMeasurementLogVehicle, db, _, userData, errorResp := helper.CustomValidatorWithRequestBody(c, reqBody)
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

	respNotif := helper.InsertNotification(
		c,
		db,
		userData,
		"Edit Vehicle Log",
		"Anda Telah Mengubah Log Kendaraan",
	)
	if respNotif != nil {
		baseResponse.Status = respNotif.Status
		baseResponse.Message = respNotif.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	baseResponse.Status = 202
	baseResponse.Message = "Edit Measurement Log Vehicle successfully"
	c.JSON(http.StatusAccepted, baseResponse)
}

package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	req "project_vehicle_log_backend/data/notification/request"
	resp "project_vehicle_log_backend/data/notification/response"
	helper "project_vehicle_log_backend/helper"
	account "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type GetNotificationResponse struct {
	Status       int                   `json:"status"`
	Message      string                `json:"message"`
	Notification *[]notif.Notification `json:"notification"`
}

func GetNotificationByUserId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	headertoken := c.Request.Header.Get("token")
	isValid, err := helper.ValidateTokenJWT(c, db, headertoken)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetNotificationResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid {

		var notificationList []notif.Notification

		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, GetNotificationResponse{
				Status:  http.StatusBadRequest,
				Message: "Empty parameter",
			})
			return
		}
		returnEmailsOrUid := helper.GetDataTokenJWT(headertoken, false)
		if returnEmailsOrUid != c.Param("id") {
			c.JSON(http.StatusBadRequest, GetNotificationResponse{
				Status:  http.StatusBadRequest,
				Message: "Wrong token",
			})
			return
		}

		//--------check id--------check id--------check id--------

		iduint64, err := strconv.ParseUint(c.Param("id"), 10, 32)

		if err != nil {
			c.JSON(http.StatusBadRequest, GetNotificationResponse{
				Status:  500,
				Message: "error parsing",
			})
			return
		}
		iduint := uint(iduint64)

		checkID := db.Table("account_user_models").Where("id = ?", c.Param("id")).Find(&account.AccountUserModel{
			ID: iduint,
		})

		if checkID.Error != nil {
			c.JSON(http.StatusBadRequest, GetNotificationResponse{
				Status:  400,
				Message: checkID.Error.Error(),
			})
			return
		}

		//--------check id--------check id--------check id--------

		result := db.Table("notifications").Where("user_id = ?", c.Param("id")).Find(&notificationList)

		if result.Value == nil {
			fmt.Println("error log notification1: ", result.Error)
			c.JSON(http.StatusInternalServerError, GetNotificationResponse{
				Status:       500,
				Message:      "get notification failed 1",
				Notification: &[]notif.Notification{},
			})
			return
		}

		c.JSON(http.StatusOK, GetNotificationResponse{
			Status:       200,
			Message:      "get notification success",
			Notification: &notificationList,
		})
	} else {
		c.JSON(http.StatusBadRequest, GetNotificationResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid token",
		})
		return
	}
}

func GetNotification(c *gin.Context) {
	baseResponse := resp.GetNotificationResponseModel{}

	reqData := req.GetNotificationaRequestModel{}

	getNotificationReq, db, _, userData, errorResp := helper.CustomValidatorWithRequestBody(c, reqData)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	resultData, errPagination := GetNotificationPaginationUsingRaw(
		db,
		getNotificationReq.CurrentPage,
		getNotificationReq.Limit,
		getNotificationReq.SortOrder,
		&userData.UserStamp,
	)
	if errPagination != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = errPagination.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = http.StatusOK
	baseResponse.Message = "Get Notification Successfully"
	baseResponse.Data = resultData
	c.JSON(http.StatusOK, baseResponse)
}

func GetNotificationPaginationUsingRaw(
	db *gorm.DB,
	currentPage int,
	limit int,
	sortOrder *string,
	userStamp *string,
) (*resp.GetNotificationPagination, error) {
	sortOrder = new(string)
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
	var result []resp.GetNotificationDataModel

	// Query to count the total number of items
	query := db.Model(&notif.Notification{}).
		Where("notifications.user_stamp = ?", userStamp)

	// Get total count
	query.Count(&totalItems)

	// Pagination calculation
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit != 0 {
		totalPages++
	}

	// If currentPage exceeds totalPages, return empty list
	if currentPage > totalPages {
		return &resp.GetNotificationPagination{
			CurrentPage: currentPage,
			NextPage:    0,
			TotalPages:  totalPages,
			TotalItems:  int(totalItems),
			Data:        &[]resp.GetNotificationDataModel{},
		}, nil
	}

	// Sorting order
	if sortOrder != nil && *sortOrder != "" {
		// query = query.Order("created_at " + *sortOrder)
		if *sortOrder != "ASC" && *sortOrder != "DESC" {
			*sortOrder = "DESC" // Default to DESC if invalid
		}
	} else {
		*sortOrder = "DESC"
	}

	// For mariadb 10.4 and below
	// rawQuery := `
	// 		SELECT
	// 			v.notification_id,
	// 			v.user_id,
	// 			v.user_stamp,
	// 			v.notification_title,
	// 			v.notification_description,
	// 			v.notification_status,
	// 			v.notification_type,
	// 			v.notification_stamp,
	// 			v.created_at,
	// 			v.updated_at
	// 		FROM notifications v
	// 		WHERE v.user_stamp = ?
	// 		ORDER BY v.created_at DESC
	// 		LIMIT ? OFFSET ?`

	rawQuery := fmt.Sprintf(`
			SELECT
				v.notification_id,
				v.user_id,
				v.user_stamp,
				v.notification_title,
				v.notification_description,
				v.notification_status,
				v.notification_type,
				v.notification_stamp,
				v.created_at,
				v.updated_at
			FROM notifications v
			WHERE v.user_stamp = ?
			ORDER BY v.created_at %s
			LIMIT ? OFFSET ?
			`, *sortOrder)

	if err := db.Raw(rawQuery, userStamp, limit, offset).Scan(&result).Error; err != nil {
		return nil, err
	}

	// // Unmarshal MeasurementTitle JSON into a slice
	// for i := range result {
	// 	// Example with string
	// 	var value2 any = result[i].MeasurementTitle
	// 	bytes2, errConvert := convertToBytes(value2)
	// 	if errConvert != nil {
	// 		fmt.Println("Error:", errConvert)
	// 		return nil, errConvert
	// 	}

	// 	var titles []string
	// 	if err := json.Unmarshal([]byte(bytes2), &titles); err != nil {
	// 		return nil, err
	// 	}
	// 	result[i].MeasurementTitle = titles // Set the unmarshaled result back to the struct
	// }

	// Calculate nextPage
	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	// Prepare the final response
	response := resp.GetNotificationPagination{
		CurrentPage: currentPage,
		NextPage:    nextPage,
		TotalPages:  totalPages,
		TotalItems:  int(totalItems),
		Data:        &result,
	}

	return &response, nil
}

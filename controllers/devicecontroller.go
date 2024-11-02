package controllers

import (
	"net/http"
	req "project_vehicle_log_backend/data/device/request"
	resp "project_vehicle_log_backend/data/device/response"
	helper "project_vehicle_log_backend/helper"
	device "project_vehicle_log_backend/models/device"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckDevice(c *gin.Context) {
	baseResponse := resp.CheckDeviceResponseModel{}

	var reqBody req.CheckDeviceRequestModel

	checkDeviceRequest, db, userStamp, _, errorResp := helper.CustomValidatorWithRequestBody(c, reqBody)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	deviceStamp := uuid.New().String()
	checkDeviceData := device.DeviceModel{
		DeviceStamp: deviceStamp,
		DeviceID:    checkDeviceRequest.DeviceID,
		UserStamp:   *userStamp,
		DeviceName:  checkDeviceRequest.DeviceName,
		DeviceType:  checkDeviceRequest.DeviceType,
		OSVersion:   checkDeviceRequest.OSVersion,
		AppVersion:  checkDeviceRequest.AppVersion,
		PushToken:   checkDeviceRequest.PushToken,
		LastActive:  checkDeviceRequest.LastActive,
	}

	result := db.Create(&checkDeviceData)
	if result.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = result.Error.Error()
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = 201
	baseResponse.Message = "Check Device Successfully"
	baseResponse.Data = &resp.CheckDeviceDataModel{
		DeviceID:    checkDeviceData.DeviceID,
		DeviceStamp: checkDeviceData.DeviceStamp,
		UserStamp:   checkDeviceData.UserStamp,
		DeviceName:  checkDeviceData.DeviceName,
		DeviceType:  checkDeviceData.DeviceType,
		OSVersion:   checkDeviceData.OSVersion,
		AppVersion:  checkDeviceData.AppVersion,
		PushToken:   checkDeviceData.PushToken,
		LastActive:  checkDeviceData.LastActive,
		CreatedAt:   checkDeviceData.CreatedAt,
		UpdatedAt:   checkDeviceData.UpdatedAt,
	}
	c.JSON(http.StatusCreated, baseResponse)
}

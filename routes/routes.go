package routes

import (
	"net/http"
	"project_vehicle_log_backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type HandleRoutesResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	r.POST("/account/signin", controllers.SignInAccount)
	r.POST("/account/signup", controllers.SignUpAccount)
	r.POST("/account/editprofile", controllers.EditProfile)
	// r.GET("/account/userdata/:id", controllers.GetUserData)
	r.GET("/account/userdata", controllers.GetUserData)

	r.GET("/vehicle/allvehicle", controllers.GetAllVehicleData)
	r.POST("/vehicle/createvehicle", controllers.CreateVehicle)
	r.POST("/vehicle/editvehicle", controllers.EditVehicle)
	r.POST("/vehicle/createlogvehicle", controllers.CreateLogVehicle)
	r.GET("/vehicle/getlogvehicle", controllers.GetLogVehicle)
	r.PUT("/vehicle/editmeasurementlogvehicle", controllers.EditMeasurementLogVehicle)
	// r.PUT("/vehicle/deletemeasurementlogvehicle", controllers.GetLogVehicle) //next update
	// r.POST("/event/submitevent", controllers.SubmitEvent)
	// r.GET("/event/getlistvolunteers/:id", controllers.GetListVolunteer) //event id

	r.GET("/notifications/:id", controllers.GetNotificationByUserId) //user id

	r.NoRoute(func(c *gin.Context) {
		// response := map[string]interface{}{
		// 	"status":  404,
		// 	"message": "Page not found",
		// }
		// c.JSON(http.StatusNotFound, response)
		// c.JSON(http.StatusNotFound, gin.H{
		// 	"status":  http.StatusNotFound,
		// 	"message": "Page not found",
		// })
		c.JSON(http.StatusNotFound, HandleRoutesResponse{
			Status:  http.StatusNotFound,
			Message: "Page not found",
		})

		return
	})

	return r
}

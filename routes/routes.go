package routes

import (
	"project_vehicle_log_backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	r.POST("/account/signin", controllers.SignInAccount)
	r.POST("/account/signup", controllers.SignUpAccount)
	r.GET("/account/userdata/:id", controllers.GetUserData)
	r.GET("/vehicle/allvehicle", controllers.GetAllVehicleData)
	r.POST("/vehicle/createvehicle", controllers.CreateVehicle)
	r.POST("/vehicle/editvehicle", controllers.EditVehicle)
	r.POST("/vehicle/createlogvehicle", controllers.CreateLogVehicle)
	r.GET("/vehicle/getlogvehicle", controllers.GetLogVehicle)
	// r.POST("/event/submitevent", controllers.SubmitEvent)
	// r.GET("/event/getlistvolunteers/:id", controllers.GetListVolunteer) //event id
	r.GET("/notifications/:id", controllers.GetNotificationByUserId) //user id

	return r
}

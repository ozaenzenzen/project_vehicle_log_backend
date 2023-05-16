package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	// r.POST("/account/signin", controllers.SignInAccount)
	// r.POST("/account/signup", controllers.SignUpAccountV2)
	// r.GET("/account/userdata/:id", controllers.GetUserData)
	// r.GET("/event/allevent", controllers.GetAllEvent)
	// r.POST("/event/createevent", controllers.CreateEvent)
	// r.POST("/event/joinevent", controllers.JoinEvent)
	// r.POST("/event/submitvolunteer", controllers.SubmitVolunteer)
	// r.POST("/event/submitevent", controllers.SubmitEvent)
	// r.GET("/event/getlistvolunteers/:id", controllers.GetListVolunteer) //event id
	// r.GET("/notifications/:id", controllers.GetNotificationByUserId)    //user id

	return r
}

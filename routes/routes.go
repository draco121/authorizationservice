package routes

import (
	"github.com/draco121/authorizationservice/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(controllers controllers.Controllers, router *gin.Engine) {
	v1 := router.Group("/v1")
	v1.POST("/authorize", controllers.Authorize)
}

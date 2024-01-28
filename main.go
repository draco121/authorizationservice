package main

import (
	"github.com/draco121/authorizationservice/controllers"
	"github.com/draco121/authorizationservice/core"
	"github.com/draco121/authorizationservice/repository"
	"github.com/draco121/authorizationservice/routes"
	"github.com/draco121/common/clients"
	"github.com/draco121/common/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func RunApp() {
	db := database.NewMongoDatabase(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_DBNAME"))
	authenticationServiceApiClient := clients.NewAuthenticationServiceApiClient(os.Getenv("AUTHENTICATION_SERVICE_BASEURL"))
	repo := repository.NewAuthorizationRepo(db, authenticationServiceApiClient)
	service := core.NewAuthorizationService(repo)
	controllers := controllers.NewControllers(service)
	router := gin.Default()
	routes.RegisterRoutes(controllers, router)
	err := router.Run()
	if err != nil {
		return
	}
}
func main() {
	_ = godotenv.Load()
	RunApp()
}

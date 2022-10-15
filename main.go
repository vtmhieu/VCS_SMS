package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/controllers"
	"github.com/vtmhieu/VCS_SMS/initializers"
	"github.com/vtmhieu/VCS_SMS/routes"
)

var (
	server *gin.Engine

	Servercontroller      controllers.Server_controller
	ServerRouteController routes.Server_Route_Controller

	Usercontroller      controllers.User_controller
	UserRouteController routes.User_Route_controller

	Authcontroller      controllers.Auth_controller
	AuthRouteController routes.Auth_Route_controller
)

// func NewOpenAPIMiddleware() gin.HandlerFunc {
// 	validator := middleware.OpenapiInputValidator("./openapi.yaml")
// 	return validator
// }

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	Servercontroller = controllers.New_server_controller(initializers.DB)
	ServerRouteController = routes.New_route_server_controller(Servercontroller)

	Usercontroller = controllers.New_user_controller(initializers.DB)
	UserRouteController = routes.New_user_route_controller(Usercontroller)

	Authcontroller = controllers.New_auth_controller(initializers.DB)
	AuthRouteController = routes.New_Auth_Route_controller(Authcontroller)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to VCS Server Management System"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	ServerRouteController.Server_Route(router)
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}

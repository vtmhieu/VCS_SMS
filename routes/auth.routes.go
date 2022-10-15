package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/controllers"
	"github.com/vtmhieu/VCS_SMS/middleware"
)

type Auth_Route_controller struct {
	authcontroller controllers.Auth_controller
}

func New_Auth_Route_controller(authcontroller controllers.Auth_controller) Auth_Route_controller {
	return Auth_Route_controller{authcontroller}
}

func (ac *Auth_Route_controller) AuthRoute(rg *gin.RouterGroup) {
	// validator := NewOpenAPIMiddleware()
	router := rg.Group("auths")
	// router.Use(validator)
	router.POST("/register", ac.authcontroller.Sign_up)
	router.POST("/login", ac.authcontroller.Sign_in)
	router.GET("/refresh", ac.authcontroller.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(), ac.authcontroller.LogoutUser)
}

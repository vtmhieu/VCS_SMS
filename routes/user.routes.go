package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/controllers"
	"github.com/vtmhieu/VCS_SMS/middleware"
)

type User_Route_controller struct {
	usercontroller controllers.User_controller
}

func New_user_route_controller(userController controllers.User_controller) User_Route_controller {
	return User_Route_controller{userController}
}

func (uc *User_Route_controller) UserRoute(rg *gin.RouterGroup) {
	validator := NewOpenAPIMiddleware()
	router := rg.Group("users")
	router.Use(validator)
	router.GET("/me", middleware.DeserializeUser(), uc.usercontroller.GetMe)
}

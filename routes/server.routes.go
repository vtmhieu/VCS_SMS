package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/controllers"
)

type Server_Route_Controller struct {
	servercontroller controllers.Server_controller
}

func New_route_server_controller(servercontroller controllers.Server_controller) *Server_Route_Controller {
	return &Server_Route_Controller{servercontroller}
}

func (c *Server_Route_Controller) Server_Route(rg *gin.RouterGroup) {
	router := rg.Group("servers")
	router.POST("/", c.servercontroller.CreateServer)

}

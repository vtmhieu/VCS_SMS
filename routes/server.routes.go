package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/controllers"
	"github.com/vtmhieu/VCS_SMS/middleware"
)

type Server_Route_Controller struct {
	servercontroller controllers.Server_controller
}

func New_route_server_controller(servercontroller controllers.Server_controller) Server_Route_Controller {
	return Server_Route_Controller{servercontroller}
}

//	func NewOpenAPIMiddleware() gin.HandlerFunc {
//		validator := middleware.OpenapiInputValidator("./openapi.yaml")
//		return validator
//	}
func (c *Server_Route_Controller) Server_Route(rg *gin.RouterGroup) {
	// validator := NewOpenAPIMiddleware()
	router := rg.Group("servers")
	// router.Use(validator)
	router.Use(middleware.DeserializeUser())
	router.POST("/", c.servercontroller.CreateServer)
	router.PUT("/:server_id", c.servercontroller.UpdateServer)
	router.GET("/:server_id", c.servercontroller.GetServer)
	router.GET("/", c.servercontroller.GetAllServer)
	router.DELETE("/:server_id", c.servercontroller.DeleteServer)
	router.DELETE("/", c.servercontroller.Delete_all_servers)
	router.POST("/all", c.servercontroller.CreatemanyServer)
	router.GET("/all/port", c.servercontroller.Check_on_off)
	router.GET("/ipv4/:ipv4", c.servercontroller.Check)
	router.GET("/excel/export", c.servercontroller.Export_Excel)
	router.POST("/excel/import", c.servercontroller.Post_by_excel)
}

package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/models"
	"gorm.io/gorm"
)

type Server_controller struct {
	DB *gorm.DB
}

func New_server_controller(DB *gorm.DB) Server_controller {
	return Server_controller{DB}
}

// create a new server
func (sc *Server_controller) CreateServer(ctx *gin.Context) {
	var payload *models.Create_server

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	var new_server models.Server
	newServer := models.Server{
		Server_id:    payload.Server_id,
		Server_name:  payload.Server_name,
		Status:       payload.Status,
		Created_time: now,
		Last_updated: now,
		Ipv4:         payload.Ipv4,
	}

	new_server = newServer
	results := sc.DB.Create(&newServer)
	if results.Error != nil {
		if strings.Contains(results.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "message": results.Error.Error()})
			return
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": new_server})
}

// Create many servers at one time
func (sc *Server_controller) CreatemanyServer(ctx *gin.Context) {
	var payload *models.Create_many_server

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	now := time.Now()
	for x, y := range payload.Create_server {
		var new_server models.Server
		newServer := models.Server{
			Server_id:    y.Server_id,
			Server_name:  y.Server_name,
			Status:       y.Status,
			Created_time: now,
			Last_updated: now,
			Ipv4:         y.Ipv4,
		}
		new_server = newServer
		results := sc.DB.Create(&newServer)
		if results.Error != nil {
			if strings.Contains(results.Error.Error(), "duplicate key") {
				ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "message": results.Error.Error()})
				return
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "number": x + 1, "data": new_server})
	}
}

//update server

func (sc *Server_controller) UpdateServer(ctx *gin.Context) {
	server_id := ctx.Param("server_id")
	var payload *models.Update_server

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var updatedServer models.Server

	result := sc.DB.First(&updatedServer, "server_id = ?", server_id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "no server found"})
		return
	}

	now := time.Now()

	server_to_updated := models.Server{
		Server_name:  payload.Server_name,
		Status:       payload.Status,
		Created_time: updatedServer.Created_time,
		Last_updated: now,
		Ipv4:         payload.Ipv4,
	}

	sc.DB.Model(&updatedServer).Updates(server_to_updated)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedServer})

}

//get specific server

func (sc *Server_controller) GetServer(ctx *gin.Context) {
	server_id := ctx.Param("server_id")

	var foundServer models.Server
	result := sc.DB.First(&foundServer, "server_id=?", server_id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "failed to find server"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "server found", "data": foundServer})
}

//get all server

func (sc *Server_controller) GetAllServer(ctx *gin.Context) {
	var from = ctx.DefaultQuery("from", "1")
	var to = ctx.DefaultQuery("to", "10000")

	int_from, _ := strconv.Atoi(from)
	int_to, _ := strconv.Atoi(to)
	var Servers []models.Server

	results := sc.DB.Offset(int_from - 1).Limit(int_to - int_from + 1).Find(&Servers)

	if results.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "message": "no connection"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "number of servers": len(Servers), "data": Servers})
}

//delete server

func (sc *Server_controller) DeleteServer(ctx *gin.Context) {
	server_id := ctx.Param("server_id")

	var Server_to_delete models.Server

	result := sc.DB.Offset(0).Delete(&Server_to_delete, "server_id = ?", server_id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad", "message": "no server id found."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"server_id": server_id, "status": "deleted"})
}

//delete all servers

func (sc *Server_controller) Delete_all_servers(ctx *gin.Context) {
	var Servers []models.Server

	results := sc.DB.Delete(&Servers)
	if results.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": results.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"results": "all data has been deleted successfully"})
}

func (sc *Server_controller) Post_by_excel(ctx *gin.Context) {

}

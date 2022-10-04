package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/models"
	"gorm.io/gorm"
)

type Server_controller struct {
	DB *gorm.DB
}

func New_server_controller(DB *gorm.DB) *Server_controller {
	return &Server_controller{DB}
}

// create a new server
func (sc *Server_controller) CreateServer(ctx *gin.Context) {
	var payload *models.Create_server

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()

	newServer := models.Server{
		Server_name:  payload.Server_name,
		Status:       payload.Status,
		Created_time: now,
		Last_updated: now,
		Ipv4:         payload.Ipv4,
	}

	results := sc.DB.Create(&newServer)
	if strings.Contains(results.Error.Error(), "duplicate key") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "message": results.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newServer})
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

	result := sc.DB.First(&updatedServer, "id = ?", server_id)

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

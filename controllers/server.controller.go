package controllers

import (
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
	var a *models.Create_server

}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/models"
	"gorm.io/gorm"
)

type User_controller struct {
	db *gorm.DB
}

func New_user_controller(DB *gorm.DB) User_controller {
	return User_controller{DB}
}

func (uc *User_controller) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := models.User{
		User_id:         currentUser.User_id,
		User_name:       currentUser.User_name,
		User_email:      currentUser.User_email,
		User_created:    currentUser.User_created,
		User_updated_at: currentUser.User_updated_at,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

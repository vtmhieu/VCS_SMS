package controllers

import (
	"net/http"
	"strings"
	"time"

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
func (uc *User_controller) CreateUser(ctx *gin.Context) {
	var payload *models.Sign_up

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	now := time.Now()
	user := models.User{
		User_id:         payload.User_id,
		User_name:       payload.User_name,
		User_password:   payload.User_password,
		User_email:      payload.User_email,
		User_created:    now,
		User_updated_at: now,
	}

	new_user := user
	results := uc.db.Create(&user)
	if results.Error != nil {
		if strings.Contains(results.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "message": results.Error.Error()})
			return
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "data": gin.H{"user": new_user}})
}

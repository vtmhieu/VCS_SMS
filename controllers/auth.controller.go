package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/models"
	"github.com/vtmhieu/VCS_SMS/utils"
	"gorm.io/gorm"
)

type Auth_controller struct {
	db *gorm.DB
}

func New_auth_controller(db *gorm.DB) Auth_controller {
	return Auth_controller{db}
}

// Sign_up
func (ac *Auth_controller) Sign_up(ctx *gin.Context) {
	var payload *models.Sign_up
	//check JSON input
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	//confirm password
	if payload.User_password_confirmation != payload.User_password {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong password confirmation"})
	}

	hashed_password, err := utils.HashPassword(payload.User_password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
	}

	now := time.Now()
	newUser := models.User{
		User_id:         payload.User_id,
		User_name:       payload.User_name,
		User_password:   hashed_password,
		User_email:      strings.ToLower(payload.User_email),
		User_created:    now,
		User_updated_at: now,
	}
	result := ac.db.Create(&newUser)
	//check if duplicate
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "message": err.Error()})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "message": err.Error()})
		return
	} //check if cant connect

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": newUser}})
}

// Sign in
func (ac *Auth_controller) Sign_in(ctx *gin.Context) {
	var payload *models.Sign_in
	//check JSON input
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}
	var user models.User
	result := ac.db.First(&user, "email = ?", strings.ToLower(payload.User_email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email"})
		return
	}

	if err := utils.VerifyPassword(user.User_password, payload.User_password); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "message": "Invalid password"})
		return
	}

}

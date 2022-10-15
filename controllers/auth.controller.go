package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/initializers"
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
		return
	}

	hashed_password, err := utils.HashPassword(payload.User_password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User_response{
		// User_id:         payload.User_id,
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

	New_user := newUser
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": gin.H{"user": New_user}})
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

	config, _ := initializers.LoadConfig(".")
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.User_id, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.User_id, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "failed", "message": err.Error()})
	}
	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})

}

// refresh access token

func (ac *Auth_controller) RefreshAccessToken(ctx *gin.Context) {
	message := "Couldn't refresh access token"
	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := ac.db.First(&user, "user_id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.User_id, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}
func (ac *Auth_controller) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

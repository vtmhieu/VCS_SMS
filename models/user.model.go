package models

import "time"

type User struct {
	User_id         int       `gorm:"uniqueIndex;not null;primary key" json:"user_id"`
	User_name       string    `gorm:"type:varchar(225)" json:"user_name"`
	User_password   string    ` gorm:"not null" json:"password"`
	User_email      string    `gorm:"uniqueIndex;not null" json:"email"`
	User_created    time.Time `gorm:"not null" json:"created"`
	User_updated_at time.Time `gorm:"not null" json:"updated_at"`
}

type Sign_up struct {
	User_id       int    `json:"user_id" binding:"required"`
	User_name     string `json:"user_name" binding:"required"`
	User_password string `json:"password" binding:"required"`
	User_email    string `json:"email" binding:"required"`
}

type Sign_in struct {
	User_name     string `json:"user_name" binding:"required"`
	User_password string `json:"password" binding:"required"`
}

type User_response struct {
	User_id         int       `json:"user_id"`
	User_name       string    `json:"user_name"`
	User_password   string    `json:"password"`
	User_email      string    `json:"email"`
	User_created    time.Time `json:"created"`
	User_updated_at time.Time `json:"updated_at"`
}
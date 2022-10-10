package models

import (
	"time"
)

type Server struct {
	Server_id    string    `gorm:"uniqueIndex;not null;primary_key" json:"server_id,omitempty" excel:"server_id,omitempty"`
	Server_name  string    `gorm:"uniqueIndex;not null" json:"server_name,omitempty" excel:"server_name,omitempty"`
	Status       string    `gorm:"not null" json:"status,omitempty" excel:"status,omitempty"`
	Created_time time.Time `gorm:"not null" json:"created_time,omitempty" excel:"created_time,omitempty"`
	Last_updated time.Time `gorm:"not null" json:"last_updated,omitempty" excel:"last_updated,omitempty"`
	Ipv4         string    `gorm:"not null" json:"ipv4,omitempty" excel:"ipv4,omitempty"`
}

type Create_server struct {
	Server_id    string    `json:"server_id" binding:"required"`
	Server_name  string    `json:"server_name" binding:"required"`
	Status       string    `json:"status" binding:"required"`
	Created_time time.Time `json:"created_time,omitempty"`
	Last_updated time.Time `json:"last_updated,omitempty"`
	Ipv4         string    `json:"ipv4,omitempty"`
}

type Create_many_server []Create_server

type Update_server struct {
	Server_name  string    `json:"server_name,omitempty"`
	Status       string    `json:"status,omitempty"`
	Created_time time.Time `json:"created_time,omitempty"`
	Last_updated time.Time `json:"last_updated,omitempty"`
	Ipv4         string    `json:"ipv4,omitempty"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Server struct {
	Server_id    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"server_id,omitempty"`
	Server_name  string    `gorm:"uniqueIndex;not null" json:"server_name,omitempty"`
	Status       string    `gorm:"not null" json:"status,omitempty"`
	Created_time time.Time `gorm:"not null" json:"created_time,omitempty"`
	Last_updated time.Time `gorm:"not null" json:"last_updated,omitempty"`
	Ipv4         string    `gorm:"not null" json:"ipv4,omitempty"`
}

type Create_server struct {
	Server_name  string    `json:"server_name" binding:"required"`
	Status       string    `json:"status" binding:"required"`
	Created_time time.Time `json:"created_time,omitempty"`
	Last_updated time.Time `json:"last_updated,omitempty"`
	Ipv4         string    `json:"ipv4,omitempty"`
}

type Update_server struct {
	Server_name  string    `json:"server_name,omitempty"`
	Status       string    `json:"status,omitempty"`
	Created_time time.Time `json:"created_time,omitempty"`
	Last_updated time.Time `json:"last_updated,omitempty"`
	Ipv4         string    `json:"ipv4,omitempty"`
}

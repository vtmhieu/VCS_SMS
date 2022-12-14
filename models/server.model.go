package models

import (
	"time"
)

type Server struct {
	Server_id    string    `gorm:"uniqueIndex;not null;primary_key" json:"server_id,omitempty" xml:"server_id,omitempty"`
	Server_name  string    `gorm:"uniqueIndex;not null" json:"server_name,omitempty" xml:"server_name,omitempty"`
	Status       string    `gorm:"not null" json:"status,omitempty" xml:"status,omitempty"`
	User_id      string    `gorm:"" json:"user_id,omitempty"`
	Created_time time.Time `gorm:"not null" json:"created_time,omitempty" xml:"created_time,omitempty"`
	Last_updated time.Time `gorm:"not null" json:"last_updated,omitempty" xml:"last_updated,omitempty"`
	Ipv4         string    `gorm:"not null" json:"ipv4,omitempty" xml:"ipv4,omitempty"`
}

type Create_server struct {
	Server_id    string    `json:"server_id" binding:"required"`
	Server_name  string    `json:"server_name" binding:"required"`
	Status       string    `json:"status" binding:"required"`
	User_id      string    `json:"user_id,omitempty"`
	Created_time time.Time `json:"created_time,omitempty"`
	Last_updated time.Time `json:"last_updated,omitempty"`
	Ipv4         string    `json:"ipv4,omitempty"`
}

type Create_many_server struct {
	Create_server []struct {
		Server_id    string    `json:"server_id" binding:"required"`
		Server_name  string    `json:"server_name" binding:"required"`
		Status       string    `json:"status" binding:"required"`
		User_id      string    `json:"user_id,omitempty"`
		Created_time time.Time `json:"created_time,omitempty"`
		Last_updated time.Time `json:"last_updated,omitempty"`
		Ipv4         string    `json:"ipv4,omitempty" binding:"required"`
	} `json:"servers"`
}

type Update_server struct {
	Server_name  string    `json:"server_name,omitempty"`
	Status       string    `json:"status,omitempty"`
	User_id      string    `json:"user_id,omitempty"`
	Created_time time.Time `json:"created_time,omitempty"`
	Last_updated time.Time `json:"last_updated,omitempty"`
	Ipv4         string    `json:"ipv4,omitempty"`
}
type Response_API struct {
	Status    string `json:"status,omitempty"`
	Number    int    `json:"number,omitempty"`
	Server_ID string `json:"server_id,omitempty"`
}

type Daily_API struct {
	Start      string   `json:"start,omitempty"`
	End        string   `json:"end,omitempty"`
	List_Email []string `json:"list_email,omitempty"`
}

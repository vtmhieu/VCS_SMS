package main

import (
	"fmt"
	"log"

	"github.com/vtmhieu/VCS_SMS/initializers"
	"github.com/vtmhieu/VCS_SMS/models"
)

// loaded the environment variables and created a connection pool to the Postgres database in the init() function
func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

//Then, evoked the AutoMigrate() function provided by GORM to create the database migration and push the changes to the database.

func main() {
	initializers.DB.AutoMigrate(&models.Server{}, &models.User{})
	fmt.Println("👍 Migration complete")
}

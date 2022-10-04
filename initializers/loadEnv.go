package initializers

import "github.com/spf13/viper"

const (
	POSTGRES_HOST     = "localhost"
	POSTGRES_USER     = "postgres"
	POSTGRES_PASSWORD = "password123"
	POSTGRES_DB       = "VCS_SMS"
	POSTGRES_PORT     = "6500"
	PORT              = "8000"

	CLIENT_ORIGIN = "http://localhost:3000"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
}

// create function to load the environment variables
// from the app.env file and make them accessible in other files and packages within the application code.

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

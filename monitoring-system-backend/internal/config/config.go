package config

import (
	"github.com/joho/godotenv"
	"log"
)

func InitEnvVariables() {
	//Loads the ENV file
	if err := godotenv.Load(); err != nil {
		log.Println(("No .env file found"))
	}
}


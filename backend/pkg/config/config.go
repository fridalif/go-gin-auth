package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var moduleName string = "config"

type AppConfig struct {
	WebPort    string
	DbUrl      string
	DbUser     string
	DbPass     string
	DbName     string
	DbPort     string
	SecretKey  string
	Email      string
	EmailToken string
}

func NewAppConfig() (*AppConfig, error) {
	var functionName string = "NewAppConfig"

	err := godotenv.Load()
	if err != nil {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %v", moduleName, functionName, err)
	}

	webPort := os.Getenv("WEB_PORT")
	dbPort := os.Getenv("DB_PORT")
	dbUrl := os.Getenv("DB_URL")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	secretKey := os.Getenv("SECRET_KEY")
	email := os.Getenv("EMAIL")
	emailToken := os.Getenv("EMAIL_TOKEN")
	if webPort == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "WEB_PORT")
	}
	if dbPort == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "DB_PORT")
	}
	if dbUrl == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "DB_URL")
	}
	if dbUser == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "DB_USER")
	}
	if dbPass == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "DB_PASS")
	}
	if dbName == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "DB_NAME")
	}
	if secretKey == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "SECRET_KEY")
	}
	if email == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "EMAIL")
	}
	if emailToken == "" {
		return &AppConfig{}, fmt.Errorf("%s.%s:ERROR: %s not exists", moduleName, functionName, "EMAIL_TOKEN")
	}
	return &AppConfig{
		WebPort:    webPort,
		DbUrl:      dbUrl,
		DbUser:     dbUser,
		DbPass:     dbPass,
		DbName:     dbName,
		DbPort:     dbPort,
		SecretKey:  secretKey,
		EmailToken: emailToken,
		Email:      email,
	}, nil
}

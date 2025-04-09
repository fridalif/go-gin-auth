package main

import (
	"fmt"
	"hitenok/pkg/config"
	"hitenok/pkg/domain"
	"hitenok/pkg/handlers"
	"hitenok/pkg/repository"
	"hitenok/pkg/services"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func runServer(db *gorm.DB, appConfig *config.AppConfig) {
	router := gin.Default()
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("runserver.AutoMigrate.Error: %v", err)
	}
	api := router.Group("/api")
	v1 := api.Group("/v1")
	auth := v1.Group("/auth")

	userRepo := repository.NewUserRepository(db, appConfig)

	mailAuthenticationService := services.NewMailAuthenticationService(userRepo, appConfig)
	otpService := services.NewMailOTPService(userRepo, appConfig)
	jwtService := services.NewJWTService(appConfig, userRepo)
	hashService := services.NewHashService(userRepo)
	userService := services.NewUserService(userRepo)

	mailAuthenticationHandler := handlers.NewMailAuthHandler(mailAuthenticationService, otpService, jwtService, appConfig)
	mailAuthenticationHandler.RegisterRoutes(auth)
	activateServiceHandler := handlers.NewActivateHandler(otpService, hashService, jwtService, userService, appConfig)
	activateServiceHandler.RegisterRoutes(auth)
	router.Run(fmt.Sprintf(":%s", appConfig.WebPort))
}

func main() {
	appConfig, err := config.NewAppConfig()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", appConfig.DbUrl, appConfig.DbUser, appConfig.DbPass, appConfig.DbName, appConfig.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("main.connection_to_database.Error: %v", err)
		return
	}
	runServer(db, appConfig)
}

package main

import (
	testTaskObjects "JWTService"
	db2 "JWTService/internal/db"
	"JWTService/internal/handler"
	repository "JWTService/internal/repository"
	"JWTService/internal/service"
	"JWTService/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	err := logger.InitLogger("app.log", viper.GetString("logLevel"))
	if err != nil {
		log.Fatal("Logger init error", err)
	}

	err = initConfig()
	if err != nil {
		logger.GetLogger().Error("Config init error", err)
	}
	logger.GetLogger().Info("Config loaded")

	db, err := db2.Connect(db2.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})


	if err != nil {
		logger.GetLogger().Error("DB init error", err)
	}
	logger.GetLogger().Info("DB connected")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(testTaskObjects.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logger.GetLogger().Error("Server error", err)
	}

	logger.GetLogger().Info("Server started")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

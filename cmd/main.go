package main

import (
	testTaskObjects "JWTService"
	db2 "JWTService/internal/db"
	"JWTService/internal/handler"
	repository "JWTService/internal/repository"
	"JWTService/internal/service"
	"JWTService/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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

	db, err := db2.Connect(db2.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(testTaskObjects.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

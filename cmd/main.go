package main

import (
	testTaskObjects "JWTService"
	db2 "JWTService/internal/db"
	"JWTService/internal/handler"
	repository "JWTService/internal/repository"
	"JWTService/internal/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	db, err := db2.NewMongoDB(db2.Config{
		Host: "localhost",
		Port: "27017",
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(testTaskObjects.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}
}

package main

import (
	"time"

	"github.com/george007361/db-course-proj/app/handler"
	"github.com/george007361/db-course-proj/app/repository"
	"github.com/george007361/db-course-proj/app/repository/postgres"
	"github.com/george007361/db-course-proj/app/server"
	"github.com/george007361/db-course-proj/app/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     "localhost",
		Port:     "5438",
		Username: "postgres",
		Password: "12345678",
		DBName:   "postgres",
		// Host:            "localhost",
		// Port:            "5432",
		// Username:        "george",
		// Password:        "12345678",
		// DBName:          "george_forum_db",
		SSLMode:         "disable",
		MaxOpenConns:    100,
		MaxIdleConns:    50,
		ConnMaxLifeTime: time.Second * 100,
	})

	if err != nil {
		logrus.Fatalf("Cant connect to db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(server.Server)
	if err := server.Run("5000", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Cant run server: %s", err.Error())
	}
}

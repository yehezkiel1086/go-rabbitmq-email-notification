package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/adapter/handler"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/adapter/storage/rabbitmq"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/core/domain"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/core/service"
)

func handleError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %v", msg, err)
		os.Exit(1)
	}
}

func main() {
	// load .env configs
	conf, err := config.New()
	handleError(err, "unable to load .env configs")
	fmt.Println(".env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init db connection
	db, err := postgres.New(ctx, conf.DB)
	handleError(err, "unable to connect to db")
	fmt.Println("DB connection established successfully")

	// migrate dbs
	err = db.Migrate(&domain.User{})
	handleError(err, "unable to migrate db")
	fmt.Println("DB migrated successfully")

	// init rabbitmq
	mq, err := rabbitmq.New(conf.RabbitMQ)
	handleError(err, "rabbitmq connection failed")
	fmt.Println("rabbitmq connected successfully")

	defer mq.Close()

	// dependency injections
	userRepo := repository.NewUserRepository(db)

	userSvc, err := service.NewUserService(userRepo, mq)
	handleError(err, "failed to init user service")

	userHandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	// init router
	r := handler.NewRouter(
		conf.HTTP,
		userHandler,
		authHandler,
	)

	// start server
	err = r.Serve(conf.HTTP)
	handleError(err, "unable to start server")
}

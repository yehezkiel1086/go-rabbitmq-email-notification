package service

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/adapter/storage/rabbitmq"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/core/domain"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/core/port"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/user-service/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
	mq *rabbitmq.Rabbitmq
	q *amqp.Queue
}

func NewUserService(repo port.UserRepository, mq *rabbitmq.Rabbitmq) (*UserService, error) {
	q, err := mq.DeclareQueue("notif_queue")
	if err != nil {
		return nil, err
	}

	return &UserService{
		repo,
		mq,
		q,
	}, nil
}

func (us *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.UserResponse, error) {
	// hash password
	hashedPwd, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPwd

	// create user
	createdUser, err := us.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// send email notification
	if err := us.mq.Send(us.q, fmt.Sprintf("%s: user registered successfully", createdUser.Email)); err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (us *UserService) GetUsers(ctx context.Context) ([]domain.UserResponse, error) {
	return us.repo.GetUsers(ctx)
}

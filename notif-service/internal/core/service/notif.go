package service

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/notif-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/notif-service/internal/core/port"
	"github.com/yehezkiel1086/go-rabbitmq-email-notification/notif-service/internal/core/util"
	"gopkg.in/gomail.v2"
)

type NotifService struct {
	repo port.NotifRepository
	conf *config.GMAIL
}

func NewNotifService(repo port.NotifRepository, conf *config.GMAIL) *NotifService {
	return &NotifService{
		repo,
		conf,
	}
}

func (s *NotifService) ReceiveNotif(ctx context.Context) (<-chan amqp.Delivery, error) {
	return s.repo.ReceiveNotif(ctx)
}

func (s *NotifService) SendConfirmationEmail(ctx context.Context, msg []byte) error {
	var data map[string]string
	if err := util.Deserialize(msg, &data); err != nil {
		log.Printf("failed to deserialized json data: %v", err)
		return err
	}

	url := util.GenerateConfirmationURL(data["confirmation_token"])

	log.Printf("[+] Sending confirmation link %s to %s", url, data["email"])

	// TODO: add gomail to send confirmation link to email
	m := gomail.NewMessage()

	m.SetHeader("From", s.conf.SenderEmail)
	m.SetHeader("To", data["email"])
	m.SetHeader("Subject", "User registration confirmation")

	m.SetBody("text/html", `
		<h2>Email Confirmation</h2>
		<p>Please confirm your email by clicking the link below:</p>
		<p>
			<a href="`+url+`">Confirm Email</a>
		</p>
		<p>This link will expire in 15 minutes.</p>
	`)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		s.conf.SenderEmail,
		s.conf.Password,
	)

	d.SSL = false // STARTTLS

	return d.DialAndSend(m)
}

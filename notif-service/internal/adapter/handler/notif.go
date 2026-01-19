package handler

import (
	"context"
	"log"

	"github.com/yehezkiel1086/go-rabbitmq-email-notification/notif-service/internal/core/port"
)

type NotifHandler struct {
	svc port.NotifService
}

func NewNotifHandler(svc port.NotifService) *NotifHandler {
	return &NotifHandler{
		svc: svc,
	}
}

func (h *NotifHandler) ReceiveNotif(ctx context.Context) {
	msgs, err := h.svc.ReceiveNotif(ctx)
	if err != nil {
		return
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			if err := h.svc.SendConfirmationEmail(ctx, msg.Body); err != nil {
				log.Fatalf("failed to send email: %v", err)
				return
			}
		}
	}()
	
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

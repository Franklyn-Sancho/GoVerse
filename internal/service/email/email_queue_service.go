package email

import (
	"GoVersi/internal/infrastrucuture/queue"
	"encoding/json"
)

type EmailQueueService interface {
	PublishEmail(msg EmailMessage) error
}

type emailQueueService struct {
	rabbitmq queue.RabbitMQClient
}

func NewEmailQueueService(rabbitmq queue.RabbitMQClient) EmailQueueService {
	return &emailQueueService{rabbitmq: rabbitmq}
}

func (s *emailQueueService) PublishEmail(msg EmailMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return s.rabbitmq.Publish("email_queue", body)
}

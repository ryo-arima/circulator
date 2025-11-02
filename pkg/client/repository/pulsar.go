package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
)

// PulsarRepository handles Pulsar messaging for Client
type PulsarRepository struct {
	config   config.BaseConfig
	client   pulsar.Client
	producer pulsar.Producer
	consumer pulsar.Consumer
}

// NewPulsarRepository creates a new PulsarRepository instance
func NewPulsarRepository(cfg config.BaseConfig, pulsarURL string) (*PulsarRepository, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: pulsarURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pulsar client: %w", err)
	}

	// Create producer
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "client-commands",
		Name:  "client-producer",
	})
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create pulsar producer: %w", err)
	}

	// Create consumer
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "client-notifications",
		SubscriptionName: "client-consumer",
		Type:             pulsar.Shared,
	})
	if err != nil {
		producer.Close()
		client.Close()
		return nil, fmt.Errorf("failed to create pulsar consumer: %w", err)
	}

	repo := &PulsarRepository{
		config:   cfg,
		client:   client,
		producer: producer,
		consumer: consumer,
	}

	cfg.Logger.INFO(config.CRPINIT, "Client Pulsar repository initialized", map[string]interface{}{
		"pulsar_url": pulsarURL,
	})

	return repo, nil
}

// PublishCommand publishes a command to Pulsar
func (r *PulsarRepository) PublishCommand(ctx context.Context, command *model.Command) error {
	r.config.Logger.DEBUG(config.CRPPUB, "Publishing command to Pulsar", map[string]interface{}{
		"command_type": command.Type,
		"command_id":   command.ID,
	})

	payload, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	msgID, err := r.producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: payload,
		Key:     command.ID,
		Properties: map[string]string{
			"type":      command.Type,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
	if err != nil {
		r.config.Logger.ERROR(config.CRPERR, "Failed to publish command", map[string]interface{}{
			"error":        err.Error(),
			"command_type": command.Type,
			"command_id":   command.ID,
		})
		return fmt.Errorf("failed to publish command: %w", err)
	}

	r.config.Logger.INFO(config.CRPSUCC, "Command published successfully", map[string]interface{}{
		"message_id":   msgID.String(),
		"command_type": command.Type,
		"command_id":   command.ID,
	})

	return nil
}

// ConsumeNotifications consumes notifications from Pulsar
func (r *PulsarRepository) ConsumeNotifications(ctx context.Context, handler func(*model.Notification) error) error {
	r.config.Logger.INFO(config.CRPCONS, "Starting notification consumption from Pulsar", nil)

	for {
		select {
		case <-ctx.Done():
			r.config.Logger.INFO(config.CRPSTOP, "Stopping notification consumption", nil)
			return ctx.Err()
		default:
			msg, err := r.consumer.Receive(ctx)
			if err != nil {
				r.config.Logger.ERROR(config.CRPERR, "Failed to receive message", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			var notification model.Notification
			if err := json.Unmarshal(msg.Payload(), &notification); err != nil {
				r.config.Logger.ERROR(config.CRPERR, "Failed to unmarshal notification", map[string]interface{}{
					"error":      err.Error(),
					"message_id": msg.ID().String(),
				})
				r.consumer.Ack(msg)
				continue
			}

			r.config.Logger.DEBUG(config.CRPREC, "Received notification", map[string]interface{}{
				"notification_type": notification.Type,
				"notification_id":   notification.ID,
				"message_id":        msg.ID().String(),
			})

			if err := handler(&notification); err != nil {
				r.config.Logger.ERROR(config.CRPERR, "Failed to handle notification", map[string]interface{}{
					"error":             err.Error(),
					"notification_type": notification.Type,
					"notification_id":   notification.ID,
				})
				r.consumer.Nack(msg)
				continue
			}

			r.consumer.Ack(msg)
			r.config.Logger.DEBUG(config.CRPSUCC, "Notification processed successfully", map[string]interface{}{
				"notification_type": notification.Type,
				"notification_id":   notification.ID,
			})
		}
	}
}

// Close closes the Pulsar repository
func (r *PulsarRepository) Close() {
	if r.consumer != nil {
		r.consumer.Close()
	}
	if r.producer != nil {
		r.producer.Close()
	}
	if r.client != nil {
		r.client.Close()
	}
	r.config.Logger.INFO(config.CRPCLOSE, "Client Pulsar repository closed", nil)
}
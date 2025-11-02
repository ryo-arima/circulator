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

// PulsarRepository handles Pulsar messaging for Server
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

	// Create producer for server events
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "server-events",
		Name:  "server-producer",
	})
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create pulsar producer: %w", err)
	}

	// Create consumer for agent reports
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "agent-reports",
		SubscriptionName: "server-consumer",
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

	cfg.Logger.INFO(config.SRPINIT, "Server Pulsar repository initialized", map[string]interface{}{
		"pulsar_url": pulsarURL,
	})

	return repo, nil
}

// PublishEvent publishes an event to Pulsar
func (r *PulsarRepository) PublishEvent(ctx context.Context, event *model.ServerEvent) error {
	r.config.Logger.DEBUG(config.SRPPUB, "Publishing event to Pulsar", map[string]interface{}{
		"event_type": event.Type,
		"event_id":   event.ID,
	})

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msgID, err := r.producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: payload,
		Key:     event.ID,
		Properties: map[string]string{
			"type":      event.Type,
			"source":    "server",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
	if err != nil {
		r.config.Logger.ERROR(config.SRPERR, "Failed to publish event", map[string]interface{}{
			"error":      err.Error(),
			"event_type": event.Type,
			"event_id":   event.ID,
		})
		return fmt.Errorf("failed to publish event: %w", err)
	}

	r.config.Logger.INFO(config.SRPSUCC, "Event published successfully", map[string]interface{}{
		"message_id": msgID.String(),
		"event_type": event.Type,
		"event_id":   event.ID,
	})

	return nil
}

// PublishNotification publishes a notification to clients
func (r *PulsarRepository) PublishNotification(ctx context.Context, notification *model.Notification) error {
	r.config.Logger.DEBUG(config.SRPNOT, "Publishing notification to Pulsar", map[string]interface{}{
		"notification_type": notification.Type,
		"notification_id":   notification.ID,
	})

	// Create producer for client notifications if not exists
	notificationProducer, err := r.client.CreateProducer(pulsar.ProducerOptions{
		Topic: "client-notifications",
		Name:  "server-notification-producer",
	})
	if err != nil {
		return fmt.Errorf("failed to create notification producer: %w", err)
	}
	defer notificationProducer.Close()

	payload, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msgID, err := notificationProducer.Send(ctx, &pulsar.ProducerMessage{
		Payload: payload,
		Key:     notification.ID,
		Properties: map[string]string{
			"type":      notification.Type,
			"source":    "server",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
	if err != nil {
		r.config.Logger.ERROR(config.SRPERR, "Failed to publish notification", map[string]interface{}{
			"error":             err.Error(),
			"notification_type": notification.Type,
			"notification_id":   notification.ID,
		})
		return fmt.Errorf("failed to publish notification: %w", err)
	}

	r.config.Logger.INFO(config.SRPSUCC, "Notification published successfully", map[string]interface{}{
		"message_id":        msgID.String(),
		"notification_type": notification.Type,
		"notification_id":   notification.ID,
	})

	return nil
}

// ConsumeAgentReports consumes agent reports from Pulsar
func (r *PulsarRepository) ConsumeAgentReports(ctx context.Context, handler func(*model.AgentReport) error) error {
	r.config.Logger.INFO(config.SRPCONS, "Starting agent report consumption from Pulsar", nil)

	for {
		select {
		case <-ctx.Done():
			r.config.Logger.INFO(config.SRPSTOP, "Stopping agent report consumption", nil)
			return ctx.Err()
		default:
			msg, err := r.consumer.Receive(ctx)
			if err != nil {
				r.config.Logger.ERROR(config.SRPERR, "Failed to receive message", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			var report model.AgentReport
			if err := json.Unmarshal(msg.Payload(), &report); err != nil {
				r.config.Logger.ERROR(config.SRPERR, "Failed to unmarshal agent report", map[string]interface{}{
					"error":      err.Error(),
					"message_id": msg.ID().String(),
				})
				r.consumer.Ack(msg)
				continue
			}

			r.config.Logger.DEBUG(config.SRPREC, "Received agent report", map[string]interface{}{
				"report_type": report.Type,
				"report_id":   report.ID,
				"agent_id":    report.AgentID,
				"message_id":  msg.ID().String(),
			})

			if err := handler(&report); err != nil {
				r.config.Logger.ERROR(config.SRPERR, "Failed to handle agent report", map[string]interface{}{
					"error":      err.Error(),
					"report_type": report.Type,
					"report_id":   report.ID,
					"agent_id":    report.AgentID,
				})
				r.consumer.Nack(msg)
				continue
			}

			r.consumer.Ack(msg)
			r.config.Logger.DEBUG(config.SRPSUCC, "Agent report processed successfully", map[string]interface{}{
				"report_type": report.Type,
				"report_id":   report.ID,
				"agent_id":    report.AgentID,
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
	r.config.Logger.INFO(config.SRPCLOSE, "Server Pulsar repository closed", nil)
}
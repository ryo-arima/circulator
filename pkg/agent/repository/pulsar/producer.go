package pulsar

import (
"context"
"encoding/json"
"time"

"github.com/apache/pulsar-client-go/pulsar"
"github.com/ryo-arima/circulator/pkg/config"
"github.com/ryo-arima/circulator/pkg/entity/model"
)

// ProducerRepository defines the interface for Pulsar producer operations from agent
type ProducerRepository interface {
	PublishReport(report *model.AgentReport) error
	PublishNotification(notification *model.Notification) error
	Close() error
}

// producerRepository implements ProducerRepository
type producerRepository struct {
	config   *config.BaseConfig
	client   pulsar.Client
	producer pulsar.Producer
}

// NewProducerRepository creates a new Pulsar producer repository for agent
func NewProducerRepository(c *config.BaseConfig) (ProducerRepository, error) {
	c.Logger.DEBUG(config.ARPPINIT, "Initializing Agent Pulsar producer", map[string]interface{}{
"pulsar_url": c.YamlConfig.Pulsar.URL,
})

	client, err := pulsar.NewClient(pulsar.ClientOptions{
URL: c.YamlConfig.Pulsar.URL,
})
	if err != nil {
		c.Logger.ERROR(config.ARPERR, "Failed to create Pulsar client", map[string]interface{}{
"error": err.Error(),
		})
		return nil, err
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
Topic: "agent-reports",
SendTimeout: time.Duration(c.YamlConfig.Pulsar.Producer.SendTimeout) * time.Second,
})
	if err != nil {
		c.Logger.ERROR(config.ARPERR, "Failed to create Pulsar producer", map[string]interface{}{
"error": err.Error(),
		})
		client.Close()
		return nil, err
	}

	repo := &producerRepository{
		config:   c,
		client:   client,
		producer: producer,
	}

	c.Logger.DEBUG(config.ARPSUCC, "Agent Pulsar producer initialized successfully", nil)
	return repo, nil
}

// PublishReport publishes an agent report to Pulsar
func (r *producerRepository) PublishReport(report *model.AgentReport) error {
	r.config.Logger.DEBUG(config.ARPREP, "Agent publishing report to Pulsar", map[string]interface{}{
"agent_id":   report.AgentID,
"report_id":  report.ID,
"report_type": report.Type,
})

	data, err := json.Marshal(report)
	if err != nil {
		r.config.Logger.ERROR(config.ARPERR, "Failed to marshal agent report", map[string]interface{}{
"error": err.Error(),
		})
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = r.producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: data,
		Key:     report.AgentID,
		Properties: map[string]string{
			"type":      "agent_report",
			"agent_id":  report.AgentID,
			"report_id": report.ID,
		},
	})

	if err != nil {
		r.config.Logger.ERROR(config.ARPERR, "Failed to publish agent report", map[string]interface{}{
"error": err.Error(),
			"agent_id": report.AgentID,
		})
		return err
	}

	r.config.Logger.DEBUG(config.ARPSUCC, "Agent report published successfully", map[string]interface{}{
"agent_id":  report.AgentID,
"report_id": report.ID,
})

	return nil
}

// PublishNotification publishes a notification to Pulsar
func (r *producerRepository) PublishNotification(notification *model.Notification) error {
	r.config.Logger.DEBUG(config.ARPNOT, "Agent publishing notification to Pulsar", map[string]interface{}{
"agent_id":        notification.AgentID,
"notification_id": notification.ID,
"type":           notification.Type,
})

	data, err := json.Marshal(notification)
	if err != nil {
		r.config.Logger.ERROR(config.ARPERR, "Failed to marshal notification", map[string]interface{}{
"error": err.Error(),
		})
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = r.producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: data,
		Key:     notification.AgentID,
		Properties: map[string]string{
			"type":            "notification",
			"agent_id":        notification.AgentID,
			"notification_id": notification.ID,
		},
	})

	if err != nil {
		r.config.Logger.ERROR(config.ARPERR, "Failed to publish notification", map[string]interface{}{
"error": err.Error(),
			"agent_id": notification.AgentID,
		})
		return err
	}

	r.config.Logger.DEBUG(config.ARPSUCC, "Notification published successfully", map[string]interface{}{
"agent_id":        notification.AgentID,
"notification_id": notification.ID,
})

	return nil
}

// Close closes the producer and client
func (r *producerRepository) Close() error {
	r.config.Logger.DEBUG(config.ARPCLOSE, "Closing Agent Pulsar producer", nil)

	if r.producer != nil {
		r.producer.Close()
	}
	if r.client != nil {
		r.client.Close()
	}

	r.config.Logger.DEBUG(config.ARPSUCC, "Agent Pulsar producer closed successfully", nil)
	return nil
}

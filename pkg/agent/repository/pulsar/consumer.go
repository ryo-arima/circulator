package pulsar

import (
	"context"
	"encoding/json"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
)// ConsumerRepository defines the interface for Pulsar consumer operations from agent
type ConsumerRepository interface {
	ConsumeCommands(ctx context.Context, handler func(*model.Command) error) error
	ConsumeServerEvents(ctx context.Context, handler func(*model.ServerEvent) error) error
	Close() error
}

// consumerRepository implements ConsumerRepository
type consumerRepository struct {
	config           *config.BaseConfig
	client           pulsar.Client
	commandConsumer  pulsar.Consumer
	eventConsumer    pulsar.Consumer
}

// NewConsumerRepository creates a new Pulsar consumer repository for agent
func NewConsumerRepository(c *config.BaseConfig, agentID string) (ConsumerRepository, error) {
	c.Logger.DEBUG(config.ARCCINIT, "Initializing Agent Pulsar consumer", map[string]interface{}{
"pulsar_url": c.YamlConfig.Pulsar.URL,
"agent_id":   agentID,
})

	client, err := pulsar.NewClient(pulsar.ClientOptions{
URL: c.YamlConfig.Pulsar.URL,
})
	if err != nil {
		c.Logger.ERROR(config.ARCERR, "Failed to create Pulsar client", map[string]interface{}{
"error": err.Error(),
		})
		return nil, err
	}

	// Create consumer for commands directed to this agent
	commandConsumer, err := client.Subscribe(pulsar.ConsumerOptions{
Topic:            "agent-commands",
SubscriptionName: "agent-" + agentID,
Type:             pulsar.Exclusive,
})
	if err != nil {
		c.Logger.ERROR(config.ARCERR, "Failed to create command consumer", map[string]interface{}{
"error": err.Error(),
		})
		client.Close()
		return nil, err
	}

	// Create consumer for server events
	eventConsumer, err := client.Subscribe(pulsar.ConsumerOptions{
Topic:            "server-events",
SubscriptionName: "agent-events-" + agentID,
Type:             pulsar.Shared,
})
	if err != nil {
		c.Logger.ERROR(config.ARCERR, "Failed to create event consumer", map[string]interface{}{
"error": err.Error(),
		})
		commandConsumer.Close()
		client.Close()
		return nil, err
	}

	repo := &consumerRepository{
		config:          c,
		client:          client,
		commandConsumer: commandConsumer,
		eventConsumer:   eventConsumer,
	}

	c.Logger.DEBUG(config.ARCSUCC, "Agent Pulsar consumer initialized successfully", map[string]interface{}{
"agent_id": agentID,
})
	return repo, nil
}

// ConsumeCommands consumes commands from Pulsar
func (r *consumerRepository) ConsumeCommands(ctx context.Context, handler func(*model.Command) error) error {
	r.config.Logger.DEBUG(config.ARCCONS, "Agent starting command consumption", nil)

	for {
		select {
		case <-ctx.Done():
			r.config.Logger.DEBUG(config.ARCSTOP, "Agent stopping command consumption", nil)
			return ctx.Err()
		default:
			msg, err := r.commandConsumer.Receive(ctx)
			if err != nil {
				r.config.Logger.ERROR(config.ARCERR, "Failed to receive command message", map[string]interface{}{
"error": err.Error(),
				})
				continue
			}

			r.config.Logger.DEBUG(config.ARCREC, "Agent received command message", map[string]interface{}{
"message_id": msg.ID().String(),
			})

			var command model.Command
			if err := json.Unmarshal(msg.Payload(), &command); err != nil {
				r.config.Logger.ERROR(config.ARCERR, "Failed to unmarshal command", map[string]interface{}{
"error": err.Error(),
				})
				r.commandConsumer.Nack(msg)
				continue
			}

			r.config.Logger.DEBUG(config.ARCPROC, "Agent processing command", map[string]interface{}{
				"command_id":   command.ID,
				"command_type": command.Type,
				"target":       command.Target,
			})

			if err := handler(&command); err != nil {
				r.config.Logger.ERROR(config.ARCERR, "Failed to process command", map[string]interface{}{
"error":        err.Error(),
					"command_id":   command.ID,
					"command_type": command.Type,
				})
				r.commandConsumer.Nack(msg)
				continue
			}

			r.commandConsumer.Ack(msg)
			r.config.Logger.DEBUG(config.ARCSUCC, "Agent processed command successfully", map[string]interface{}{
"command_id": command.ID,
})
		}
	}
}

// ConsumeServerEvents consumes server events from Pulsar
func (r *consumerRepository) ConsumeServerEvents(ctx context.Context, handler func(*model.ServerEvent) error) error {
	r.config.Logger.DEBUG(config.ARCCONS, "Agent starting server event consumption", nil)

	for {
		select {
		case <-ctx.Done():
			r.config.Logger.DEBUG(config.ARCSTOP, "Agent stopping server event consumption", nil)
			return ctx.Err()
		default:
			msg, err := r.eventConsumer.Receive(ctx)
			if err != nil {
				r.config.Logger.ERROR(config.ARCERR, "Failed to receive server event message", map[string]interface{}{
"error": err.Error(),
				})
				continue
			}

			r.config.Logger.DEBUG(config.ARCREC, "Agent received server event message", map[string]interface{}{
"message_id": msg.ID().String(),
			})

			var event model.ServerEvent
			if err := json.Unmarshal(msg.Payload(), &event); err != nil {
				r.config.Logger.ERROR(config.ARCERR, "Failed to unmarshal server event", map[string]interface{}{
"error": err.Error(),
				})
				r.eventConsumer.Nack(msg)
				continue
			}

			r.config.Logger.DEBUG(config.ARCPROC, "Agent processing server event", map[string]interface{}{
"event_id":   event.ID,
"event_type": event.Type,
"agent_id":   event.AgentID,
})

			if err := handler(&event); err != nil {
				r.config.Logger.ERROR(config.ARCERR, "Failed to process server event", map[string]interface{}{
"error":      err.Error(),
					"event_id":   event.ID,
					"event_type": event.Type,
				})
				r.eventConsumer.Nack(msg)
				continue
			}

			r.eventConsumer.Ack(msg)
			r.config.Logger.DEBUG(config.ARCSUCC, "Agent processed server event successfully", map[string]interface{}{
"event_id": event.ID,
})
		}
	}
}

// Close closes all consumers and client
func (r *consumerRepository) Close() error {
	r.config.Logger.DEBUG(config.ARCCLOSE, "Closing Agent Pulsar consumer", nil)

	if r.commandConsumer != nil {
		r.commandConsumer.Close()
	}
	if r.eventConsumer != nil {
		r.eventConsumer.Close()
	}
	if r.client != nil {
		r.client.Close()
	}

	r.config.Logger.DEBUG(config.ARCSUCC, "Agent Pulsar consumer closed successfully", nil)
	return nil
}

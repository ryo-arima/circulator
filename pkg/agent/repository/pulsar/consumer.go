package pulsar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
)

// StreamDataConsumer handles consuming streaming data from external sources
type StreamDataConsumer struct {
	client   pulsar.Client
	consumer pulsar.Consumer
	config   config.BaseConfig
}

// NewStreamDataConsumer creates a new consumer for external streaming data
func NewStreamDataConsumer(conf config.BaseConfig) (*StreamDataConsumer, error) {
	// Create Pulsar client using configuration
	client, err := config.NewPulsarClient(conf.YamlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Pulsar client: %w", err)
	}

	// Create consumer for external data stream using configuration
	consumerOptions := conf.YamlConfig.GetPulsarConsumerOptions(conf.YamlConfig.Pulsar.Topics.ExternalSensorData)
	consumer, err := client.Subscribe(consumerOptions)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &StreamDataConsumer{
		client:   client,
		consumer: consumer,
		config:   conf,
	}, nil
}

// ConsumeStreamData consumes streaming data from external sources
func (c *StreamDataConsumer) ConsumeStreamData(ctx context.Context, dataHandler func(model.IncomingStreamData) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Receive message with timeout
			msg, err := c.consumer.Receive(ctx)
			if err != nil {
				c.config.Logger.ERROR(config.APCERM, "Error receiving message", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			// Parse incoming data
			var streamData model.IncomingStreamData
			if err := json.Unmarshal(msg.Payload(), &streamData); err != nil {
				c.config.Logger.ERROR(config.APCPSD, "Error parsing stream data", map[string]interface{}{
					"error": err.Error(),
				})
				c.consumer.Nack(msg)
				continue
			}

			// Process the data using the provided handler
			if err := dataHandler(streamData); err != nil {
				c.config.Logger.ERROR(config.APCPSD2, "Error processing stream data", map[string]interface{}{
					"error": err.Error(),
				})
				c.consumer.Nack(msg)
				continue
			}

			// Acknowledge successful processing
			c.consumer.Ack(msg)
		}
	}
}

// Close closes the consumer and client
func (c *StreamDataConsumer) Close() error {
	if c.consumer != nil {
		c.consumer.Close()
	}
	if c.client != nil {
		c.client.Close()
	}
	return nil
}

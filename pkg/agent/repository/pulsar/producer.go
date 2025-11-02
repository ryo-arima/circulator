package pulsar

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
) // StreamDataProducer handles producing processed streaming data
type StreamDataProducer struct {
	client                   pulsar.Client
	processedDataProducer    pulsar.Producer
	systemMetricsProducer    pulsar.Producer
	alertDataProducer        pulsar.Producer
	processingResultProducer pulsar.Producer
	config                   config.BaseConfig
}

// NewStreamDataProducer creates a new producer for processed streaming data
func NewStreamDataProducer(conf config.BaseConfig) (*StreamDataProducer, error) {
	// Create Pulsar client using configuration
	client, err := config.NewPulsarClient(conf.YamlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Pulsar client: %w", err)
	}

	// Create producers for different topics
	processedDataOptions := conf.YamlConfig.GetPulsarProducerOptions(conf.YamlConfig.Pulsar.Topics.ProcessedSensorData)
	processedDataProducer, err := client.CreateProducer(processedDataOptions)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create processed data producer: %w", err)
	}

	systemMetricsOptions := conf.YamlConfig.GetPulsarProducerOptions(conf.YamlConfig.Pulsar.Topics.SystemMetrics)
	systemMetricsProducer, err := client.CreateProducer(systemMetricsOptions)
	if err != nil {
		processedDataProducer.Close()
		client.Close()
		return nil, fmt.Errorf("failed to create system metrics producer: %w", err)
	}

	alertDataOptions := conf.YamlConfig.GetPulsarProducerOptions(conf.YamlConfig.Pulsar.Topics.AlertData)
	alertDataProducer, err := client.CreateProducer(alertDataOptions)
	if err != nil {
		processedDataProducer.Close()
		systemMetricsProducer.Close()
		client.Close()
		return nil, fmt.Errorf("failed to create alert data producer: %w", err)
	}

	processingResultOptions := conf.YamlConfig.GetPulsarProducerOptions(conf.YamlConfig.Pulsar.Topics.ProcessingResults)
	processingResultProducer, err := client.CreateProducer(processingResultOptions)
	if err != nil {
		processedDataProducer.Close()
		systemMetricsProducer.Close()
		alertDataProducer.Close()
		client.Close()
		return nil, fmt.Errorf("failed to create processing result producer: %w", err)
	}

	return &StreamDataProducer{
		client:                   client,
		processedDataProducer:    processedDataProducer,
		systemMetricsProducer:    systemMetricsProducer,
		alertDataProducer:        alertDataProducer,
		processingResultProducer: processingResultProducer,
		config:                   conf,
	}, nil
}

// SendProcessedData sends processed data to output streams
func (p *StreamDataProducer) SendProcessedData(ctx context.Context, data model.ProcessedStreamData) error {
	// Serialize the processed data
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize processed data: %w", err)
	}

	// Create message properties
	properties := map[string]string{
		"agent_uuid": data.AgentUUID,
		"timestamp":  data.Timestamp.Format(time.RFC3339),
		"anomaly":    fmt.Sprintf("%t", data.Anomaly),
		"confidence": fmt.Sprintf("%.2f", data.Confidence),
	}

	// Send message
	_, err = p.processedDataProducer.Send(ctx, &pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	})
	if err != nil {
		return fmt.Errorf("failed to send processed data: %w", err)
	}

	return nil
}

// SendAlert sends alert data for anomalies
func (p *StreamDataProducer) SendAlert(ctx context.Context, alert model.AlertData) error {
	payload, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to serialize alert data: %w", err)
	}

	properties := map[string]string{
		"agent_uuid": alert.AgentUUID,
		"severity":   alert.Severity,
		"timestamp":  alert.Timestamp.Format(time.RFC3339),
	}

	_, err = p.alertDataProducer.Send(ctx, &pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	})
	if err != nil {
		return fmt.Errorf("failed to send alert: %w", err)
	}

	return nil
}

// SendSystemMetrics sends system metrics data
func (p *StreamDataProducer) SendSystemMetrics(ctx context.Context, metrics model.SystemMetrics) error {
	payload, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to serialize system metrics: %w", err)
	}

	properties := map[string]string{
		"agent_uuid":   metrics.AgentUUID,
		"timestamp":    metrics.Timestamp.Format(time.RFC3339),
		"cpu_usage":    fmt.Sprintf("%.2f", metrics.CPUUsage),
		"memory_usage": fmt.Sprintf("%.2f", metrics.MemoryUsage),
		"disk_usage":   fmt.Sprintf("%.2f", metrics.DiskUsage),
	}

	_, err = p.systemMetricsProducer.Send(ctx, &pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	})
	if err != nil {
		return fmt.Errorf("failed to send system metrics: %w", err)
	}

	return nil
}

// SendProcessingResult sends processing results
func (p *StreamDataProducer) SendProcessingResult(ctx context.Context, result model.StreamProcessingResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to serialize processing result: %w", err)
	}

	properties := map[string]string{
		"agent_uuid":      result.AgentUUID,
		"processing_type": result.ProcessingType,
		"success":         fmt.Sprintf("%t", result.Success),
		"timestamp":       result.Timestamp.Format(time.RFC3339),
		"processing_time": fmt.Sprintf("%d", result.ProcessingTime),
	}

	_, err = p.processingResultProducer.Send(ctx, &pulsar.ProducerMessage{
		Payload:    payload,
		Properties: properties,
	})
	if err != nil {
		return fmt.Errorf("failed to send processing result: %w", err)
	}

	return nil
}

// Close closes all producers and client
func (p *StreamDataProducer) Close() error {
	if p.processedDataProducer != nil {
		p.processedDataProducer.Close()
	}
	if p.systemMetricsProducer != nil {
		p.systemMetricsProducer.Close()
	}
	if p.alertDataProducer != nil {
		p.alertDataProducer.Close()
	}
	if p.processingResultProducer != nil {
		p.processingResultProducer.Close()
	}
	if p.client != nil {
		p.client.Close()
	}
	return nil
}

# Pulsar Configuration Guide

## Overview

Circulator uses Apache Pulsar for asynchronous message processing. Pulsar configuration has been centralized to `pkg/config` for better maintainability and flexibility.

## Configuration Structure

### YAML Configuration (`etc/app.yaml`)

```yaml
Pulsar:
  url: "pulsar://localhost:6650"
  connection_timeout: 10  # seconds
  operation_timeout: 5    # seconds
  topics:
    external_sensor_data: "external-sensor-data"
    processed_sensor_data: "processed-sensor-data"
    system_metrics: "system-metrics"
    alert_data: "alert-data"
    processing_results: "processing-results"
  consumer:
    subscription_name: "agent-processor"
    type: "Shared"  # Shared, Exclusive, Failover, KeyShared
  producer:
    send_timeout: 30  # seconds
```

### Go Configuration Structures

```go
type Pulsar struct {
    URL               string         `yaml:"url"`
    ConnectionTimeout int            `yaml:"connection_timeout"`
    OperationTimeout  int            `yaml:"operation_timeout"`
    Topics            PulsarTopics   `yaml:"topics"`
    Consumer          PulsarConsumer `yaml:"consumer"`
    Producer          PulsarProducer `yaml:"producer"`
}
```

## Topics

| Topic Name | Purpose | Data Type |
|------------|---------|-----------|
| external-sensor-data | Raw sensor data from external sources | `model.IncomingStreamData` |
| processed-sensor-data | Processed sensor data with transformations | `model.ProcessedStreamData` |
| system-metrics | System performance metrics | `model.SystemMetrics` |
| alert-data | Alert notifications for anomalies | `model.AlertData` |
| processing-results | Processing operation results | `model.StreamProcessingResult` |

## Usage

### Creating Pulsar Client

```go
client, err := config.NewPulsarClient(conf.YamlConfig)
if err != nil {
    // Handle error
}
```

### Creating Producer

```go
producer, err := NewStreamDataProducer(conf)
if err != nil {
    // Handle error
}

// Send processed data
err = producer.SendProcessedData(ctx, processedData)

// Send system metrics
err = producer.SendSystemMetrics(ctx, systemMetrics)

// Send alerts
err = producer.SendAlert(ctx, alertData)
```

### Creating Consumer

```go
consumer, err := NewStreamDataConsumer(conf)
if err != nil {
    // Handle error
}

// Consume data with handler
err = consumer.ConsumeStreamData(ctx, dataHandler)
```

## Environment-specific Configuration

### Development (Docker Compose)

```yaml
Pulsar:
  url: "pulsar://localhost:6650"
  connection_timeout: 10
  operation_timeout: 5
```

### Production

```yaml
Pulsar:
  url: "pulsar://pulsar-cluster.example.com:6650"
  connection_timeout: 30
  operation_timeout: 10
  topics:
    external_sensor_data: "prod-external-sensor-data"
    processed_sensor_data: "prod-processed-sensor-data"
    # ... other topics with prod prefixes
```

## Migration Notes

- All hardcoded Pulsar URLs have been removed
- Topic names are now configurable
- Timeout settings are now centralized
- Multiple producers are created for different topic types
- Configuration is loaded once and reused across components

## Error Handling

The system gracefully handles Pulsar connection failures:

- If Pulsar is unavailable during startup, the agent continues without Pulsar functionality
- All Pulsar operations include proper error handling and logging
- Failed message sends are logged but don't crash the application

## Docker Configuration

The system works with both standalone and cluster Pulsar deployments. The default docker-compose.yaml uses standalone mode for simplicity:

```yaml
pulsar:
  image: apachepulsar/pulsar:latest
  ports:
    - "6650:6650"  # Binary protocol
    - "8080:8080"  # Admin REST API
  command: bin/pulsar standalone
```
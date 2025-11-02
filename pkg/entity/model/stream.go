package model

import (
	"time"
)

// Pulsar-only data structures for streaming (no GORM tags)

// IncomingStreamData represents data from external sources for Pulsar streaming
type IncomingStreamData struct {
	UUID       string    `json:"uuid"`
	Source     string    `json:"source"`
	SensorType string    `json:"sensor_type"`
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
	RawPayload []byte    `json:"raw_payload"`
}

// ProcessedStreamData represents processed stream data for Pulsar streaming
type ProcessedStreamData struct {
	UUID           string    `json:"uuid"`
	AgentUUID      string    `json:"agent_uuid"`
	OriginalValue  float64   `json:"original_value"`
	ProcessedValue float64   `json:"processed_value"`
	Anomaly        bool      `json:"anomaly"`
	Confidence     float64   `json:"confidence"`
	ProcessingTime int64     `json:"processing_time"` // microseconds
	Timestamp      time.Time `json:"timestamp"`
}

// AgentProcessingConfig represents processing configuration for agents
type AgentProcessingConfig struct {
	UUID            string                   `json:"uuid"`
	AgentUUID       string                   `json:"agent_uuid"`
	SensorType      string                   `json:"sensor_type"`
	ProcessingRules []map[string]interface{} `json:"processing_rules"` // Simplified for now
	OutputStreams   []string                 `json:"output_streams"`
	Enabled         bool                     `json:"enabled"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// SystemMetrics represents real-time system metrics for Pulsar streaming
type SystemMetrics struct {
	UUID        string    `json:"uuid"`
	AgentUUID   string    `json:"agent_uuid"`
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
	Timestamp   time.Time `json:"timestamp"`
}

// AlertData represents alert information for anomalies for Pulsar streaming
type AlertData struct {
	UUID           string    `json:"uuid"`
	AgentUUID      string    `json:"agent_uuid"`
	SensorType     string    `json:"sensor_type"`
	OriginalValue  float64   `json:"original_value"`
	ProcessedValue float64   `json:"processed_value"`
	Threshold      float64   `json:"threshold"`
	Severity       string    `json:"severity"` // "low", "medium", "high", "critical"
	Message        string    `json:"message"`
	Timestamp      time.Time `json:"timestamp"`
}

// StreamProcessingResult represents the result of stream processing for Pulsar
type StreamProcessingResult struct {
	UUID           string    `json:"uuid"`
	AgentUUID      string    `json:"agent_uuid"`
	ProcessingType string    `json:"processing_type"`
	Success        bool      `json:"success"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	ProcessingTime int64     `json:"processing_time"` // microseconds
	Timestamp      time.Time `json:"timestamp"`
}

// BatchProcessingRequest represents a batch processing request for Pulsar
type BatchProcessingRequest struct {
	UUID        string               `json:"uuid"`
	AgentUUID   string               `json:"agent_uuid"`
	BatchSize   int                  `json:"batch_size"`
	DataSources []string             `json:"data_sources"`
	StreamData  []IncomingStreamData `json:"stream_data"`
	Timestamp   time.Time            `json:"timestamp"`
}

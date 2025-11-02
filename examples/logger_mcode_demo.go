package main

// Example usage of the new MCode system with proper naming convention
// Naming rule: "Package initial letters" - "Function initial letters" - "Number"
// Examples: PC-L1 (Package Config - Load - 1), SR-M2 (Server Repository - MySQL - 2)

import (
	"github.com/ryo-arima/circulator/pkg/config"
)

func demonstrateLoggerUsage() {
	// Create logger configuration
	loggerConfig := config.LoggerConfig{
		Component:    "server",
		Service:      "circulator-demo",
		Level:        "INFO",
		Structured:   true,
		EnableCaller: true,
		Output:       "stdout",
	}

	// Create base config
	baseConfig := &config.BaseConfig{}

	// Create logger
	logger := config.NewLogger(loggerConfig, baseConfig)

	// =========================
	// Configuration Layer Examples (PC - Package Config)
	// =========================

	// Config (C) examples
	logger.DEBUG(config.CL1, "")
	logger.INFO(config.CL1, "Database connection established")
	logger.WARN(config.CV2, "")
	logger.ERROR(config.CL2, "Database host: localhost") // =========================
	// Server Repository Examples (SR - Server Repository)
	// =========================

	// MySQL operation success
	logger.INFO(config.SRM1, "user created successfully", map[string]interface{}{
		"user_id":       "user-12345",
		"query_time_ms": 23,
		"table":         "users",
	})

	// MySQL operation failed
	logger.ERROR(config.SRM2, "connection timeout", map[string]interface{}{
		"host":            "mysql-cluster-01",
		"timeout_seconds": 30,
		"retry_count":     3,
	})

	// =========================
	// Agent Repository Examples (AR - Agent Repository)
	// =========================

	// Pulsar producer success
	logger.INFO(config.ARP1, "message published successfully", map[string]interface{}{
		"topic":      "circulator-events",
		"message_id": "msg-67890",
		"partition":  2,
	})

	// Agent-Server HTTP communication failed
	logger.WARN(config.ARS2, "retry scheduled", map[string]interface{}{
		"server_endpoint":     "http://circulator-server:8080/api/v1/agent/status",
		"status_code":         503,
		"retry_after_seconds": 60,
	})

	// =========================
	// File System Examples (FS - File System)
	// =========================

	// File write success
	logger.INFO(config.FSW1, "configuration backup created", map[string]interface{}{
		"backup_file": "/var/backups/circulator-config-20241102.yaml",
		"size_bytes":  2048,
	})

	// Directory creation failed
	logger.ERROR(config.FSM2, "permission denied", map[string]interface{}{
		"directory":            "/opt/circulator/data",
		"required_permissions": "755",
		"current_user":         "circulator",
	})

	// =========================
	// Communication Pattern Examples (CP - Communication Pattern)
	// =========================

	// Direct client-server communication
	logger.INFO(config.CP01, "API request completed", map[string]interface{}{
		"method":           "POST",
		"endpoint":         "/api/v1/agents",
		"response_time_ms": 120,
		"status":           "success",
	})

	// Complex end-to-end flow
	logger.INFO(config.CP23, "data processing pipeline completed", map[string]interface{}{
		"external_source":    "sensor-network-01",
		"processed_records":  1500,
		"processing_time_ms": 3200,
		"output_topic":       "processed-sensor-data",
	})

	// =========================
	// Controller Layer Examples
	// =========================

	// Server Controller - Client handling success
	logger.INFO(config.SCC1, "REST API request processed", map[string]interface{}{
		"client_ip":  "192.168.1.100",
		"user_agent": "circulator-cli/1.0.0",
		"request_id": "req-abc123",
	})

	// Agent Controller - Pulsar message handling
	logger.INFO(config.ACP1, "configuration update applied", map[string]interface{}{
		"message_type":   "config_update",
		"config_version": "v2.1.0",
		"applied_at":     "2024-11-02T12:00:00Z",
	})

	// =========================
	// Additional examples
	// =========================

	logger.INFO(config.CL1, "Demonstration completed successfully", map[string]interface{}{
		"examples_run": "all",
		"status":       "success",
	})
}

func main() {
	demonstrateLoggerUsage()
}

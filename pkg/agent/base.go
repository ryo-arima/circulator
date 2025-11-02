package agent

import (
	"context"

	"github.com/ryo-arima/circulator/pkg/agent/repository/api"
	"github.com/ryo-arima/circulator/pkg/agent/repository/local"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
	"github.com/ryo-arima/circulator/pkg/entity/request"
)

// Main handles agent operations - registers with server and starts gRPC server
func Main(conf config.BaseConfig) {
	conf.Logger.INFO(config.ABM, "Starting Agent")

	// Register agent with server
	if err := registerAgent(conf); err != nil {
		conf.Logger.FATAL(config.ABME2, "Failed to register agent", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Start gRPC server with all registered services
	if err := StartGRPCServer(conf, "50051"); err != nil {
		conf.Logger.FATAL(config.ABME3, "Failed to start gRPC server", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// registerAgent registers this agent with the server on startup
func registerAgent(conf config.BaseConfig) error {
	conf.Logger.INFO(config.ABRA, "Registering agent with server")

	// Create API client for registration
	apiClient := api.NewAPICommonRepository(conf)

	// Create local repo for system info
	localRepo := local.NewLocalDataRepository()

	// Get local system info to populate agent information
	systemInfo, err := localRepo.GetSystemInfo()
	if err != nil {
		conf.Logger.ERROR(config.ABRAE3, "Failed to get system info", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Create agent info from system and configuration
	agentInfo := model.AgentInfo{
		Hostname:       systemInfo.Hostname,
		IPAddress:      "127.0.0.1", // TODO: Get actual IP address
		Port:           50051,
		ThreadCount:    4,
		MaxThreadCount: 8,
		Version:        "1.0.0",
		Capabilities: []string{
			"stream_processing",
			"anomaly_detection",
			"system_monitoring",
		},
		Metadata: map[string]string{
			"os":           systemInfo.OS,
			"architecture": systemInfo.Architecture,
			"cpu_count":    string(rune(systemInfo.CPUCount + '0')),
		},
	}

	// Convert to RegisterAgentRequest
	registerReq := request.RegisterAgentRequest{
		UUID:           agentInfo.UUID,
		Hostname:       agentInfo.Hostname,
		IPAddress:      agentInfo.IPAddress,
		Port:           agentInfo.Port,
		ThreadCount:    agentInfo.ThreadCount,
		MaxThreadCount: agentInfo.MaxThreadCount,
		Version:        agentInfo.Version,
		Capabilities:   agentInfo.Capabilities,
		Metadata:       agentInfo.Metadata,
	}

	// Register with server using API client
	ctx := context.Background()
	_, err = apiClient.RegisterAgent(ctx, registerReq)
	if err != nil {
		conf.Logger.ERROR(config.ABRAE4, "Failed to register agent", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	conf.Logger.INFO(config.ABRAS, "Agent registration completed", map[string]interface{}{
		"hostname":   agentInfo.Hostname,
		"ip_address": agentInfo.IPAddress,
		"port":       agentInfo.Port,
		"version":    agentInfo.Version,
	})

	return nil
}

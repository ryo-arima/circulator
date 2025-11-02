package usecase

import (
	"context"
	"time"

	"github.com/ryo-arima/circulator/pkg/agent/repository/api"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
	"github.com/ryo-arima/circulator/pkg/entity/response"
)

// AgentUsecase handles stream processing business logic
type AgentUsecase struct {
	config config.BaseConfig
	repo   api.APIAgentRepository
}

// AgentRepositoryInterface defines the interface for agent data operations
type AgentRepositoryInterface interface {
	SetProcessingConfig(ctx context.Context, agentUUID string, config *model.AgentProcessingConfig) error
	GetProcessingConfig(ctx context.Context, agentUUID string) (*response.ProcessingConfigResponse, error)
}

// NewAgentUsecase creates a new AgentUsecase instance
func NewAgentUsecase(conf config.BaseConfig, repo api.APIAgentRepository) *AgentUsecase {
	return &AgentUsecase{
		config: conf,
		repo:   repo,
	}
}

// SetProcessingConfig sets the processing configuration for an agent
func (u *AgentUsecase) SetProcessingConfig(ctx context.Context, agentUUID string, processingConfig *model.AgentProcessingConfig) error {
	u.config.Logger.INFO(config.AUASPC, "Setting processing config", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.SetProcessingConfig(ctx, agentUUID, processingConfig)
}

// GetProcessingConfig gets the current processing configuration for an agent
func (u *AgentUsecase) GetProcessingConfig(ctx context.Context, agentUUID string) (*model.AgentProcessingConfig, error) {
	u.config.Logger.DEBUG(config.AUAGPC, "Getting processing config", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	resp, err := u.repo.GetProcessingConfig(ctx, agentUUID)
	if err != nil {
		return nil, err
	}
	return resp.Config, nil
}

// ProcessAgentData processes incoming stream data
func (u *AgentUsecase) ProcessAgentData(ctx context.Context, data config.IncomingAgentData) (*config.ProcessedAgentData, error) {
	startTime := time.Now()

	// Get processing configuration for this agent
	processingConfigResp, err := u.repo.GetProcessingConfig(ctx, "current-agent-uuid") // Would get from context
	if err != nil {
		return nil, err
	}

	processingConfig := processingConfigResp.Config
	if processingConfig == nil {
		// Use default processing config if none found
		processingConfig = &model.AgentProcessingConfig{
			ProcessingRules: []map[string]interface{}{},
		}
	}

	// Apply processing rules
	processedValue := u.applyProcessingRules(data.Value, processingConfig.ProcessingRules)
	isAnomaly := u.detectAnomaly(processedValue, processingConfig.ProcessingRules)
	confidence := u.calculateConfidence(processedValue, isAnomaly)

	processingTime := time.Since(startTime).Microseconds()

	result := &config.ProcessedAgentData{
		AgentUUID:      "current-agent-uuid", // Would be retrieved from context
		OriginalValue:  data.Value,
		ProcessedValue: processedValue,
		Anomaly:        isAnomaly,
		Confidence:     confidence,
		ProcessingTime: processingTime,
	}

	// Log processed data instead of storing (since StoreProcessedData doesn't exist in Repository)
	u.config.Logger.INFO(config.AUAPAD, "Processed agent data", map[string]interface{}{
		"agent_uuid":      result.AgentUUID,
		"original_value":  result.OriginalValue,
		"processed_value": result.ProcessedValue,
		"anomaly":         result.Anomaly,
		"confidence":      result.Confidence,
		"processing_time": result.ProcessingTime,
	})

	return result, nil
}

// applyProcessingRules applies configured processing rules
func (u *AgentUsecase) applyProcessingRules(value float64, rules []map[string]interface{}) float64 {
	processedValue := value

	for _, rule := range rules {
		enabled, ok := rule["enabled"].(bool)
		if !ok || !enabled {
			continue
		}

		name, ok := rule["name"].(string)
		if !ok {
			continue
		}

		switch name {
		case "moving_average":
			// Apply moving average (simplified implementation)
			if params, ok := rule["params"].(map[string]interface{}); ok {
				processedValue = u.applyMovingAverage(processedValue, params)
			}
		case "outlier_detection":
			// Outlier detection doesn't modify value, just detects
			continue
		}
	}

	return processedValue
}

// applyMovingAverage applies moving average filter
func (u *AgentUsecase) applyMovingAverage(value float64, params map[string]interface{}) float64 {
	// Simple implementation - in production would maintain buffer
	return value // Placeholder for now
}

// detectAnomaly detects anomalies based on processing rules
func (u *AgentUsecase) detectAnomaly(value float64, rules []map[string]interface{}) bool {
	for _, rule := range rules {
		enabled, ok := rule["enabled"].(bool)
		if !ok || !enabled {
			continue
		}

		name, ok := rule["name"].(string)
		if !ok || name != "outlier_detection" {
			continue
		}

		if params, ok := rule["params"].(map[string]interface{}); ok {
			if sigma, ok := params["threshold_sigma"].(float64); ok {
				// Simple threshold-based detection
				return value < (25.0-sigma*5.0) || value > (25.0+sigma*5.0)
			}
		}
	}

	// Default anomaly detection
	return value < 10.0 || value > 50.0
}

// calculateConfidence calculates confidence score
func (u *AgentUsecase) calculateConfidence(value float64, isAnomaly bool) float64 {
	if isAnomaly {
		return 0.9 // High confidence for anomalies
	}
	return 0.7 // Normal confidence for regular data
}

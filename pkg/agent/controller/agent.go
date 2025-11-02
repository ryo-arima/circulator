package controller

import (
	"context"

	"github.com/ryo-arima/circulator/pkg/agent/repository/api"
	"github.com/ryo-arima/circulator/pkg/agent/usecase"
	"github.com/ryo-arima/circulator/pkg/config"
)

// StreamController handles gRPC requests for stream processing
type StreamController struct {
	config       config.BaseConfig
	agentUsecase *usecase.AgentUsecase
}

// NewStreamController creates a new StreamController instance
func NewStreamController(conf config.BaseConfig) (*StreamController, error) {
	// Initialize API repository for usecase dependency
	apiRepo := api.NewAPIAgentRepository(conf)

	agentUsecase := usecase.NewAgentUsecase(conf, apiRepo)
	return &StreamController{
		config:       conf,
		agentUsecase: agentUsecase,
	}, nil
}

// ProcessStreamData processes incoming stream data through the agent usecase
func (c *StreamController) ProcessStreamData(ctx context.Context, data config.IncomingAgentData) (*config.ProcessedAgentData, error) {
	c.config.Logger.DEBUG(config.ACAPSD, "Processing stream data via controller", map[string]interface{}{
		"source":      data.Source,
		"sensor_type": data.SensorType,
	})

	return c.agentUsecase.ProcessAgentData(ctx, data)
}

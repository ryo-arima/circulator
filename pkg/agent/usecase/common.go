package usecase

import (
	"context"

	"github.com/ryo-arima/circulator/pkg/config"
)

// CommonUsecase handles common agent operations
type CommonUsecase struct {
	config config.BaseConfig
}

// NewCommonUsecase creates a new CommonUsecase instance
func NewCommonUsecase(conf config.BaseConfig) *CommonUsecase {
	return &CommonUsecase{
		config: conf,
	}
}

// HealthCheck performs agent health check
func (u *CommonUsecase) HealthCheck(ctx context.Context) error {
	u.config.Logger.DEBUG(config.AUCHC, "Performing agent health check")
	// Add health check logic here
	return nil
}

// GetStatus returns current agent status
func (u *CommonUsecase) GetStatus(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"status":    "running",
		"component": "agent",
		"timestamp": "now",
	}
}

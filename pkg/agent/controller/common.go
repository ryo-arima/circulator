package controller

import (
	"context"

	"github.com/ryo-arima/circulator/pkg/agent/usecase"
	"github.com/ryo-arima/circulator/pkg/config"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CommonController handles common gRPC operations for agent
type CommonController struct {
	config        config.BaseConfig
	commonUsecase *usecase.CommonUsecase
}

// NewCommonController creates a new CommonController instance
func NewCommonController(conf config.BaseConfig) *CommonController {
	return &CommonController{
		config:        conf,
		commonUsecase: usecase.NewCommonUsecase(conf),
	}
}

// GetStatus handles status request
func (c *CommonController) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	c.config.Logger.DEBUG(config.ACCGS, "Status requested via controller")
	return c.commonUsecase.GetStatus(ctx), nil
}

// Ping handles ping request
func (c *CommonController) Ping(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	c.config.Logger.DEBUG(config.ACCP, "Ping requested")
	return &emptypb.Empty{}, nil
}

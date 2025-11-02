package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/server/controller"
	"github.com/ryo-arima/circulator/pkg/server/middleware"
	"github.com/ryo-arima/circulator/pkg/server/repository"
)

func InitRouter(conf config.BaseConfig) *gin.Engine {
	conf.Logger.INFO(config.SRIR, "")

	// Initialize required repositories with config injection
	commonRepository := repository.NewCommonRepository(conf)
	agentRepository := repository.NewAgentRepository(conf)

	// Initialize required controllers with config injection
	commonController := controller.NewCommonController(conf, commonRepository)
	agentController := controller.NewAgentController(conf, agentRepository, commonRepository)

	conf.Logger.DEBUG(config.SRCARI, "", map[string]interface{}{
		"common_controller": "initialized",
		"agent_controller":  "initialized",
	})

	router := gin.Default()

	// API versioning
	v1 := router.Group("/v1")

	// Authentication endpoints
	common := v1.Group("/common")
	common.Use(middleware.Logger(conf.Logger))
	{
		conf.Logger.DEBUG(config.SRRAE, "")
		common.POST("/tokens", commonController.Login)
		common.DELETE("/tokens", commonController.Logout)
		common.GET("/tokens/validate", commonController.ValidateToken)
		common.POST("/tokens/refresh", commonController.RefreshToken)
		common.GET("/tokens/user", commonController.GetUserInfo)
	}

	// API endpoints - Authentication required for all endpoints
	v1.Use(middleware.Logger(conf.Logger), middleware.Auth())
	{
		conf.Logger.DEBUG(config.SRRPAE, "")

		// ============ AGENT ENDPOINTS ============
		v1.GET("/agents", agentController.GetAgents)
		v1.GET("/agents/count", agentController.CountAgents)
		v1.POST("/agent", agentController.CreateAgent)
		v1.PUT("/agent/:id", agentController.UpdateAgent)
		v1.DELETE("/agent/:id", agentController.DeleteAgent)

		// ============ AGENT DATA MANAGEMENT ENDPOINTS ============
		// Agent Info management
		v1.GET("/agent/:id/info", agentController.GetAgentInfo)
		v1.POST("/agent/:id/info", agentController.CreateAgentInfo)
		v1.PUT("/agent/:id/info", agentController.UpdateAgentInfo)
		v1.DELETE("/agent/:id/info", agentController.DeleteAgentInfo)

		// System Info management
		v1.GET("/agent/:id/system", agentController.GetAgentSystem)
		v1.POST("/agent/:id/system", agentController.CreateAgentSystem)
		v1.PUT("/agent/:id/system", agentController.UpdateAgentSystem)
		v1.DELETE("/agent/:id/system", agentController.DeleteAgentSystem)

		// Stream Processing Config management
		v1.GET("/agent/:id/config", agentController.GetAgentConfig)
		v1.POST("/agent/:id/config", agentController.CreateAgentConfig)
		v1.PUT("/agent/:id/config", agentController.UpdateAgentConfig)
		v1.DELETE("/agent/:id/config", agentController.DeleteAgentConfig)

		// Processing Rules management
		v1.GET("/agent/:id/config/rules", agentController.GetAgentConfigRules)
		v1.POST("/agent/:id/config/rules", agentController.CreateAgentConfigRules)
		v1.PUT("/agent/:id/config/rules/:rule_id", agentController.UpdateAgentConfigRules)
		v1.DELETE("/agent/:id/config/rules/:rule_id", agentController.DeleteAgentConfigRules)
	}

	conf.Logger.INFO(config.SRHRIS, "", map[string]interface{}{
		"endpoints_registered": "all",
		"middleware_applied":   "authentication and logging",
	})

	return router
}

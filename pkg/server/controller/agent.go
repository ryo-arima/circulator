package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
	"github.com/ryo-arima/circulator/pkg/server/repository"
	"github.com/ryo-arima/circulator/pkg/server/usecase"
)

type AgentController interface {
	GetAgents(c *gin.Context)
	CountAgents(c *gin.Context)
	GetAgent(c *gin.Context)
	CreateAgent(c *gin.Context)
	UpdateAgent(c *gin.Context)
	DeleteAgent(c *gin.Context)

	// Agent Info operations
	GetAgentInfo(c *gin.Context)
	CreateAgentInfo(c *gin.Context)
	UpdateAgentInfo(c *gin.Context)
	DeleteAgentInfo(c *gin.Context)

	// System Info operations
	GetAgentSystem(c *gin.Context)
	CreateAgentSystem(c *gin.Context)
	UpdateAgentSystem(c *gin.Context)
	DeleteAgentSystem(c *gin.Context)

	// Stream Processing Config operations
	GetAgentConfig(c *gin.Context)
	CreateAgentConfig(c *gin.Context)
	UpdateAgentConfig(c *gin.Context)
	DeleteAgentConfig(c *gin.Context)

	// Processing Rules operations
	GetAgentConfigRules(c *gin.Context)
	CreateAgentConfigRules(c *gin.Context)
	UpdateAgentConfigRules(c *gin.Context)
	DeleteAgentConfigRules(c *gin.Context)
}

type agentController struct {
	config       config.BaseConfig
	agentUsecase usecase.AgentUsecase
}

func NewAgentController(conf config.BaseConfig, agentRepo repository.AgentRepository, commonRepo repository.CommonRepository) AgentController {
	agentUsecase := usecase.NewAgentUsecase(conf, agentRepo)
	return &agentController{
		config:       conf,
		agentUsecase: agentUsecase,
	}
}

func (ctrl *agentController) GetAgents(c *gin.Context) {
	agents := ctrl.agentUsecase.GetAgents()

	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
	})
}

type counter interface {
	CountAgents() int64
}

func (ctrl *agentController) CountAgents(c *gin.Context) {
	count := ctrl.agentUsecase.CountAgents()
	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

func (ctrl *agentController) GetAgent(c *gin.Context) {
	id := c.Param("id")

	// Get basic agent by UUID
	agent := ctrl.agentUsecase.GetAgent(id)
	if agent.UUID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Agent not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"agent": agent,
	})
}

func (ctrl *agentController) CreateAgent(c *gin.Context) {
	var req request.AgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent := ctrl.agentUsecase.CreateAgent(req)

	c.JSON(http.StatusCreated, gin.H{"agent": agent})
}

func (ctrl *agentController) UpdateAgent(c *gin.Context) {
	id := c.Param("id")
	var req request.AgentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent := ctrl.agentUsecase.UpdateAgent(id, req)

	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

func (ctrl *agentController) DeleteAgent(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.agentUsecase.DeleteAgent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
} // ============ AGENT INFO OPERATIONS ============

func (ctrl *agentController) GetAgentInfo(c *gin.Context) {
	id := c.Param("id")

	agentInfo, err := ctrl.agentUsecase.GetAgentInfo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.AgentInfoResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.AgentInfo{
		ID:             agentInfo.ID,
		UUID:           agentInfo.UUID,
		Hostname:       agentInfo.Hostname,
		IPAddress:      agentInfo.IPAddress,
		Port:           agentInfo.Port,
		ThreadCount:    agentInfo.ThreadCount,
		MaxThreadCount: agentInfo.MaxThreadCount,
		Version:        agentInfo.Version,
		Capabilities:   agentInfo.Capabilities,
		Metadata:       agentInfo.Metadata,
		CreatedAt:      agentInfo.CreatedAt,
		UpdatedAt:      agentInfo.UpdatedAt,
		DeletedAt:      agentInfo.DeletedAt,
	}

	c.JSON(http.StatusOK, response.AgentInfoResponse{
		Code:    "SUCCESS",
		Message: "Agent info retrieved successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) CreateAgentInfo(c *gin.Context) {
	var req request.AgentInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentInfoResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	agentInfo, err := ctrl.agentUsecase.CreateAgentInfo(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.AgentInfoResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.AgentInfo{
		ID:             agentInfo.ID,
		UUID:           agentInfo.UUID,
		Hostname:       agentInfo.Hostname,
		IPAddress:      agentInfo.IPAddress,
		Port:           agentInfo.Port,
		ThreadCount:    agentInfo.ThreadCount,
		MaxThreadCount: agentInfo.MaxThreadCount,
		Version:        agentInfo.Version,
		Capabilities:   agentInfo.Capabilities,
		Metadata:       agentInfo.Metadata,
		CreatedAt:      agentInfo.CreatedAt,
		UpdatedAt:      agentInfo.UpdatedAt,
		DeletedAt:      agentInfo.DeletedAt,
	}

	c.JSON(http.StatusCreated, response.AgentInfoResponse{
		Code:    "SUCCESS",
		Message: "Agent info created successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) UpdateAgentInfo(c *gin.Context) {
	id := c.Param("id")
	var req request.AgentInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentInfoResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	agentInfo, err := ctrl.agentUsecase.UpdateAgentInfo(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, response.AgentInfoResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.AgentInfo{
		ID:             agentInfo.ID,
		UUID:           agentInfo.UUID,
		Hostname:       agentInfo.Hostname,
		IPAddress:      agentInfo.IPAddress,
		Port:           agentInfo.Port,
		ThreadCount:    agentInfo.ThreadCount,
		MaxThreadCount: agentInfo.MaxThreadCount,
		Version:        agentInfo.Version,
		Capabilities:   agentInfo.Capabilities,
		Metadata:       agentInfo.Metadata,
		CreatedAt:      agentInfo.CreatedAt,
		UpdatedAt:      agentInfo.UpdatedAt,
		DeletedAt:      agentInfo.DeletedAt,
	}

	c.JSON(http.StatusOK, response.AgentInfoResponse{
		Code:    "SUCCESS",
		Message: "Agent info updated successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) DeleteAgentInfo(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.agentUsecase.DeleteAgentInfo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.AgentInfoResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.AgentInfoResponse{
		Code:    "SUCCESS",
		Message: "Agent info deleted successfully",
	})
}

// ============ SYSTEM INFO OPERATIONS ============

func (ctrl *agentController) GetAgentSystem(c *gin.Context) {
	id := c.Param("id")

	systemInfo, err := ctrl.agentUsecase.GetAgentSystem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.SystemInfoResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.SystemInfo{
		ID:           systemInfo.ID,
		UUID:         systemInfo.UUID,
		AgentUUID:    systemInfo.AgentUUID,
		Hostname:     systemInfo.Hostname,
		OS:           systemInfo.OS,
		Architecture: systemInfo.Architecture,
		CPUCount:     systemInfo.CPUCount,
		Timestamp:    systemInfo.Timestamp,
		CreatedAt:    systemInfo.CreatedAt,
	}

	c.JSON(http.StatusOK, response.SystemInfoResponse{
		Code:    "SUCCESS",
		Message: "System info retrieved successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) CreateAgentSystem(c *gin.Context) {
	var req request.AgentSystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentSystemResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	// Set agent UUID from URL parameter
	req.AgentUUID = c.Param("id")

	systemInfo, err := ctrl.agentUsecase.CreateAgentSystem(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.AgentSystemResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.AgentSystem{
		ID:           systemInfo.ID,
		UUID:         systemInfo.UUID,
		AgentUUID:    systemInfo.AgentUUID,
		Hostname:     systemInfo.Hostname,
		OS:           systemInfo.OS,
		Architecture: systemInfo.Architecture,
		CPUCount:     systemInfo.CPUCount,
		Timestamp:    systemInfo.Timestamp,
		CreatedAt:    systemInfo.CreatedAt,
	}

	c.JSON(http.StatusCreated, response.AgentSystemResponse{
		Code:    "SUCCESS",
		Message: "System info created successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) UpdateAgentSystem(c *gin.Context) {
	id := c.Param("id")
	var req request.AgentSystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentSystemResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	systemInfo, err := ctrl.agentUsecase.UpdateAgentSystem(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, response.AgentSystemResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.AgentSystem{
		ID:           systemInfo.ID,
		UUID:         systemInfo.UUID,
		AgentUUID:    systemInfo.AgentUUID,
		Hostname:     systemInfo.Hostname,
		OS:           systemInfo.OS,
		Architecture: systemInfo.Architecture,
		CPUCount:     systemInfo.CPUCount,
		Timestamp:    systemInfo.Timestamp,
		CreatedAt:    systemInfo.CreatedAt,
	}

	c.JSON(http.StatusOK, response.SystemInfoResponse{
		Code:    "SUCCESS",
		Message: "System info updated successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) DeleteAgentSystem(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.agentUsecase.DeleteAgentSystem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.SystemInfoResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.SystemInfoResponse{
		Code:    "SUCCESS",
		Message: "System info deleted successfully",
	})
}

// ============ STREAM PROCESSING CONFIG OPERATIONS ============

func (ctrl *agentController) GetAgentConfig(c *gin.Context) {
	id := c.Param("id")

	config, err := ctrl.agentUsecase.GetStreamProcessingConfig(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.StreamProcessingConfigResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	// Convert ProcessingRules from model to response
	var responseRules []response.ProcessingRule
	for _, rule := range config.ProcessingRules {
		responseRules = append(responseRules, response.ProcessingRule{
			ID:        rule.ID,
			UUID:      rule.UUID,
			ConfigID:  rule.ConfigID,
			Name:      rule.Name,
			Enabled:   rule.Enabled,
			Params:    rule.Params,
			CreatedAt: rule.CreatedAt,
			UpdatedAt: rule.UpdatedAt,
			DeletedAt: rule.DeletedAt,
		})
	}

	// Convert model to response struct
	responseData := &response.StreamProcessingConfig{
		ID:              config.ID,
		UUID:            config.UUID,
		AgentUUID:       config.AgentUUID,
		SensorType:      config.SensorType,
		ProcessingRules: responseRules,
		OutputStreams:   config.OutputStreams,
		CreatedAt:       config.CreatedAt,
		UpdatedAt:       config.UpdatedAt,
		DeletedAt:       config.DeletedAt,
	}

	c.JSON(http.StatusOK, response.StreamProcessingConfigResponse{
		Code:    "SUCCESS",
		Message: "Stream processing config retrieved successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) CreateAgentConfig(c *gin.Context) {
	var req request.AgentConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentConfigResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	// Set agent UUID from URL parameter
	req.AgentUUID = c.Param("id")

	config, err := ctrl.agentUsecase.CreateStreamProcessingConfig(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.AgentConfigResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Convert ProcessingRules from model to response
	var responseRules []response.AgentConfigRules
	for _, rule := range config.ProcessingRules {
		responseRules = append(responseRules, response.AgentConfigRules{
			ID:        rule.ID,
			UUID:      rule.UUID,
			ConfigID:  rule.ConfigID,
			Name:      rule.Name,
			Enabled:   rule.Enabled,
			Params:    rule.Params,
			CreatedAt: rule.CreatedAt,
			UpdatedAt: rule.UpdatedAt,
			DeletedAt: rule.DeletedAt,
		})
	}

	// Convert model to response struct
	responseData := &response.AgentConfig{
		ID:              config.ID,
		UUID:            config.UUID,
		AgentUUID:       config.AgentUUID,
		SensorType:      config.SensorType,
		ProcessingRules: responseRules,
		OutputStreams:   config.OutputStreams,
		CreatedAt:       config.CreatedAt,
		UpdatedAt:       config.UpdatedAt,
		DeletedAt:       config.DeletedAt,
	}

	c.JSON(http.StatusCreated, response.AgentConfigResponse{
		Code:    "SUCCESS",
		Message: "Agent config created successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) UpdateAgentConfig(c *gin.Context) {
	id := c.Param("id")
	var req request.AgentConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentConfigResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	config, err := ctrl.agentUsecase.UpdateStreamProcessingConfig(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, response.AgentConfigResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	// Convert ProcessingRules from model to response
	var responseRules []response.AgentConfigRules
	for _, rule := range config.ProcessingRules {
		responseRules = append(responseRules, response.AgentConfigRules{
			ID:        rule.ID,
			UUID:      rule.UUID,
			ConfigID:  rule.ConfigID,
			Name:      rule.Name,
			Enabled:   rule.Enabled,
			Params:    rule.Params,
			CreatedAt: rule.CreatedAt,
			UpdatedAt: rule.UpdatedAt,
			DeletedAt: rule.DeletedAt,
		})
	}

	// Convert model to response struct
	responseData := &response.AgentConfig{
		ID:              config.ID,
		UUID:            config.UUID,
		AgentUUID:       config.AgentUUID,
		SensorType:      config.SensorType,
		ProcessingRules: responseRules,
		OutputStreams:   config.OutputStreams,
		CreatedAt:       config.CreatedAt,
		UpdatedAt:       config.UpdatedAt,
		DeletedAt:       config.DeletedAt,
	}

	c.JSON(http.StatusOK, response.AgentConfigResponse{
		Code:    "SUCCESS",
		Message: "Agent config updated successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) DeleteAgentConfig(c *gin.Context) {
	id := c.Param("id")

	err := ctrl.agentUsecase.DeleteStreamProcessingConfig(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.StreamProcessingConfigResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.StreamProcessingConfigResponse{
		Code:    "SUCCESS",
		Message: "Stream processing config deleted successfully",
	})
}

// ============ PROCESSING RULES OPERATIONS ============

func (ctrl *agentController) GetAgentConfigRules(c *gin.Context) {
	id := c.Param("id")

	rules, err := ctrl.agentUsecase.GetProcessingRules(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ProcessingRuleResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "Processing rules retrieved successfully",
		"data":    rules,
	})
}

func (ctrl *agentController) CreateAgentConfigRules(c *gin.Context) {
	var req request.AgentConfigRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentConfigRulesResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	id := c.Param("id")
	rule, err := ctrl.agentUsecase.CreateProcessingRule(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.AgentConfigRulesResponse{
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.AgentConfigRules{
		ID:        rule.ID,
		UUID:      rule.UUID,
		ConfigID:  rule.ConfigID,
		Name:      rule.Name,
		Enabled:   rule.Enabled,
		Params:    rule.Params,
		CreatedAt: rule.CreatedAt,
		UpdatedAt: rule.UpdatedAt,
		DeletedAt: rule.DeletedAt,
	}

	c.JSON(http.StatusCreated, response.AgentConfigRulesResponse{
		Code:    "SUCCESS",
		Message: "Agent config rules created successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) UpdateAgentConfigRules(c *gin.Context) {
	id := c.Param("id")
	ruleID := c.Param("rule_id")
	var req request.AgentConfigRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.AgentConfigRulesResponse{
			Code:    "BAD_REQUEST",
			Message: err.Error(),
		})
		return
	}

	rule, err := ctrl.agentUsecase.UpdateProcessingRule(id, ruleID, req)
	if err != nil {
		c.JSON(http.StatusNotFound, response.AgentConfigRulesResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	// Convert model to response struct
	responseData := &response.AgentConfigRules{
		ID:        rule.ID,
		UUID:      rule.UUID,
		ConfigID:  rule.ConfigID,
		Name:      rule.Name,
		Enabled:   rule.Enabled,
		Params:    rule.Params,
		CreatedAt: rule.CreatedAt,
		UpdatedAt: rule.UpdatedAt,
		DeletedAt: rule.DeletedAt,
	}

	c.JSON(http.StatusOK, response.AgentConfigRulesResponse{
		Code:    "SUCCESS",
		Message: "Agent config rules updated successfully",
		Data:    responseData,
	})
}

func (ctrl *agentController) DeleteAgentConfigRules(c *gin.Context) {
	id := c.Param("id")
	ruleID := c.Param("rule_id")

	err := ctrl.agentUsecase.DeleteProcessingRule(id, ruleID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ProcessingRuleResponse{
			Code:    "NOT_FOUND",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.ProcessingRuleResponse{
		Code:    "SUCCESS",
		Message: "Processing rule deleted successfully",
	})
}

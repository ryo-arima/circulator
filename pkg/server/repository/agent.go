package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"gorm.io/gorm"
)

type AgentRepository interface {
	GetAgents() []model.Agent
	GetAgentByUUID(uuid string) model.Agent
	CountAgents() int64
	CreateAgent(req request.AgentRequest) model.Agent
	UpdateAgent(uuid string, req request.AgentUpdateRequest) model.Agent
	DeleteAgent(uuid string) *gorm.DB

	// Agent Info operations
	GetAgentInfo(agentUUID string) (*model.AgentInfo, error)
	CreateAgentInfo(req request.AgentInfoRequest) (*model.AgentInfo, error)
	UpdateAgentInfo(agentUUID string, req request.AgentInfoRequest) (*model.AgentInfo, error)
	DeleteAgentInfo(agentUUID string) error

	// System Info operations
	GetSystemInfo(agentUUID string) (*model.SystemInfo, error)
	CreateSystemInfo(req request.AgentSystemRequest) (*model.SystemInfo, error)
	UpdateSystemInfo(agentUUID string, req request.AgentSystemRequest) (*model.SystemInfo, error)
	DeleteSystemInfo(agentUUID string) error

	// Stream Processing Config operations
	GetStreamProcessingConfig(agentUUID string) (*model.StreamProcessingConfig, error)
	CreateStreamProcessingConfig(req request.AgentConfigRequest) (*model.StreamProcessingConfig, error)
	UpdateStreamProcessingConfig(agentUUID string, req request.AgentConfigRequest) (*model.StreamProcessingConfig, error)
	DeleteStreamProcessingConfig(agentUUID string) error

	// Processing Rules operations
	GetProcessingRules(agentUUID string) ([]model.ProcessingRule, error)
	CreateProcessingRule(agentUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error)
	UpdateProcessingRule(agentUUID string, ruleUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error)
	DeleteProcessingRule(agentUUID string, ruleUUID string) error
}

type agentRepository struct {
	BaseConfig config.BaseConfig
}

func NewAgentRepository(conf config.BaseConfig) AgentRepository {
	return &agentRepository{
		BaseConfig: conf,
	}
}

func (r *agentRepository) GetAgents() []model.Agent {
	var agents []model.Agent
	r.BaseConfig.DBConnection.Find(&agents)
	return agents
}

func (r *agentRepository) GetAgentByUUID(uuid string) model.Agent {
	var agent model.Agent
	r.BaseConfig.DBConnection.Where("uuid = ?", uuid).First(&agent)
	return agent
}

func (r *agentRepository) CountAgents() int64 {
	var count int64
	r.BaseConfig.DBConnection.Model(&model.Agent{}).Count(&count)
	return count
}

func (r *agentRepository) CreateAgent(req request.AgentRequest) model.Agent {
	// Populate model from request via JSON to respect json tags, and assign UUID
	agent := model.Agent{}
	// Assign UUID (models use `UUID` field)
	agent.UUID = uuid.New().String()
	// Copy request fields into model using JSON tags
	if b, err := json.Marshal(req); err == nil {
		_ = json.Unmarshal(b, &agent)
	}
	r.BaseConfig.DBConnection.Create(&agent)
	return agent
}

func (r *agentRepository) UpdateAgent(uuid string, req request.AgentUpdateRequest) model.Agent {
	var agent model.Agent
	r.BaseConfig.DBConnection.Where("uuid = ?", uuid).First(&agent)
	// Merge request fields using JSON (partial/overwrite semantics)
	if b, err := json.Marshal(req); err == nil {
		_ = json.Unmarshal(b, &agent)
	}
	r.BaseConfig.DBConnection.Save(&agent)
	return agent
}

func (r *agentRepository) DeleteAgent(uuid string) *gorm.DB {
	return r.BaseConfig.DBConnection.Where("uuid = ?", uuid).Delete(&model.Agent{})
}

// ============ AGENT INFO OPERATIONS ============

func (r *agentRepository) GetAgentInfo(agentUUID string) (*model.AgentInfo, error) {
	var agentInfo model.AgentInfo
	result := r.BaseConfig.DBConnection.Where("uuid = ? OR agent_uuid = ?", agentUUID, agentUUID).First(&agentInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &agentInfo, nil
}

func (r *agentRepository) CreateAgentInfo(req request.AgentInfoRequest) (*model.AgentInfo, error) {
	agentInfo := &model.AgentInfo{
		UUID:           uuid.New().String(),
		Hostname:       req.Hostname,
		IPAddress:      req.IPAddress,
		Port:           req.Port,
		ThreadCount:    req.ThreadCount,
		MaxThreadCount: req.MaxThreadCount,
		Version:        req.Version,
		Capabilities:   req.Capabilities,
		Metadata:       req.Metadata,
	}

	result := r.BaseConfig.DBConnection.Create(agentInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return agentInfo, nil
}

func (r *agentRepository) UpdateAgentInfo(agentUUID string, req request.AgentInfoRequest) (*model.AgentInfo, error) {
	var agentInfo model.AgentInfo
	result := r.BaseConfig.DBConnection.Where("uuid = ?", agentUUID).First(&agentInfo)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update fields
	agentInfo.Hostname = req.Hostname
	agentInfo.IPAddress = req.IPAddress
	agentInfo.Port = req.Port
	agentInfo.ThreadCount = req.ThreadCount
	agentInfo.MaxThreadCount = req.MaxThreadCount
	agentInfo.Version = req.Version
	agentInfo.Capabilities = req.Capabilities
	agentInfo.Metadata = req.Metadata

	result = r.BaseConfig.DBConnection.Save(&agentInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &agentInfo, nil
}

func (r *agentRepository) DeleteAgentInfo(agentUUID string) error {
	result := r.BaseConfig.DBConnection.Where("uuid = ?", agentUUID).Delete(&model.AgentInfo{})
	return result.Error
}

// ============ SYSTEM INFO OPERATIONS ============

func (r *agentRepository) GetSystemInfo(agentUUID string) (*model.SystemInfo, error) {
	var systemInfo model.SystemInfo
	result := r.BaseConfig.DBConnection.Where("agent_uuid = ?", agentUUID).First(&systemInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &systemInfo, nil
}

func (r *agentRepository) CreateSystemInfo(req request.AgentSystemRequest) (*model.SystemInfo, error) {
	systemInfo := &model.SystemInfo{
		UUID:         uuid.New().String(),
		AgentUUID:    req.AgentUUID,
		Hostname:     req.Hostname,
		OS:           req.OS,
		Architecture: req.Architecture,
		CPUCount:     req.CPUCount,
		Timestamp:    time.Now(),
	}

	result := r.BaseConfig.DBConnection.Create(systemInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return systemInfo, nil
}

func (r *agentRepository) UpdateSystemInfo(agentUUID string, req request.AgentSystemRequest) (*model.SystemInfo, error) {
	var systemInfo model.SystemInfo
	result := r.BaseConfig.DBConnection.Where("agent_uuid = ?", agentUUID).First(&systemInfo)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update fields
	systemInfo.Hostname = req.Hostname
	systemInfo.OS = req.OS
	systemInfo.Architecture = req.Architecture
	systemInfo.CPUCount = req.CPUCount
	systemInfo.Timestamp = time.Now()

	result = r.BaseConfig.DBConnection.Save(&systemInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &systemInfo, nil
}

func (r *agentRepository) DeleteSystemInfo(agentUUID string) error {
	result := r.BaseConfig.DBConnection.Where("agent_uuid = ?", agentUUID).Delete(&model.SystemInfo{})
	return result.Error
}

// ============ STREAM PROCESSING CONFIG OPERATIONS ============

func (r *agentRepository) GetStreamProcessingConfig(agentUUID string) (*model.StreamProcessingConfig, error) {
	var config model.StreamProcessingConfig
	result := r.BaseConfig.DBConnection.Where("agent_uuid = ?", agentUUID).Preload("ProcessingRules").First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

func (r *agentRepository) CreateStreamProcessingConfig(req request.AgentConfigRequest) (*model.StreamProcessingConfig, error) {
	config := &model.StreamProcessingConfig{
		UUID:          uuid.New().String(),
		AgentUUID:     req.AgentUUID,
		SensorType:    req.SensorType,
		OutputStreams: req.OutputStreams,
	}

	// Create the config first
	result := r.BaseConfig.DBConnection.Create(config)
	if result.Error != nil {
		return nil, result.Error
	}

	// Create associated processing rules
	for _, ruleReq := range req.ProcessingRules {
		rule := &model.ProcessingRule{
			UUID:     uuid.New().String(),
			ConfigID: config.ID,
			Name:     ruleReq.Name,
			Enabled:  ruleReq.Enabled,
			Params:   ruleReq.Params,
		}
		r.BaseConfig.DBConnection.Create(rule)
	}

	// Reload with rules
	r.BaseConfig.DBConnection.Where("id = ?", config.ID).Preload("ProcessingRules").First(config)
	return config, nil
}

func (r *agentRepository) UpdateStreamProcessingConfig(agentUUID string, req request.AgentConfigRequest) (*model.StreamProcessingConfig, error) {
	var config model.StreamProcessingConfig
	result := r.BaseConfig.DBConnection.Where("agent_uuid = ?", agentUUID).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update fields
	config.SensorType = req.SensorType
	config.OutputStreams = req.OutputStreams

	result = r.BaseConfig.DBConnection.Save(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

func (r *agentRepository) DeleteStreamProcessingConfig(agentUUID string) error {
	// Delete associated rules first
	r.BaseConfig.DBConnection.Where("config_id IN (SELECT id FROM stream_processing_configs WHERE agent_uuid = ?)", agentUUID).Delete(&model.ProcessingRule{})
	// Delete config
	result := r.BaseConfig.DBConnection.Where("agent_uuid = ?", agentUUID).Delete(&model.StreamProcessingConfig{})
	return result.Error
}

// ============ PROCESSING RULES OPERATIONS ============

func (r *agentRepository) GetProcessingRules(agentUUID string) ([]model.ProcessingRule, error) {
	var rules []model.ProcessingRule
	result := r.BaseConfig.DBConnection.Joins("JOIN stream_processing_configs ON processing_rules.config_id = stream_processing_configs.id").Where("stream_processing_configs.agent_uuid = ?", agentUUID).Find(&rules)
	if result.Error != nil {
		return nil, result.Error
	}
	return rules, nil
}

func (r *agentRepository) CreateProcessingRule(agentUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error) {
	// Get the config first
	var config model.StreamProcessingConfig
	result := r.BaseConfig.DBConnection.Where("agent_uuid = ?", agentUUID).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}

	rule := &model.ProcessingRule{
		UUID:     uuid.New().String(),
		ConfigID: config.ID,
		Name:     req.Name,
		Enabled:  req.Enabled,
		Params:   req.Params,
	}

	result = r.BaseConfig.DBConnection.Create(rule)
	if result.Error != nil {
		return nil, result.Error
	}
	return rule, nil
}

func (r *agentRepository) UpdateProcessingRule(agentUUID string, ruleUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error) {
	var rule model.ProcessingRule
	result := r.BaseConfig.DBConnection.Joins("JOIN stream_processing_configs ON processing_rules.config_id = stream_processing_configs.id").Where("processing_rules.uuid = ? AND stream_processing_configs.agent_uuid = ?", ruleUUID, agentUUID).First(&rule)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update fields
	rule.Name = req.Name
	rule.Enabled = req.Enabled
	rule.Params = req.Params

	result = r.BaseConfig.DBConnection.Save(&rule)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rule, nil
}

func (r *agentRepository) DeleteProcessingRule(agentUUID string, ruleUUID string) error {
	result := r.BaseConfig.DBConnection.Joins("JOIN stream_processing_configs ON processing_rules.config_id = stream_processing_configs.id").Where("processing_rules.uuid = ? AND stream_processing_configs.agent_uuid = ?", ruleUUID, agentUUID).Delete(&model.ProcessingRule{})
	return result.Error
}

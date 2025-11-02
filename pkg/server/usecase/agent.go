package usecase

import (
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/server/repository"
)

type AgentUsecase interface {
	// Basic CRUD operations
	GetAgents() []model.Agent
	CountAgents() int64
	GetAgent(uuid string) model.Agent
	CreateAgent(req request.AgentRequest) model.Agent
	UpdateAgent(uuid string, req request.AgentUpdateRequest) model.Agent
	DeleteAgent(uuid string) error

	// Agent Info operations
	GetAgentInfo(agentUUID string) (*model.AgentInfo, error)
	CreateAgentInfo(req request.AgentInfoRequest) (*model.AgentInfo, error)
	UpdateAgentInfo(agentUUID string, req request.AgentInfoRequest) (*model.AgentInfo, error)
	DeleteAgentInfo(agentUUID string) error

	// System Info operations
	GetAgentSystem(agentUUID string) (*model.SystemInfo, error)
	CreateAgentSystem(req request.AgentSystemRequest) (*model.SystemInfo, error)
	UpdateAgentSystem(agentUUID string, req request.AgentSystemRequest) (*model.SystemInfo, error)
	DeleteAgentSystem(agentUUID string) error

	// Stream Processing Config operations
	GetStreamProcessingConfig(agentUUID string) (*model.StreamProcessingConfig, error)
	CreateStreamProcessingConfig(req request.AgentConfigRequest) (*model.StreamProcessingConfig, error)
	UpdateStreamProcessingConfig(agentUUID string, req request.AgentConfigRequest) (*model.StreamProcessingConfig, error)
	DeleteStreamProcessingConfig(agentUUID string) error

	// Processing Rules operations
	GetProcessingRules(agentUUID string) ([]model.ProcessingRule, error)
	CreateProcessingRule(agentUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error)
	UpdateProcessingRule(agentUUID, ruleUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error)
	DeleteProcessingRule(agentUUID, ruleUUID string) error
}

type agentUsecase struct {
	config config.BaseConfig
	repo   repository.AgentRepository
}

func NewAgentUsecase(conf config.BaseConfig, repo repository.AgentRepository) AgentUsecase {
	return &agentUsecase{
		config: conf,
		repo:   repo,
	}
}

// Basic CRUD operations
func (u *agentUsecase) GetAgents() []model.Agent {
	u.config.Logger.DEBUG(config.SUAGA, "Getting all agents")
	return u.repo.GetAgents()
}

func (u *agentUsecase) CountAgents() int64 {
	u.config.Logger.DEBUG(config.SUACA, "Counting agents")
	return u.repo.CountAgents()
}

func (u *agentUsecase) GetAgent(uuid string) model.Agent {
	u.config.Logger.DEBUG(config.SUAGA1, "Getting agent by UUID", map[string]interface{}{
		"agent_uuid": uuid,
	})
	return u.repo.GetAgentByUUID(uuid)
}

func (u *agentUsecase) CreateAgent(req request.AgentRequest) model.Agent {
	u.config.Logger.INFO(config.SUACA1, "Creating new agent", map[string]interface{}{
		"agent_name": req.Name,
	})
	return u.repo.CreateAgent(req)
}

func (u *agentUsecase) UpdateAgent(uuid string, req request.AgentUpdateRequest) model.Agent {
	u.config.Logger.INFO(config.SUAUA, "Updating agent", map[string]interface{}{
		"agent_uuid": uuid,
		"agent_name": req.Name,
	})
	return u.repo.UpdateAgent(uuid, req)
}

func (u *agentUsecase) DeleteAgent(uuid string) error {
	u.config.Logger.INFO(config.SUADA, "Deleting agent", map[string]interface{}{
		"agent_uuid": uuid,
	})
	result := u.repo.DeleteAgent(uuid)
	return result.Error
}

// Agent Info operations
func (u *agentUsecase) GetAgentInfo(agentUUID string) (*model.AgentInfo, error) {
	u.config.Logger.DEBUG(config.SUAGAI, "Getting agent info", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.GetAgentInfo(agentUUID)
}

func (u *agentUsecase) CreateAgentInfo(req request.AgentInfoRequest) (*model.AgentInfo, error) {
	u.config.Logger.INFO(config.SUACAI, "Creating agent info", map[string]interface{}{
		"hostname": req.Hostname,
	})
	return u.repo.CreateAgentInfo(req)
}

func (u *agentUsecase) UpdateAgentInfo(agentUUID string, req request.AgentInfoRequest) (*model.AgentInfo, error) {
	u.config.Logger.INFO(config.SUAUAI, "Updating agent info", map[string]interface{}{
		"agent_uuid": agentUUID,
		"hostname":   req.Hostname,
	})
	return u.repo.UpdateAgentInfo(agentUUID, req)
}

func (u *agentUsecase) DeleteAgentInfo(agentUUID string) error {
	u.config.Logger.INFO(config.SUADAI, "Deleting agent info", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.DeleteAgentInfo(agentUUID)
}

// System Info operations
func (u *agentUsecase) GetAgentSystem(agentUUID string) (*model.SystemInfo, error) {
	u.config.Logger.DEBUG(config.SUAGAS, "Getting agent system info", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.GetSystemInfo(agentUUID)
}

func (u *agentUsecase) CreateAgentSystem(req request.AgentSystemRequest) (*model.SystemInfo, error) {
	u.config.Logger.INFO(config.SUACAS, "Creating agent system info", map[string]interface{}{
		"os": req.OS,
	})
	return u.repo.CreateSystemInfo(req)
}

func (u *agentUsecase) UpdateAgentSystem(agentUUID string, req request.AgentSystemRequest) (*model.SystemInfo, error) {
	u.config.Logger.INFO(config.SUAUAS, "Updating agent system info", map[string]interface{}{
		"agent_uuid": agentUUID,
		"os":         req.OS,
	})
	return u.repo.UpdateSystemInfo(agentUUID, req)
}

func (u *agentUsecase) DeleteAgentSystem(agentUUID string) error {
	u.config.Logger.INFO(config.SUADAS, "Deleting agent system info", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.DeleteSystemInfo(agentUUID)
}

// Stream Processing Config operations
func (u *agentUsecase) GetStreamProcessingConfig(agentUUID string) (*model.StreamProcessingConfig, error) {
	u.config.Logger.DEBUG(config.SUAGSPC, "Getting stream processing config", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.GetStreamProcessingConfig(agentUUID)
}

func (u *agentUsecase) CreateStreamProcessingConfig(req request.AgentConfigRequest) (*model.StreamProcessingConfig, error) {
	u.config.Logger.INFO(config.SUACSPC, "Creating stream processing config", map[string]interface{}{
		"agent_uuid":  req.AgentUUID,
		"sensor_type": req.SensorType,
	})
	return u.repo.CreateStreamProcessingConfig(req)
}

func (u *agentUsecase) UpdateStreamProcessingConfig(agentUUID string, req request.AgentConfigRequest) (*model.StreamProcessingConfig, error) {
	u.config.Logger.INFO(config.SUAUSPC, "Updating stream processing config", map[string]interface{}{
		"agent_uuid":  agentUUID,
		"sensor_type": req.SensorType,
	})
	return u.repo.UpdateStreamProcessingConfig(agentUUID, req)
}

func (u *agentUsecase) DeleteStreamProcessingConfig(agentUUID string) error {
	u.config.Logger.INFO(config.SUADSPC, "Deleting stream processing config", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.DeleteStreamProcessingConfig(agentUUID)
}

// Processing Rules operations
func (u *agentUsecase) GetProcessingRules(agentUUID string) ([]model.ProcessingRule, error) {
	u.config.Logger.DEBUG(config.SUAGPR, "Getting processing rules", map[string]interface{}{
		"agent_uuid": agentUUID,
	})
	return u.repo.GetProcessingRules(agentUUID)
}

func (u *agentUsecase) CreateProcessingRule(agentUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error) {
	u.config.Logger.INFO(config.SUACPR, "Creating processing rule", map[string]interface{}{
		"agent_uuid": agentUUID,
		"rule_name":  req.Name,
	})
	return u.repo.CreateProcessingRule(agentUUID, req)
}

func (u *agentUsecase) UpdateProcessingRule(agentUUID, ruleUUID string, req request.AgentConfigRulesRequest) (*model.ProcessingRule, error) {
	u.config.Logger.INFO(config.SUAUPR, "Updating processing rule", map[string]interface{}{
		"agent_uuid": agentUUID,
		"rule_uuid":  ruleUUID,
		"rule_name":  req.Name,
	})
	return u.repo.UpdateProcessingRule(agentUUID, ruleUUID, req)
}

func (u *agentUsecase) DeleteProcessingRule(agentUUID, ruleUUID string) error {
	u.config.Logger.INFO(config.SUADPR, "Deleting processing rule", map[string]interface{}{
		"agent_uuid": agentUUID,
		"rule_uuid":  ruleUUID,
	})
	return u.repo.DeleteProcessingRule(agentUUID, ruleUUID)
}

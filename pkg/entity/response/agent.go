package response

import (
	"time"

	"github.com/ryo-arima/circulator/pkg/entity/model"
)

// Enveloped responses to align with locker-style outputs
type AgentResponse struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Agents  []Agent `json:"agents"`
}

type AgentListResponse = AgentResponse

type Agent struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AgentInfoResponse represents a response for agent information operations
type AgentInfoResponse struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Data    *AgentInfo `json:"data,omitempty"`
}

// AgentRegistrationResponse represents a response for agent registration operations
type AgentRegistrationResponse struct {
	Code         string            `json:"code"`
	Message      string            `json:"message"`
	RegistrationInfo map[string]string `json:"registration_info,omitempty"`
	Data    *AgentInfo  `json:"data,omitempty"`
	List    []AgentInfo `json:"list,omitempty"`
}

type AgentInfo struct {
	ID             uint              `json:"id"`
	UUID           string            `json:"uuid"`
	Hostname       string            `json:"hostname"`
	IPAddress      string            `json:"ip_address"`
	Port           int               `json:"port"`
	ThreadCount    int               `json:"thread_count"`
	MaxThreadCount int               `json:"max_thread_count"`
	Version        string            `json:"version"`
	Capabilities   []string          `json:"capabilities"`
	Metadata       map[string]string `json:"metadata"`
	CreatedAt      *time.Time        `json:"created_at"`
	UpdatedAt      *time.Time        `json:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at,omitempty"`
}

// ProcessingConfigResponse represents processing configuration response from server
type ProcessingConfigResponse struct {
	Code    string                       `json:"code"`
	Message string                       `json:"message"`
	Config  *model.AgentProcessingConfig `json:"config"`
}

// AgentSystemResponse represents a response for system information operations
type AgentSystemResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Data    *AgentSystem  `json:"data,omitempty"`
	List    []AgentSystem `json:"list,omitempty"`
}

type AgentSystem struct {
	ID           uint       `json:"id"`
	UUID         string     `json:"uuid"`
	AgentUUID    string     `json:"agent_uuid"`
	Hostname     string     `json:"hostname"`
	OS           string     `json:"os"`
	Architecture string     `json:"architecture"`
	CPUCount     int        `json:"cpu_count"`
	Timestamp    time.Time  `json:"timestamp"`
	CreatedAt    *time.Time `json:"created_at"`
}

// AgentConfigResponse represents a response for stream processing configuration operations
type AgentConfigResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Data    *AgentConfig  `json:"data,omitempty"`
	List    []AgentConfig `json:"list,omitempty"`
}

type AgentConfig struct {
	ID              uint               `json:"id"`
	UUID            string             `json:"uuid"`
	AgentUUID       string             `json:"agent_uuid"`
	SensorType      string             `json:"sensor_type"`
	ProcessingRules []AgentConfigRules `json:"processing_rules"`
	OutputStreams   []string           `json:"output_streams"`
	CreatedAt       *time.Time         `json:"created_at"`
	UpdatedAt       *time.Time         `json:"updated_at"`
	DeletedAt       *time.Time         `json:"deleted_at,omitempty"`
}

// AgentConfigRulesResponse represents a response for processing rule operations
type AgentConfigRulesResponse struct {
	Code    string             `json:"code"`
	Message string             `json:"message"`
	Data    *AgentConfigRules  `json:"data,omitempty"`
	List    []AgentConfigRules `json:"list,omitempty"`
}

type AgentConfigRules struct {
	ID        uint                   `json:"id"`
	UUID      string                 `json:"uuid"`
	ConfigID  uint                   `json:"config_id"`
	Name      string                 `json:"name"`
	Enabled   bool                   `json:"enabled"`
	Params    map[string]interface{} `json:"params"`
	CreatedAt *time.Time             `json:"created_at"`
	UpdatedAt *time.Time             `json:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at,omitempty"`
}

// Legacy type aliases for backward compatibility
type (
	SystemInfoResponse             = AgentSystemResponse
	SystemInfo                     = AgentSystem
	StreamProcessingConfigResponse = AgentConfigResponse
	StreamProcessingConfig         = AgentConfig
	ProcessingRuleResponse         = AgentConfigRulesResponse
	ProcessingRule                 = AgentConfigRules
)

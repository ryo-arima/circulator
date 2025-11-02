package request

import "time"

type AgentRequest struct {
	Hostname       string         `json:"hostname"`
	IpAddress      string         `json:"ip_address"`
	Port           int            `json:"port"`
	ThreadCount    int            `json:"thread_count"`
	MaxThreadCount int            `json:"max_thread_count"`
	Version        string         `json:"version"`
	Capabilities   []string       `json:"capabilities"`
	Metadata       map[string]any `json:"metadata"`
	UUID           string         `json:"uuid,omitempty"`
	Name           string         `json:"name,omitempty"`
	Description    string         `json:"description,omitempty"`
}

// AgentUpdateRequest for partial updates
type AgentUpdateRequest struct {
	Hostname       *string         `json:"hostname,omitempty"`
	IpAddress      *string         `json:"ip_address,omitempty"`
	Port           *int            `json:"port,omitempty"`
	ThreadCount    *int            `json:"thread_count,omitempty"`
	MaxThreadCount *int            `json:"max_thread_count,omitempty"`
	Version        *string         `json:"version,omitempty"`
	Capabilities   *[]string       `json:"capabilities,omitempty"`
	Metadata       *map[string]any `json:"metadata,omitempty"`
	Name           *string         `json:"name,omitempty"`
	Description    *string         `json:"description,omitempty"`
}

// AgentRegistrationRequest for agent registration
type AgentRegistrationRequest struct {
	AgentID      string            `json:"agent_id"`
	Hostname     string            `json:"hostname"`
	IpAddress    string            `json:"ip_address"`
	Port         int               `json:"port"`
	Version      string            `json:"version"`
	Capabilities []string          `json:"capabilities"`
	Metadata     map[string]string `json:"metadata"`
}

// AgentStatusReportRequest for agent status reporting
type AgentStatusReportRequest struct {
	AgentID     string            `json:"agent_id"`
	Status      string            `json:"status"` // "online", "offline", "busy", "idle"
	Metrics     map[string]string `json:"metrics"`
	LastUpdated time.Time         `json:"last_updated"`
}

// AgentInfoRequest represents a request for agent information operations
type AgentInfoRequest struct {
	UUID           string            `json:"uuid,omitempty"`
	Hostname       string            `json:"hostname" binding:"required"`
	IPAddress      string            `json:"ip_address" binding:"required"`
	Port           int               `json:"port" binding:"required"`
	ThreadCount    int               `json:"thread_count"`
	MaxThreadCount int               `json:"max_thread_count"`
	Version        string            `json:"version"`
	Capabilities   []string          `json:"capabilities"`
	Metadata       map[string]string `json:"metadata"`
}

// AgentSystemRequest represents a request for system information operations
type AgentSystemRequest struct {
	UUID         string `json:"uuid,omitempty"`
	AgentUUID    string `json:"agent_uuid" binding:"required"`
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	CPUCount     int    `json:"cpu_count"`
}

// HeartbeatRequest represents a heartbeat request to the server
type HeartbeatRequest struct {
	AgentUUID string    `json:"agent_uuid"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// AgentConfigRequest represents a request for stream processing configuration operations
type AgentConfigRequest struct {
	UUID            string                    `json:"uuid,omitempty"`
	AgentUUID       string                    `json:"agent_uuid" binding:"required"`
	SensorType      string                    `json:"sensor_type"`
	ProcessingRules []AgentConfigRulesRequest `json:"processing_rules"`
	OutputStreams   []string                  `json:"output_streams"`
}

// AgentConfigRulesRequest represents a request for processing rule operations
type AgentConfigRulesRequest struct {
	UUID     string                 `json:"uuid,omitempty"`
	ConfigID uint                   `json:"config_id,omitempty"`
	Name     string                 `json:"name" binding:"required"`
	Enabled  bool                   `json:"enabled"`
	Params   map[string]interface{} `json:"params"`
}

type RegisterAgentRequest struct {
	UUID           string            `json:"uuid"`
	Hostname       string            `json:"hostname"`
	IPAddress      string            `json:"ip_address"`
	Port           int               `json:"port"`
	ThreadCount    int               `json:"thread_count"`
	MaxThreadCount int               `json:"max_thread_count"`
	Version        string            `json:"version"`
	Capabilities   []string          `json:"capabilities"`
	Metadata       map[string]string `json:"metadata"`
}

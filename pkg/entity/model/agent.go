package model

import (
	"time"
)

// MySQL/GORM persistent data structures for agent information

// Agent represents basic agent information (legacy structure)
type Agent struct {
	ID             int            `gorm:"primarykey" json:"id"`
	UUID           string         `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	Hostname       string         `gorm:"type:varchar(255)" json:"hostname"`
	IpAddress      string         `gorm:"type:varchar(255)" json:"ip_address"`
	Port           int            `json:"port"`
	Status         string         `gorm:"type:varchar(255)" json:"status"`
	ThreadCount    int            `json:"thread_count"`
	MaxThreadCount int            `json:"max_thread_count"`
	Version        string         `gorm:"type:varchar(255)" json:"version"`
	Capabilities   []string       `gorm:"type:json" json:"capabilities"`
	Metadata       map[string]any `gorm:"type:json" json:"metadata"`
	HeartbeatAt    *time.Time     `gorm:"type:datetime" json:"heartbeat_at"`
	CreatedAt      *time.Time     `gorm:"type:datetime" json:"created_at"`
	UpdatedAt      *time.Time     `gorm:"type:datetime" json:"updated_at"`
}

func (Agent) TableName() string {
	return "agents"
}

// AgentInfo represents comprehensive agent information stored in MySQL
type AgentInfo struct {
	ID             uint              `gorm:"primarykey" json:"id"`
	UUID           string            `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	Hostname       string            `gorm:"type:varchar(255)" json:"hostname"`
	IPAddress      string            `gorm:"type:varchar(45)" json:"ip_address"`
	Port           int               `json:"port"`
	ThreadCount    int               `json:"thread_count"`
	MaxThreadCount int               `json:"max_thread_count"`
	Version        string            `gorm:"type:varchar(100)" json:"version"`
	Capabilities   []string          `gorm:"type:json" json:"capabilities"`
	Metadata       map[string]string `gorm:"type:json" json:"metadata"`
	CreatedAt      *time.Time        `json:"created_at"`
	UpdatedAt      *time.Time        `json:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at,omitempty"`
}

func (AgentInfo) TableName() string {
	return "agent_info"
}

// SystemInfo represents persistent system information stored in MySQL
type SystemInfo struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	UUID         string     `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	AgentUUID    string     `gorm:"type:varchar(36);index" json:"agent_uuid"`
	Hostname     string     `gorm:"type:varchar(255)" json:"hostname"`
	OS           string     `gorm:"type:varchar(100)" json:"os"`
	Architecture string     `gorm:"type:varchar(100)" json:"architecture"`
	CPUCount     int        `json:"cpu_count"`
	Timestamp    time.Time  `json:"timestamp"`
	CreatedAt    *time.Time `json:"created_at"`
}

func (SystemInfo) TableName() string {
	return "system_info"
}

// StreamProcessingConfig represents stream processing configuration stored in MySQL
type StreamProcessingConfig struct {
	ID              uint             `gorm:"primarykey" json:"id"`
	UUID            string           `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	AgentUUID       string           `gorm:"type:varchar(36);index" json:"agent_uuid"`
	SensorType      string           `gorm:"type:varchar(100)" json:"sensor_type"`
	ProcessingRules []ProcessingRule `gorm:"foreignKey:ConfigID" json:"processing_rules"`
	OutputStreams   []string         `gorm:"type:json" json:"output_streams"`
	CreatedAt       *time.Time       `json:"created_at"`
	UpdatedAt       *time.Time       `json:"updated_at"`
	DeletedAt       *time.Time       `json:"deleted_at,omitempty"`
}

func (StreamProcessingConfig) TableName() string {
	return "stream_processing_configs"
}

// ProcessingRule represents a single processing rule stored in MySQL
type ProcessingRule struct {
	ID        uint                   `gorm:"primarykey" json:"id"`
	UUID      string                 `gorm:"type:varchar(36);uniqueIndex" json:"uuid"`
	ConfigID  uint                   `gorm:"index" json:"config_id"`
	Name      string                 `gorm:"type:varchar(100)" json:"name"`
	Enabled   bool                   `json:"enabled"`
	Params    map[string]interface{} `gorm:"type:json" json:"params"`
	CreatedAt *time.Time             `json:"created_at"`
	UpdatedAt *time.Time             `json:"updated_at"`
	DeletedAt *time.Time             `json:"deleted_at,omitempty"`
}

func (ProcessingRule) TableName() string {
	return "processing_rules"
}

package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
)

// SystemRepository defines the interface for local system operations from agent
type SystemRepository interface {
	GetSystemInfo() (*model.AgentInfo, error)
	GetSystemStatus() (*model.AgentStatus, error)
	GetRegistrationInfo() (*model.Agent, error)
	StoreRegistrationInfo(agent *model.Agent) error
	Close() error
}

// systemRepository implements SystemRepository
type systemRepository struct {
	config       *config.BaseConfig
	dataFilePath string
}

// NewSystemRepository creates a new local system repository for agent
func NewSystemRepository(c *config.BaseConfig, dataDir string) SystemRepository {
	repo := &systemRepository{
		config:       c,
		dataFilePath: fmt.Sprintf("%s/agent_data.json", dataDir),
	}

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		c.Logger.ERROR(config.ALSERR, "Failed to create data directory", map[string]interface{}{
			"error":    err.Error(),
			"data_dir": dataDir,
		})
	}

	c.Logger.DEBUG(config.ALSINIT, "Agent local system repository initialized", map[string]interface{}{
		"data_dir":  dataDir,
		"data_file": repo.dataFilePath,
	})

	return repo
}

// GetSystemInfo retrieves current system information
func (r *systemRepository) GetSystemInfo() (*model.AgentInfo, error) {
	r.config.Logger.DEBUG(config.ALSGINFO, "Getting system information", nil)

	hostname, err := os.Hostname()
	if err != nil {
		r.config.Logger.ERROR(config.ALSERR, "Failed to get hostname", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	agentInfo := &model.AgentInfo{
		Hostname:       hostname,
		IPAddress:      "127.0.0.1", // Simplified, should be actual IP detection
		Port:           8080,        // Default port, should be configurable
		ThreadCount:    runtime.NumCPU(),
		MaxThreadCount: runtime.NumCPU() * 2,
		Version:        "1.0.0", // Should be from build info
		Capabilities:   []string{"monitoring", "data_processing", "alerting"},
		Metadata: map[string]string{
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"go_version": runtime.Version(),
			"started_at": time.Now().Format(time.RFC3339),
		},
	}

	r.config.Logger.DEBUG(config.ALSSUCC, "System information retrieved successfully", map[string]interface{}{
		"hostname":     hostname,
		"thread_count": agentInfo.ThreadCount,
		"max_threads":  agentInfo.MaxThreadCount,
	})

	return agentInfo, nil
}

// GetSystemStatus retrieves current system status
func (r *systemRepository) GetSystemStatus() (*model.AgentStatus, error) {
	r.config.Logger.DEBUG(config.ALSGSTAT, "Getting system status", nil)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	status := &model.AgentStatus{
		Status:      "online",
		ThreadCount: runtime.NumGoroutine(),
		Metrics: map[string]interface{}{
			"memory_alloc": memStats.Alloc,
			"memory_total": memStats.TotalAlloc,
			"memory_sys":   memStats.Sys,
			"gc_cycles":    memStats.NumGC,
			"goroutines":   runtime.NumGoroutine(),
			"last_updated": time.Now().Format(time.RFC3339),
		},
		LastUpdated: time.Now(),
	}

	r.config.Logger.DEBUG(config.ALSSUCC, "System status retrieved successfully", map[string]interface{}{
		"status":     status.Status,
		"goroutines": status.ThreadCount,
		"memory_mb":  memStats.Alloc / 1024 / 1024,
	})

	return status, nil
}

// GetRegistrationInfo retrieves stored registration information
func (r *systemRepository) GetRegistrationInfo() (*model.Agent, error) {
	r.config.Logger.DEBUG(config.ALSGREG, "Getting registration information from local storage", map[string]interface{}{
		"file_path": r.dataFilePath,
	})

	if _, err := os.Stat(r.dataFilePath); os.IsNotExist(err) {
		r.config.Logger.DEBUG(config.ALSSUCC, "No registration data found", map[string]interface{}{
			"file_path": r.dataFilePath,
		})
		return nil, nil
	}

	data, err := ioutil.ReadFile(r.dataFilePath)
	if err != nil {
		r.config.Logger.ERROR(config.ALSERR, "Failed to read registration data", map[string]interface{}{
			"error":     err.Error(),
			"file_path": r.dataFilePath,
		})
		return nil, err
	}

	var agent model.Agent
	if err := json.Unmarshal(data, &agent); err != nil {
		r.config.Logger.ERROR(config.ALSERR, "Failed to unmarshal registration data", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	r.config.Logger.DEBUG(config.ALSSUCC, "Registration information retrieved successfully", map[string]interface{}{
		"agent_uuid": agent.UUID,
		"hostname":   agent.Hostname,
	})

	return &agent, nil
}

// StoreRegistrationInfo stores registration information locally
func (r *systemRepository) StoreRegistrationInfo(agent *model.Agent) error {
	r.config.Logger.DEBUG(config.ALSSREG, "Storing registration information to local storage", map[string]interface{}{
		"agent_uuid": agent.UUID,
		"hostname":   agent.Hostname,
		"file_path":  r.dataFilePath,
	})

	data, err := json.MarshalIndent(agent, "", "  ")
	if err != nil {
		r.config.Logger.ERROR(config.ALSERR, "Failed to marshal registration data", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	if err := ioutil.WriteFile(r.dataFilePath, data, 0644); err != nil {
		r.config.Logger.ERROR(config.ALSERR, "Failed to write registration data", map[string]interface{}{
			"error":     err.Error(),
			"file_path": r.dataFilePath,
		})
		return err
	}

	r.config.Logger.DEBUG(config.ALSSUCC, "Registration information stored successfully", map[string]interface{}{
		"agent_uuid": agent.UUID,
		"file_path":  r.dataFilePath,
	})

	return nil
}

// Close cleans up the system repository
func (r *systemRepository) Close() error {
	r.config.Logger.DEBUG(config.ALSSUCC, "Agent local system repository closed", nil)
	return nil
}

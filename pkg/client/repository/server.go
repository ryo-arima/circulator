package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/model"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
)

// ServerRepository handles HTTP communication with Server
type ServerRepository struct {
	config  config.BaseConfig
	baseURL string
	client  *http.Client
}

// NewServerRepository creates a new ServerRepository instance
func NewServerRepository(config config.BaseConfig, baseURL string) *ServerRepository {
	return &ServerRepository{
		config:  config,
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Login authenticates with the server and returns tokens
func (r *ServerRepository) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	r.config.Logger.DEBUG(config.CRSLOGIN, "Attempting server login", map[string]interface{}{
		"email": req.Email,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal login request: %w", err)
	}

	url := r.baseURL + "/v1/login"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(httpReq)
	if err != nil {
		r.config.Logger.ERROR(config.CRSERR, "Server login request failed", map[string]interface{}{
			"error": err.Error(),
			"url":   url,
		})
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var loginResp response.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	r.config.Logger.INFO(config.CRSSUCC, "Server login successful", nil)
	return &loginResp, nil
}

// GetAllAgents retrieves all agents from the server
func (r *ServerRepository) GetAllAgents(ctx context.Context) ([]*model.Agent, error) {
	r.config.Logger.DEBUG(config.CRSGA, "Getting all agents from server", nil)

	url := r.baseURL + "/v1/agents"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	resp, err := r.client.Do(httpReq)
	if err != nil {
		r.config.Logger.ERROR(config.CRSERR, "Get all agents request failed", map[string]interface{}{
			"error": err.Error(),
			"url":   url,
		})
		return nil, fmt.Errorf("get all agents request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get all agents failed with status %d: %s", resp.StatusCode, string(body))
	}

	var agents []*model.Agent
	if err := json.Unmarshal(body, &agents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal agents response: %w", err)
	}

	r.config.Logger.INFO(config.CRSSUCC, "Successfully retrieved agents from server", map[string]interface{}{
		"agent_count": len(agents),
	})
	return agents, nil
}

// CreateAgent creates a new agent on the server
func (r *ServerRepository) CreateAgent(ctx context.Context, agent *model.Agent) (*model.Agent, error) {
	r.config.Logger.INFO(config.CRSCA, "Creating agent on server", map[string]interface{}{
		"agent_uuid":     agent.UUID,
		"agent_hostname": agent.Hostname,
	})

	jsonData, err := json.Marshal(agent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent: %w", err)
	}

	url := r.baseURL + "/v1/agents"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(httpReq)
	if err != nil {
		r.config.Logger.ERROR(config.CRSERR, "Create agent request failed", map[string]interface{}{
			"error": err.Error(),
			"url":   url,
		})
		return nil, fmt.Errorf("create agent request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create agent failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdAgent model.Agent
	if err := json.Unmarshal(body, &createdAgent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal created agent response: %w", err)
	}

	r.config.Logger.INFO(config.CRSSUCC, "Successfully created agent on server", map[string]interface{}{
		"agent_uuid": createdAgent.UUID,
	})
	return &createdAgent, nil
}

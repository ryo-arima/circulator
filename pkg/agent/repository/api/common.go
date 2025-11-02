package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
)

// APICommonRepository handles API communications with server
type APICommonRepository interface {
	Login(ctx context.Context, email, password string) (*response.LoginResponse, error)
	RegisterAgent(ctx context.Context, req request.RegisterAgentRequest) (*response.RegisterAgentResponse, error)
	SendHeartbeat(ctx context.Context, req request.HeartbeatRequest) error
	ValidateToken(ctx context.Context) bool
	RefreshToken(ctx context.Context) error
}

type apiCommonRepository struct {
	config config.BaseConfig
}

// NewAPICommonRepository creates a new API repository with HTTP client
func NewAPICommonRepository(conf config.BaseConfig) APICommonRepository {
	repo := &apiCommonRepository{
		config: conf,
	}
	return repo
}

// getBaseURL returns the base URL from configuration
func (r *apiCommonRepository) getBaseURL() string {
	baseURL := r.config.YamlConfig.Application.Agent.ServerEndpoint
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return baseURL
}

// getHTTPClient creates a new HTTP client
func (r *apiCommonRepository) getHTTPClient() *http.Client {
	return &http.Client{}
}

func (r *apiCommonRepository) Login(ctx context.Context, email, password string) (*response.LoginResponse, error) {
	r.config.Logger.INFO(config.ARACLOG, "Attempting login", map[string]interface{}{"email": email})
	payload := map[string]string{"email": email, "password": password}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", r.getBaseURL()+"/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.getHTTPClient().Do(req)
	if err != nil {
		r.config.Logger.ERROR(config.ARACLOG, "Login request failed", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed: %s", string(body))
	}
	var loginResp response.LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	r.config.Logger.INFO(config.ARACLOG, "Login successful", nil)
	return &loginResp, nil
}

func (r *apiCommonRepository) RegisterAgent(ctx context.Context, req request.RegisterAgentRequest) (*response.RegisterAgentResponse, error) {
	r.config.Logger.INFO(config.ARACREG, "Registering agent", map[string]interface{}{"uuid": req.UUID})
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", r.getBaseURL()+"/agents/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.getHTTPClient().Do(httpReq)
	if err != nil {
		r.config.Logger.ERROR(config.ARACREG, "Agent registration failed", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("agent registration failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("registration failed: %s", string(body))
	}
	var regResp response.RegisterAgentResponse
	if err := json.Unmarshal(body, &regResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	r.config.Logger.INFO(config.ARACREG, "Agent registration successful", nil)
	return &regResp, nil
}

func (r *apiCommonRepository) SendHeartbeat(ctx context.Context, req request.HeartbeatRequest) error {
	r.config.Logger.DEBUG(config.ARACHB, "Sending heartbeat", map[string]interface{}{"agent_uuid": req.AgentUUID})
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", r.getBaseURL()+"/agents/heartbeat", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := r.getHTTPClient().Do(httpReq)
	if err != nil {
		r.config.Logger.ERROR(config.ARACHB, "Heartbeat failed", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("heartbeat failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("heartbeat failed: %s", string(body))
	}
	r.config.Logger.DEBUG(config.ARACHB, "Heartbeat successful", nil)
	return nil
}

func (r *apiCommonRepository) ValidateToken(ctx context.Context) bool {
	// 簡素化: 常にtrueを返す（認証が不要な場合）
	return true
}

func (r *apiCommonRepository) RefreshToken(ctx context.Context) error {
	// 簡素化: 何もしない（認証が不要な場合）
	r.config.Logger.DEBUG(config.ARACRT, "Token refresh not implemented", nil)
	return nil
}

package api

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

// APIAgentRepository interface defines all agent-related operations
type APIAgentRepository interface {
	// Basic agent operations
	GetAgents(ctx context.Context) (*response.AgentResponse, error)
	CountAgents(ctx context.Context) (*response.AgentResponse, error)
	CreateAgent(ctx context.Context, req request.AgentRequest) (*response.AgentResponse, error)
	UpdateAgent(ctx context.Context, id int, req request.AgentRequest) (*response.AgentResponse, error)
	DeleteAgent(ctx context.Context, id int) (*response.AgentResponse, error)

	// Agent Info management
	GetAgentInfo(ctx context.Context, id int) (*response.AgentResponse, error)
	CreateAgentInfo(ctx context.Context, id int, req request.AgentInfoRequest) (*response.AgentResponse, error)
	UpdateAgentInfo(ctx context.Context, id int, req request.AgentInfoRequest) (*response.AgentResponse, error)
	DeleteAgentInfo(ctx context.Context, id int) (*response.AgentResponse, error)

	// System Info management
	GetAgentSystem(ctx context.Context, id int) (*response.AgentSystemResponse, error)
	CreateAgentSystem(ctx context.Context, id int, req request.AgentSystemRequest) (*response.AgentSystemResponse, error)
	UpdateAgentSystem(ctx context.Context, id int, req request.AgentSystemRequest) (*response.AgentSystemResponse, error)
	DeleteAgentSystem(ctx context.Context, id int) (*response.AgentSystemResponse, error)

	// Stream Processing Config management
	GetAgentConfig(ctx context.Context, id int) (*response.AgentConfigResponse, error)
	CreateAgentConfig(ctx context.Context, id int, req request.AgentConfigRequest) (*response.AgentConfigResponse, error)
	UpdateAgentConfig(ctx context.Context, id int, req request.AgentConfigRequest) (*response.AgentConfigResponse, error)
	DeleteAgentConfig(ctx context.Context, id int) (*response.AgentConfigResponse, error)

	// Processing Rules management
	GetAgentConfigRules(ctx context.Context, id int) (*response.AgentConfigRulesResponse, error)
	CreateAgentConfigRules(ctx context.Context, id int, req request.AgentConfigRulesRequest) (*response.AgentConfigRulesResponse, error)
	UpdateAgentConfigRules(ctx context.Context, id int, ruleID int, req request.AgentConfigRulesRequest) (*response.AgentConfigRulesResponse, error)
	DeleteAgentConfigRules(ctx context.Context, id int, ruleID int) (*response.AgentConfigRulesResponse, error)

	// Processing Configuration legacy methods
	SetProcessingConfig(ctx context.Context, agentUUID string, config *model.AgentProcessingConfig) error
	GetProcessingConfig(ctx context.Context, agentUUID string) (*response.ProcessingConfigResponse, error)
}

type apiAgentRepository struct {
	config     config.BaseConfig
	repository APICommonRepository
	baseURL    string
}

// NewAPIAgentRepository creates a new API agent repository
func NewAPIAgentRepository(conf config.BaseConfig) APIAgentRepository {
	baseURL := conf.YamlConfig.Application.Agent.ServerEndpoint
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default endpoint
	}

	return &apiAgentRepository{
		config:     conf,
		repository: NewAPICommonRepository(conf),
		baseURL:    baseURL,
	}
}

// ============ AGENT ENDPOINTS ============

// GetAgents gets all agents (GET /agents)
func (r *apiAgentRepository) GetAgents(ctx context.Context) (*response.AgentResponse, error) {
	r.config.Logger.DEBUG(config.AREGA, "Getting all agents", nil)

	url := r.baseURL + "/v1/agents"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CountAgents counts agents (GET /agents/count)
func (r *apiAgentRepository) CountAgents(ctx context.Context) (*response.AgentResponse, error) {
	r.config.Logger.DEBUG(config.ARECA, "Counting agents", nil)

	url := r.baseURL + "/v1/agents/count"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateAgent creates a new agent (POST /agent)
func (r *apiAgentRepository) CreateAgent(ctx context.Context, req request.AgentRequest) (*response.AgentResponse, error) {
	r.config.Logger.INFO(config.ARECRA, "Creating agent", map[string]interface{}{
		"agent_id": req.UUID,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + "/v1/agent"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateAgent updates an existing agent (PUT /agent/:id)
func (r *apiAgentRepository) UpdateAgent(ctx context.Context, id int, req request.AgentRequest) (*response.AgentResponse, error) {
	r.config.Logger.INFO(config.AREUA, "Updating agent", map[string]interface{}{
		"id":       id,
		"agent_id": req.UUID,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d", id)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make PUT request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteAgent deletes an agent (DELETE /agent/:id)
func (r *apiAgentRepository) DeleteAgent(ctx context.Context, id int) (*response.AgentResponse, error) {
	r.config.Logger.INFO(config.AREDA, "Deleting agent", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d", id)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============ AGENT INFO MANAGEMENT ENDPOINTS ============

// GetAgentInfo gets agent information (GET /agent/:id/info)
func (r *apiAgentRepository) GetAgentInfo(ctx context.Context, id int) (*response.AgentResponse, error) {
	r.config.Logger.DEBUG(config.AREGAI, "Getting agent info", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/info", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateAgentInfo creates agent information (POST /agent/:id/info)
func (r *apiAgentRepository) CreateAgentInfo(ctx context.Context, id int, req request.AgentInfoRequest) (*response.AgentResponse, error) {
	r.config.Logger.INFO(config.ARECRAI, "Creating agent info", map[string]interface{}{
		"id": id,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/info", id)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateAgentInfo updates agent information (PUT /agent/:id/info)
func (r *apiAgentRepository) UpdateAgentInfo(ctx context.Context, id int, req request.AgentInfoRequest) (*response.AgentResponse, error) {
	r.config.Logger.INFO(config.AREUAI, "Updating agent info", map[string]interface{}{
		"id": id,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/info", id)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make PUT request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteAgentInfo deletes agent information (DELETE /agent/:id/info)
func (r *apiAgentRepository) DeleteAgentInfo(ctx context.Context, id int) (*response.AgentResponse, error) {
	r.config.Logger.INFO(config.AREDAI, "Deleting agent info", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/info", id)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============ SYSTEM INFO MANAGEMENT ENDPOINTS ============

// GetAgentSystem gets agent system information (GET /agent/:id/system)
func (r *apiAgentRepository) GetAgentSystem(ctx context.Context, id int) (*response.AgentSystemResponse, error) {
	r.config.Logger.DEBUG(config.AREGAS, "Getting agent system info", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/system", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentSystemResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateAgentSystem creates agent system information (POST /agent/:id/system)
func (r *apiAgentRepository) CreateAgentSystem(ctx context.Context, id int, req request.AgentSystemRequest) (*response.AgentSystemResponse, error) {
	r.config.Logger.INFO(config.ARECRAS, "Creating agent system info", map[string]interface{}{
		"id": id,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/system", id)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentSystemResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateAgentSystem updates agent system information (PUT /agent/:id/system)
func (r *apiAgentRepository) UpdateAgentSystem(ctx context.Context, id int, req request.AgentSystemRequest) (*response.AgentSystemResponse, error) {
	r.config.Logger.INFO(config.AREUAS, "Updating agent system info", map[string]interface{}{
		"id": id,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/system", id)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make PUT request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentSystemResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteAgentSystem deletes agent system information (DELETE /agent/:id/system)
func (r *apiAgentRepository) DeleteAgentSystem(ctx context.Context, id int) (*response.AgentSystemResponse, error) {
	r.config.Logger.INFO(config.AREDAS, "Deleting agent system info", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/system", id)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentSystemResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============ STREAM PROCESSING CONFIG MANAGEMENT ENDPOINTS ============

// GetAgentConfig gets agent processing configuration (GET /agent/:id/config)
func (r *apiAgentRepository) GetAgentConfig(ctx context.Context, id int) (*response.AgentConfigResponse, error) {
	r.config.Logger.DEBUG(config.AREGAC, "Getting agent config", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateAgentConfig creates agent processing configuration (POST /agent/:id/config)
func (r *apiAgentRepository) CreateAgentConfig(ctx context.Context, id int, req request.AgentConfigRequest) (*response.AgentConfigResponse, error) {
	r.config.Logger.INFO(config.ARECRAC, "Creating agent config", map[string]interface{}{
		"id": id,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config", id)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateAgentConfig updates agent processing configuration (PUT /agent/:id/config)
func (r *apiAgentRepository) UpdateAgentConfig(ctx context.Context, id int, req request.AgentConfigRequest) (*response.AgentConfigResponse, error) {
	r.config.Logger.INFO(config.AREUAC, "Updating agent config", map[string]interface{}{
		"id": id,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config", id)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make PUT request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteAgentConfig deletes agent processing configuration (DELETE /agent/:id/config)
func (r *apiAgentRepository) DeleteAgentConfig(ctx context.Context, id int) (*response.AgentConfigResponse, error) {
	r.config.Logger.INFO(config.AREDAC, "Deleting agent config", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config", id)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============ PROCESSING RULES MANAGEMENT ENDPOINTS ============

// GetAgentConfigRules gets agent processing rules (GET /agent/:id/config/rules)
func (r *apiAgentRepository) GetAgentConfigRules(ctx context.Context, id int) (*response.AgentConfigRulesResponse, error) {
	r.config.Logger.DEBUG(config.AREGACR, "Getting agent config rules", map[string]interface{}{
		"id": id,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config/rules", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigRulesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CreateAgentConfigRules creates agent processing rules (POST /agent/:id/config/rules)
func (r *apiAgentRepository) CreateAgentConfigRules(ctx context.Context, id int, req request.AgentConfigRulesRequest) (*response.AgentConfigRulesResponse, error) {
	r.config.Logger.INFO(config.ARECRACR, "Creating agent config rules", map[string]interface{}{
		"id": id,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config/rules", id)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigRulesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateAgentConfigRules updates specific agent processing rule (PUT /agent/:id/config/rules/:rule_id)
func (r *apiAgentRepository) UpdateAgentConfigRules(ctx context.Context, id int, ruleID int, req request.AgentConfigRulesRequest) (*response.AgentConfigRulesResponse, error) {
	r.config.Logger.INFO(config.AREUACR, "Updating agent config rules", map[string]interface{}{
		"id":      id,
		"rule_id": ruleID,
	})

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config/rules/%d", id, ruleID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make PUT request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigRulesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteAgentConfigRules deletes specific agent processing rule (DELETE /agent/:id/config/rules/:rule_id)
func (r *apiAgentRepository) DeleteAgentConfigRules(ctx context.Context, id int, ruleID int) (*response.AgentConfigRulesResponse, error) {
	r.config.Logger.INFO(config.AREDACR, "Deleting agent config rules", map[string]interface{}{
		"id":      id,
		"rule_id": ruleID,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/%d/config/rules/%d", id, ruleID)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result response.AgentConfigRulesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ============ LEGACY PROCESSING CONFIG METHODS ============

// SetProcessingConfig sets processing configuration for an agent (legacy method)
func (r *apiAgentRepository) SetProcessingConfig(ctx context.Context, agentUUID string, processingConfig *model.AgentProcessingConfig) error {
	r.config.Logger.INFO(config.ARESPC, "Setting processing config", map[string]interface{}{
		"uuid": agentUUID,
	})

	jsonData, err := json.Marshal(processingConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	url := r.baseURL + fmt.Sprintf("/v1/agent/uuid/%s/config", agentUUID)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to make PUT request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to set processing config: %s", string(body))
	}

	r.config.Logger.INFO(config.ARESPC, "Processing config set successfully", nil)
	return nil
}

// GetProcessingConfig retrieves processing configuration for an agent (legacy method)
func (r *apiAgentRepository) GetProcessingConfig(ctx context.Context, agentUUID string) (*response.ProcessingConfigResponse, error) {
	r.config.Logger.DEBUG(config.AREGPC, "Getting processing config", map[string]interface{}{
		"uuid": agentUUID,
	})

	url := r.baseURL + fmt.Sprintf("/v1/agent/uuid/%s/config", agentUUID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get processing config: %s", string(body))
	}

	var result response.ProcessingConfigResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

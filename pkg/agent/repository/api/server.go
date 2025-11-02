package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
)

// ServerRepository defines the interface for server API operations from agent
type ServerRepository interface {
	Login(req *request.LoginRequest) (*response.LoginResponse, error)
	GetAgentInfo(agentID string) (*response.AgentInfoResponse, error)
	SendStatusReport(req *request.AgentStatusReportRequest) (*response.CommonResponse, error)
	GetRegistrationInfo(agentID string) (*response.AgentRegistrationResponse, error)
	SendRegistration(req *request.AgentRegistrationRequest) (*response.CommonResponse, error)
	Close() error
}

// serverRepository implements ServerRepository
type serverRepository struct {
	config     *config.BaseConfig
	serverURL  string
	httpClient *http.Client
	authToken  string
}

// NewServerRepository creates a new server API repository for agent
func NewServerRepository(c *config.BaseConfig, serverURL string) ServerRepository {
	repo := &serverRepository{
		config:     c,
		serverURL:  serverURL,
		httpClient: &http.Client{},
	}
	
	c.Logger.DEBUG(config.ARSINIT, "Agent Server API repository initialized", map[string]interface{}{
		"server_url": serverURL,
	})
	
	return repo
}

// Login authenticates the agent with the server
func (r *serverRepository) Login(req *request.LoginRequest) (*response.LoginResponse, error) {
	r.config.Logger.DEBUG(config.ARSLOGIN, "Agent attempting login to server", map[string]interface{}{
		"email": req.Email,
	})

	body, err := json.Marshal(req)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to marshal login request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	resp, err := r.httpClient.Post(
		fmt.Sprintf("%s/api/v1/login", r.serverURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to send login request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to read login response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var loginResp response.LoginResponse
	if err := json.Unmarshal(responseBody, &loginResp); err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to unmarshal login response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	// Store auth token for subsequent requests
	r.authToken = loginResp.Token

	r.config.Logger.DEBUG(config.ARSSUCC, "Agent login successful", map[string]interface{}{
		"email": req.Email,
	})

	return &loginResp, nil
}

// GetAgentInfo retrieves agent information from the server
func (r *serverRepository) GetAgentInfo(agentID string) (*response.AgentInfoResponse, error) {
	r.config.Logger.DEBUG(config.ARSGINFO, "Agent getting info from server", map[string]interface{}{
		"agent_id": agentID,
	})

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/agents/%s", r.serverURL, agentID), nil)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to create get agent info request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	if r.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.authToken)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to send get agent info request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to read agent info response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var agentInfoResp response.AgentInfoResponse
	if err := json.Unmarshal(responseBody, &agentInfoResp); err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to unmarshal agent info response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	r.config.Logger.DEBUG(config.ARSSUCC, "Agent info retrieved successfully", map[string]interface{}{
		"agent_id": agentID,
	})

	return &agentInfoResp, nil
}

// SendStatusReport sends agent status report to the server
func (r *serverRepository) SendStatusReport(req *request.AgentStatusReportRequest) (*response.CommonResponse, error) {
	r.config.Logger.DEBUG(config.ARSSPREP, "Agent sending status report to server", map[string]interface{}{
		"agent_id": req.AgentID,
		"status":   req.Status,
	})

	body, err := json.Marshal(req)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to marshal status report request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/agents/%s/status", r.serverURL, req.AgentID), bytes.NewBuffer(body))
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to create status report request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if r.authToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+r.authToken)
	}

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to send status report request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to read status report response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var commonResp response.CommonResponse
	if err := json.Unmarshal(responseBody, &commonResp); err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to unmarshal status report response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	r.config.Logger.DEBUG(config.ARSSUCC, "Agent status report sent successfully", map[string]interface{}{
		"agent_id": req.AgentID,
	})

	return &commonResp, nil
}

// GetRegistrationInfo retrieves agent registration information from the server
func (r *serverRepository) GetRegistrationInfo(agentID string) (*response.AgentRegistrationResponse, error) {
	r.config.Logger.DEBUG(config.ARSGREG, "Agent getting registration info from server", map[string]interface{}{
		"agent_id": agentID,
	})

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/agents/%s/registration", r.serverURL, agentID), nil)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to create get registration info request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	if r.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.authToken)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to send get registration info request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to read registration info response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var regResp response.AgentRegistrationResponse
	if err := json.Unmarshal(responseBody, &regResp); err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to unmarshal registration info response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	r.config.Logger.DEBUG(config.ARSSUCC, "Agent registration info retrieved successfully", map[string]interface{}{
		"agent_id": agentID,
	})

	return &regResp, nil
}

// SendRegistration sends agent registration to the server
func (r *serverRepository) SendRegistration(req *request.AgentRegistrationRequest) (*response.CommonResponse, error) {
	r.config.Logger.DEBUG(config.ARSSREG, "Agent sending registration to server", map[string]interface{}{
		"agent_id": req.AgentID,
		"hostname": req.Hostname,
	})

	body, err := json.Marshal(req)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to marshal registration request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/agents/%s/register", r.serverURL, req.AgentID), bytes.NewBuffer(body))
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to create registration request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if r.authToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+r.authToken)
	}

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to send registration request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to read registration response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var commonResp response.CommonResponse
	if err := json.Unmarshal(responseBody, &commonResp); err != nil {
		r.config.Logger.ERROR(config.ARSERR, "Failed to unmarshal registration response", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	r.config.Logger.DEBUG(config.ARSSUCC, "Agent registration sent successfully", map[string]interface{}{
		"agent_id": req.AgentID,
	})

	return &commonResp, nil
}

// Close cleans up the server repository
func (r *serverRepository) Close() error {
	r.config.Logger.DEBUG(config.ARSSUCC, "Agent Server API repository closed", nil)
	return nil
}
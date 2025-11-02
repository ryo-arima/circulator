package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
)

type AgentRepository interface {
	BootstrapAgentForDB(req request.AgentRequest) interface{}
	GetAgent(req request.AgentRequest) interface{}
	CreateAgent(req request.AgentRequest) interface{}
	UpdateAgent(req request.AgentRequest) interface{}
	DeleteAgent(req request.AgentRequest) interface{}
}

type agentRepository struct {
	BaseConfig config.BaseConfig
}

func NewAgentRepository(conf config.BaseConfig) AgentRepository {
	return &agentRepository{BaseConfig: conf}
}

func (r *agentRepository) BootstrapAgentForDB(req request.AgentRequest) interface{} {
	return map[string]any{"code": "SUCCESS", "message": "Bootstrap for Agent completed successfully"}
}

func (r *agentRepository) GetAgent(req request.AgentRequest) interface{} {
	url := fmt.Sprintf("%s/v1/agents", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint)
	client := &http.Client{}
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	bearer(httpReq)
	resp, err := client.Do(httpReq)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	var out response.AgentListResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	return out
}

func (r *agentRepository) CreateAgent(req request.AgentRequest) interface{} {
	url := fmt.Sprintf("%s/v1/agent", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint)
	b, err := json.Marshal(req)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	bearer(httpReq)
	resp, err := client.Do(httpReq)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	var out response.AgentResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	return out
}

func (r *agentRepository) UpdateAgent(req request.AgentRequest) interface{} {
	id := extractUUID(req)
	if id == "" {
		return map[string]any{"code": "error", "message": "update requires uuid; not provided in request"}
	}
	url := fmt.Sprintf("%s/v1/agent/%s", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, id)
	b, err := json.Marshal(req)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	client := &http.Client{}
	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	bearer(httpReq)
	resp, err := client.Do(httpReq)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	var out response.AgentResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	return out
}

func (r *agentRepository) DeleteAgent(req request.AgentRequest) interface{} {
	id := extractUUID(req)
	if id == "" {
		return map[string]any{"code": "error", "message": "delete requires uuid; not provided in request"}
	}
	url := fmt.Sprintf("%s/v1/agent/%s", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint, id)
	client := &http.Client{}
	httpReq, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	bearer(httpReq)
	resp, err := client.Do(httpReq)
	if err != nil {
		return map[string]any{"code": "error", "message": err.Error()}
	}
	defer resp.Body.Close()
	return map[string]any{"code": "SUCCESS", "message": "deleted"}
}

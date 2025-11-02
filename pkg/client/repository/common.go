package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
)

type CommonRepository interface {
	// Authentication methods
	Login(req request.LoginRequest) (response.LoginResponse, error)
	RefreshToken(req request.RefreshTokenRequest) (response.RefreshTokenResponse, error)
	ValidateToken(token string) (response.ValidateResponse, error)
	Logout(token string) (response.CommonResponse, error)
	GetUserInfo(token string) (response.CommonResponse, error)
}

type commonRepository struct {
	BaseConfig config.BaseConfig
}

func NewCommonRepository(conf config.BaseConfig) CommonRepository {
	return &commonRepository{
		BaseConfig: conf,
	}
}

// --- shared helpers for repository package ---
// loadAccessTokenFromFiles tries token files used by locker-style clients
func loadAccessTokenFromFiles() string {
	paths := []string{
		filepath.Join("etc", ".circulator", "client", "base", "access_token"),
		filepath.Join("etc", ".circulator", "client", "app", "access_token"),
	}
	for _, p := range paths {
		if b, err := os.ReadFile(p); err == nil && len(b) > 0 {
			return strings.TrimSpace(string(b))
		}
	}
	return ""
}

// bearer sets Authorization header from env or token files
func bearer(req *http.Request) {
	token := os.Getenv("STREAM_MANAGER_ACCESS_TOKEN")
	if token == "" {
		token = loadAccessTokenFromFiles()
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// extractUUID tries to pull `uuid` from any request struct via JSON tags
func extractUUID(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return ""
	}
	if u, ok := m["uuid"].(string); ok {
		return u
	}
	return ""
}

func (r *commonRepository) Login(req request.LoginRequest) (response.LoginResponse, error) {
	url := fmt.Sprintf("%s/v1/common/tokens", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return response.LoginResponse{}, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return response.LoginResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.LoginResponse{}, err
	}

	var result response.LoginResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return response.LoginResponse{}, err
	}

	return result, nil
}

func (r *commonRepository) RefreshToken(req request.RefreshTokenRequest) (response.RefreshTokenResponse, error) {
	url := fmt.Sprintf("%s/v1/common/tokens/refresh", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return response.RefreshTokenResponse{}, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return response.RefreshTokenResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.RefreshTokenResponse{}, err
	}

	var result response.RefreshTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return response.RefreshTokenResponse{}, err
	}

	return result, nil
}

func (r *commonRepository) ValidateToken(token string) (response.ValidateResponse, error) {
	url := fmt.Sprintf("%s/v1/common/tokens/validate", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response.ValidateResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return response.ValidateResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.ValidateResponse{}, err
	}

	var result response.ValidateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return response.ValidateResponse{}, err
	}

	return result, nil
}

func (r *commonRepository) Logout(token string) (response.CommonResponse, error) {
	url := fmt.Sprintf("%s/v1/common/tokens", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return response.CommonResponse{}, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return response.CommonResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.CommonResponse{}, err
	}

	var result response.CommonResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return response.CommonResponse{}, err
	}

	return result, nil
}

func (r *commonRepository) GetUserInfo(token string) (response.CommonResponse, error) {
	url := fmt.Sprintf("%s/v1/common/tokens/user", r.BaseConfig.YamlConfig.Application.Client.ServerEndpoint)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response.CommonResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return response.CommonResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.CommonResponse{}, err
	}

	var result response.CommonResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return response.CommonResponse{}, err
	}

	return result, nil
}

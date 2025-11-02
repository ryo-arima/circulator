package usecase

import (
	"fmt"

	"github.com/ryo-arima/circulator/pkg/client/repository"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
)

type CommonUsecase interface {
	// Authentication methods
	Login(req request.LoginRequest) response.LoginResponse
	RefreshToken(refreshToken string) response.RefreshTokenResponse
	Logout(accessToken string) response.CommonResponse
	ValidateToken(accessToken string) response.ValidateResponse
	GetUserInfo(accessToken string) response.CommonResponse
}

type commonUsecase struct {
	config config.BaseConfig
	repo   repository.CommonRepository
}

// NewCommonUsecase creates a new CommonUsecase using config and underlying repository
func NewCommonUsecase(conf config.BaseConfig) CommonUsecase {
	return &commonUsecase{
		config: conf,
		repo:   repository.NewCommonRepository(conf),
	}
}

func (u *commonUsecase) Login(req request.LoginRequest) response.LoginResponse {
	result, err := u.repo.Login(req)
	if err != nil {
		return response.LoginResponse{
			Code:    "error",
			Message: fmt.Sprintf("Login failed: %v", err),
		}
	}
	return result
}

func (u *commonUsecase) RefreshToken(refreshToken string) response.RefreshTokenResponse {
	req := request.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	result, err := u.repo.RefreshToken(req)
	if err != nil {
		return response.RefreshTokenResponse{
			Code:    "error",
			Message: fmt.Sprintf("Token refresh failed: %v", err),
		}
	}
	return result
}

func (u *commonUsecase) Logout(accessToken string) response.CommonResponse {
	// Implementation would call logout endpoint
	// For now, return success
	return response.CommonResponse{
		Code:    "success",
		Message: "Logout successful",
	}
}

func (u *commonUsecase) ValidateToken(accessToken string) response.ValidateResponse {
	result, err := u.repo.ValidateToken(accessToken)
	if err != nil {
		return response.ValidateResponse{
			Code:    "error",
			Message: fmt.Sprintf("Token validation failed: %v", err),
			Valid:   false,
		}
	}
	return result
}

func (u *commonUsecase) GetUserInfo(accessToken string) response.CommonResponse {
	// Delegate to repository implementation
	res, err := u.repo.GetUserInfo(accessToken)
	if err != nil {
		return response.CommonResponse{Code: "error", Message: fmt.Sprintf("Get user info failed: %v", err)}
	}
	return res
}

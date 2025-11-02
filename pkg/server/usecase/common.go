package usecase

import (
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
	"github.com/ryo-arima/circulator/pkg/entity/response"
	"github.com/ryo-arima/circulator/pkg/server/repository"
)

type CommonUsecase interface {
	ValidateUser(email, password string) bool
	GenerateToken(email string) (string, error)
	ValidateToken(token string) bool
	LoginUser(req request.UserRequest) response.LoginResponse
	RefreshUser(req request.UserRequest) response.RefreshTokenResponse
}

type commonUsecase struct {
	config config.BaseConfig
	repo   repository.CommonRepository
}

func NewCommonUsecase(conf config.BaseConfig, repo repository.CommonRepository) CommonUsecase {
	return &commonUsecase{
		config: conf,
		repo:   repo,
	}
}

func (u *commonUsecase) ValidateUser(email, password string) bool {
	u.config.Logger.DEBUG(config.SUCVU, "Validating user credentials", map[string]interface{}{
		"email": email,
	})
	return u.repo.ValidateUser(email, password)
}

func (u *commonUsecase) GenerateToken(email string) (string, error) {
	u.config.Logger.DEBUG(config.SUCGT, "Generating token for user", map[string]interface{}{
		"email": email,
	})
	return u.repo.GenerateToken(email)
}

func (u *commonUsecase) ValidateToken(token string) bool {
	u.config.Logger.DEBUG(config.SUCVT, "Validating token")
	return u.repo.ValidateToken(token)
}

func (u *commonUsecase) LoginUser(req request.UserRequest) response.LoginResponse {
	u.config.Logger.INFO(config.SUCLU, "Processing user login", map[string]interface{}{
		"email": req.Email,
	})
	return u.repo.LoginUser(req)
}

func (u *commonUsecase) RefreshUser(req request.UserRequest) response.RefreshTokenResponse {
	u.config.Logger.INFO(config.SUCRU, "Processing token refresh", map[string]interface{}{
		"email": req.Email,
	})
	return u.repo.RefreshUser(req)
}

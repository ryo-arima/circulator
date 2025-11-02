package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/server/repository"
	"github.com/ryo-arima/circulator/pkg/server/usecase"
)

type CommonController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	ValidateToken(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetUserInfo(c *gin.Context)
}

type commonController struct {
	config        config.BaseConfig
	commonUsecase usecase.CommonUsecase
}

func NewCommonController(conf config.BaseConfig, commonRepo repository.CommonRepository) CommonController {
	commonUsecase := usecase.NewCommonUsecase(conf, commonRepo)
	return &commonController{
		config:        conf,
		commonUsecase: commonUsecase,
	}
}

func (ctrl *commonController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Implement authentication logic
	token := "dummy-jwt-token"

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}

func (ctrl *commonController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (ctrl *commonController) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}

func (ctrl *commonController) RefreshToken(c *gin.Context) {
	// Implement refresh token logic
	newToken := "new-dummy-jwt-token"

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
		"type":  "Bearer",
	})
}

// GetUserInfo returns minimal user info from the token/context
func (ctrl *commonController) GetUserInfo(c *gin.Context) {
	// In real implementation, extract user info from JWT or session
	c.JSON(http.StatusOK, gin.H{
		"uuid":  "00000000-0000-0000-0000-000000000000",
		"email": "anonymous@example.com",
		"name":  "Anonymous",
	})
}

package usecase

import (
	"github.com/ryo-arima/circulator/pkg/client/repository"
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/entity/request"
)

type AgentUsecase interface {
	Bootstrap(req request.AgentRequest, format string) string
	Get(req request.AgentRequest, format string) string
	Create(req request.AgentRequest, format string) string
	Update(req request.AgentRequest, format string) string
	Delete(req request.AgentRequest, format string) string
}

type agentUsecase struct {
	config config.BaseConfig
	repo   repository.AgentRepository
}

func NewAgentUsecase(conf config.BaseConfig) AgentUsecase {
	return &agentUsecase{
		config: conf,
		repo:   repository.NewAgentRepository(conf),
	}
}

func (u *agentUsecase) Bootstrap(req request.AgentRequest, format string) string {
	resp := u.repo.BootstrapAgentForDB(req)
	return Format(format, resp)
}

func (u *agentUsecase) Get(req request.AgentRequest, format string) string {
	resp := u.repo.GetAgent(req)
	return Format(format, resp)
}

func (u *agentUsecase) Create(req request.AgentRequest, format string) string {
	resp := u.repo.CreateAgent(req)
	return Format(format, resp)
}

func (u *agentUsecase) Update(req request.AgentRequest, format string) string {
	resp := u.repo.UpdateAgent(req)
	return Format(format, resp)
}

func (u *agentUsecase) Delete(req request.AgentRequest, format string) string {
	resp := u.repo.DeleteAgent(req)
	return Format(format, resp)
}

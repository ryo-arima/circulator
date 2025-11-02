package main

import (
	"github.com/ryo-arima/circulator/pkg/agent"
	"github.com/ryo-arima/circulator/pkg/config"
)

func main() {
	conf := config.NewAgentConfig()
	// Base operations with DB access
	agent.Main(conf)
}

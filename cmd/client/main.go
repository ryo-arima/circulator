package main

import (
	"github.com/ryo-arima/circulator/pkg/client"
	"github.com/ryo-arima/circulator/pkg/config"
)

func main() {
	conf := config.NewBaseConfig()
	// Base operations with DB access
	client.Client(conf)
}

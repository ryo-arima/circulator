package main

import (
	"github.com/ryo-arima/circulator/pkg/config"
	"github.com/ryo-arima/circulator/pkg/server"
)

func main() {
	conf := config.NewBaseConfig()
	server.Main(conf)
}

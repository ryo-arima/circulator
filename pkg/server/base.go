package server

import (
	"github.com/ryo-arima/circulator/pkg/config"
)

func Main(conf config.BaseConfig) {
	router := InitRouter(conf)
	router.Run(":" + conf.YamlConfig.Application.Common.Port)
}

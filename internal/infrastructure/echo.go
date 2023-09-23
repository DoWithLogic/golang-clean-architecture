package infrastructure

import (
	"github.com/DoWithLogic/golang-clean-architecture/config"
	"github.com/labstack/echo/v4"
)

func NewEcho(cfg config.ServerConfig) (*echo.Echo, error) {
	echo := echo.New()
	echo.Debug = cfg.Debug

	return echo, nil
}

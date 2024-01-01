package http

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"gharabakloo/search/internal/entity"
	"gharabakloo/search/pkg/myerr"
)

type Runner interface {
	Run(addr ...string) (err error)
}

func RunServer(runner Runner, cfg entity.HTTPConfig) error {
	if err := validateConfig(cfg); err != nil {
		return myerr.Errorf(err)
	}

	addr := fmt.Sprintf("%s:%s", cfg.IP, cfg.Port)
	err := runner.Run(addr)
	return myerr.Errorf(err)
}

func validateConfig(cfg entity.HTTPConfig) error {
	_, err := strconv.Atoi(cfg.Port)
	return myerr.Errorf(err)
}

func SetupHTTPRouter(h *Handler) *gin.Engine {
	engine := gin.Default()
	engine.GET("/search", h.Search)
	return engine
}

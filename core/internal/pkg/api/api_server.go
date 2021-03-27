package api

import (
	srv "core/internal/pkg/server"
	"strconv"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

// Config API server configuration
type Config struct {
	Host    string
	APIPort int
	Verbose bool
}

// New initialize a new gin-based HTTP server
func New(sctx *srv.Ctx, cfg Config) {
	if !cfg.Verbose {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	if cfg.Verbose {
		r.Use(logger.SetLogger(logger.Config{
			Logger: sctx.Log.Logger,
		}))
	}

	ApplyRoutes(sctx, r)

	err := r.Run(cfg.Host + ":" + strconv.Itoa(cfg.APIPort))
	if err != nil {
		sctx.Log.Error.Str("reason", err.Error()).Msg("Unable to start HTTP server")
	}
}

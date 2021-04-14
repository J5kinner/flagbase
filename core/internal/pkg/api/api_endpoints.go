package api

import (
	"core/internal/app/access"
	"core/internal/app/environment"
	"core/internal/app/flag"
	"core/internal/app/healthcheck"
	"core/internal/app/identity"
	"core/internal/app/project"
	"core/internal/app/segment"
	"core/internal/app/segmentrule"
	"core/internal/app/trait"
	"core/internal/app/variation"
	"core/internal/app/workspace"
	srv "core/internal/pkg/server"

	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies route from all packages to root handler
func ApplyRoutes(sctx *srv.Ctx, r *gin.Engine) {
	ApplyMetrics(r)
	root := r.Group("/")
	access.ApplyRoutes(sctx, root)
	environment.ApplyRoutes(sctx, root)
	flag.ApplyRoutes(sctx, root)
	healthcheck.ApplyRoutes(sctx, root)
	identity.ApplyRoutes(sctx, root)
	project.ApplyRoutes(sctx, root)
	trait.ApplyRoutes(sctx, root)
	segment.ApplyRoutes(sctx, root)
	segmentrule.ApplyRoutes(sctx, root)
	variation.ApplyRoutes(sctx, root)
	workspace.ApplyRoutes(sctx, root)
}
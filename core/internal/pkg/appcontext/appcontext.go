package appcontext

import (
	"core//pkg/logger"
	"core/internal/pkg/policy"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Ctx primary app context structure
type Ctx struct {
	Cache  *redis.Client
	DB     *pgxpool.Pool
	Log    *logger.Logger
	Policy *policy.Policy
	Metric string // TODO: add metric interface
}

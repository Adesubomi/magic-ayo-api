package data

import (
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Repo struct {
	Config      *configPkg.Config
	DbClient    *gorm.DB
	RedisClient *redis.Client
}

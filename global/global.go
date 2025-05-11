package global

import (
	"github.com/HeRedBo/pkg/cache"
	"github.com/HeRedBo/pkg/es"
	"github.com/HeRedBo/pkg/nosql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ES    *es.Client
	LOG   *zap.Logger
	DB    *gorm.DB
	CACHE *cache.Redis
	Mongo *nosql.MgClient
)

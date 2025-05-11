package main

import (
	"github.com/HeRedBo/pkg/cache"
	"github.com/HeRedBo/pkg/db"
	"github.com/HeRedBo/pkg/es"
	"github.com/HeRedBo/pkg/logger"
	"github.com/HeRedBo/pkg/nosql"
	"github.com/HeRedBo/pkg/timeutil"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"shop-search-api/config"
	"shop-search-api/global"
)

func init() {
	config.LoadConfig()
	InitLog()
	initMysqlClient()
	initRedisClient()
	initMongoClient()
	initESClient()
}

func InitLog() {
	// 初始化 logger
	global.LOG = logger.InitLogger(
		//logger.WithDisableConsole(),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileRotationP(config.Cfg.App.AppLogPath),
	)
}

func initMysqlClient() {
	mysqlCfg := config.Cfg.Mysql
	err := db.InitMysqlClient(db.DefaultClient, mysqlCfg.User, mysqlCfg.Password, mysqlCfg.Host, mysqlCfg.DBName)
	if err != nil {
		global.LOG.Error("mysql init error", zap.Error(err))
		panic("initMysqlClient error")
	}
	global.DB = db.GetMysqlClient(db.DefaultClient).DB
}

func initRedisClient() {
	redisCfg := config.Cfg.Redis
	opt := redis.Options{
		Addr:         redisCfg.Host,
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		MaxRetries:   redisCfg.MaxRetries,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConn,
	}

	err := cache.InitRedis(cache.DefaultRedisClient, &opt)
	if err != nil {
		global.LOG.Error("redis init error", zap.Error(err))
		panic("initRedisClient error")
	}
	global.CACHE = cache.GetRedisClient(cache.DefaultRedisClient)
}

func initESClient() {
	ESCfg := config.Cfg.Elasticsearch
	err := es.InitClientWithOptions(es.DefaultClient, ESCfg.Host,
		ESCfg.User,
		ESCfg.Password,
		es.WithScheme("https"))
	if err != nil {
		logger.Error("InitClientWithOptions error", zap.Error(err), zap.String("client", es.DefaultClient))
		panic(err)
	}
	global.ES = es.GetClient(es.DefaultClient)
}

func initMongoClient() {
	err := nosql.InitMongoClient(nosql.DefaultMongoClient, config.Cfg.MongoDB.User,
		config.Cfg.MongoDB.Password, config.Cfg.MongoDB.Host, 200)
	if err != nil {
		logger.Error("InitMongoClient error", zap.Error(err), zap.String("client", nosql.DefaultMongoClient))
		//panic(err)
	} else {
		global.Mongo = nosql.GeMongoClient(nosql.DefaultMongoClient)
	}
}

func main() {

}

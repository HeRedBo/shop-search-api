package main

import (
	"context"
	"fmt"
	"github.com/HeRedBo/pkg/cache"
	"github.com/HeRedBo/pkg/db"
	"github.com/HeRedBo/pkg/es"
	"github.com/HeRedBo/pkg/logger"
	"github.com/HeRedBo/pkg/nosql"
	"github.com/HeRedBo/pkg/shutdown"
	"github.com/HeRedBo/pkg/timeutil"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"net/http"
	"shop-search-api/config"
	"shop-search-api/global"
	"shop-search-api/internal/server/api"
	"time"
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
		es.WithScheme("https"),
		es.WithQueryLogEnable(true), // 启用查询日志
	)
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
	router := api.InitRouter()
	listenAddr := fmt.Sprintf(":%d", config.Cfg.App.HttpPort)
	global.LOG.Warn("start http server", zap.String("listenAddr", listenAddr))
	server := &http.Server{
		Addr:           listenAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			global.LOG.Error("http server start error", zap.Error(err))
		}
	}()

	//优雅关闭
	shutdown.NewHook().Close(
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				global.LOG.Error("http server shutdown err", zap.Error(err))
			}
		},

		func() {
			es.CloseAll()
		},

		func() {
			//关闭mysql
			if err := db.CloseMysqlClient(db.DefaultClient); err != nil {
				global.LOG.Error("mysql shutdown err", zap.Error(err), zap.String("client", db.DefaultClient))
			}
		},

		func() {
			err := global.CACHE.Close()
			if err != nil {
				global.LOG.Error("redis close error", zap.Error(err), zap.String("client", cache.DefaultRedisClient))
			}
		},
		func() {
			if global.Mongo != nil {
				global.Mongo.Close()
			}
		},
		func() {
			err := global.LOG.Sync()
			if err != nil {
				fmt.Println(err)
			}
		},
	)

}

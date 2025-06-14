package api

import (
	"github.com/gin-gonic/gin"
	v1 "shop-search-api/internal/server/api/v1"
)

func InitRouter() *gin.Engine {
	engin := gin.New()
	engin.Use(gin.Logger())
	//防止panic发生，返回500
	engin.Use(gin.Recovery())
	engin.HEAD("/health", Health)

	//通过中间件进行接口签名校验
	apiv1 := engin.Group("/api/v1")
	//apiv1.Use(auth.Auth())
	apiv1.GET("/product-search", v1.ProductSearch)
	apiv1.GET("/order-search", v1.OrderSearch)

	return engin
}

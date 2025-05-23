package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"shop-search-api/global"
	"shop-search-api/internal/pkg/errcode"
	"shop-search-api/internal/repo/es/product_repo"
	"shop-search-api/internal/server/api/api_response"
	"shop-search-api/internal/service/product_service"
	"strconv"
)

type searchResponse struct {
	Total int64                        `json:"total"`
	Hits  []*product_repo.ProductIndex `json:"hits"`
}

func ProductSearch(c *gin.Context) {
	// TODO 监控上报
	appG := api_response.Gin{C: c}
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		appG.ResponseErr(errcode.ErrCodes.ErrParams)
		return
	}

	productService := product_service.Product{
		Keyword:  keyword,
		PageNum:  com.StrTo(c.Query("page_num")).MustInt(),
		PageSize: com.StrTo(c.Query("page_size")).MustInt(),
		UserID:   com.StrTo(c.Query("userid")).MustInt64(),
		Sales:    c.Query("sales_order"),
		Price:    c.Query("price_order"),
	}

	news := c.Query("news")
	newsFilter := com.StrTo(news).MustInt()
	if len(news) == 0 {
		productService.New = nil
	} else {
		productService.New = &newsFilter
	}
	res, err := productService.SearchProduct()
	global.LOG.Warn("resp", zap.Any("", res))
	if err != nil {
		global.LOG.Error("search error", zap.Error(err), zap.Any("param", productService))
		appG.ResponseErr(errcode.ErrCodes.ErrSearch)
		return
	}
	resp := searchResponse{
		Total: 0,
		Hits:  make([]*product_repo.ProductIndex, 0),
	}
	if res == nil {
		appG.ResponseOk(errcode.ErrCodes.ErrNo, resp)
		return
	}
	resp.Total = res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		index := &product_repo.ProductIndex{}
		err = json.Unmarshal(hit.Source, index)
		if err != nil {
			global.LOG.Error("Unmarshal error", zap.Error(err))
			continue
		}
		index.Id, err = strconv.ParseInt(hit.Id, 10, 64)
		if err != nil {
			global.LOG.Error("strconv.ParseInt error", zap.Error(err), zap.String("id", hit.Id))
			continue
		}
		resp.Hits = append(resp.Hits, index)
	}

	//global.LOG.Warn("resp", zap.Any("resp", resp))
	appG.ResponseOk(errcode.ErrCodes.ErrNo, resp)

}

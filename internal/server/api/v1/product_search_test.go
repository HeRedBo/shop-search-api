package v1

import (
	"fmt"
	"github.com/HeRedBo/pkg/httpclient"
	"github.com/gookit/goutil/dump"
	"net/http"
	"net/url"
	"shop-search-api/config"
	"shop-search-api/internal/pkg/sign"
	"testing"
	"time"
)

const ProductSearchHost = "http://127.0.0.1:9091"
const ProductSearchUri = "/api/v1/product-search"

var (
	ak  = "AK100523687952"
	sk  = "W1WTYvJpfeH1YpUjTpeFbEx^DnpQ&35L"
	ttl = time.Minute * 3
)

func TestProductSearch(t *testing.T) {
	params := url.Values{}
	params.Add("userid", "1")
	params.Add("keyword", "手机")
	params.Add("page_num", "1")
	params.Add("page_size", "10")
	authorization, date, err := sign.New(ak, sk, ttl).Generate(ProductSearchUri, http.MethodGet, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	dump.P(authorization, date, err)
	headerAuth := httpclient.WithHeader(config.HeaderAuthField, authorization)
	headerAuthDate := httpclient.WithHeader(config.HeaderAuthDateField, date)
	c, r, e := httpclient.Get(ProductSearchHost+ProductSearchUri, params, headerAuth, headerAuthDate)
	fmt.Println(c, string(r), e)
}

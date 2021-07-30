package handler

import (
	"ZeroProject/Api/unifiedLogin/internal/svc"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 统一登录
var (
	url  = "http://10.28.83.128:8083/login/unifiedLogin"
	body = bytes.NewBufferString(`{"username": "PD038", "password": "F379EAF3C831B04DE153469D1BEC345E"}`)
)

// 统一登录 单元测试
func TestHandler(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, url, body)
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	var c svc.Config
	c = svc.ClientConfig(c)
	ctx := svc.NewServiceContext(c)
	router := UnifiedLoginHandler(ctx)

	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
	assert.Equal(t, 200, rr.Code)
	fmt.Println(rr.Body.String())
}

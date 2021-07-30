package handler

import (
	"ZeroProject/Api/unifiedLogin/internal/svc"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	token      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJObyI6IjIxMDYxODA1IiwiY29tcGFueV9pZCI6IkMwMDEiLCJleHBpcmUiOjE2MjUxMjg2NTAsImlhdCI6MTYyNTEyNTA1MCwiaWQiOiJVc2VyMTQzNTciLCJ1c2VyTmFtZSI6ImplcnJ5In0.QxmhXG3oW_PCRIBKawQQrpMWvO9EfPtpKk2WJ0L0XwU"
	resFeshUrl = "http://10.28.83.123:8896/to_ken/refresh"
)

// token 刷新 单元测试
func TestRefresh(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, resFeshUrl, nil)
	assert.Nil(t, err)
	r.Header.Set("Authorization", "Bearer "+token)

	var c svc.Config
	c = svc.ClientConfig(c)
	ctx := svc.NewServiceContext(c)
	router := AuthTokenRefresh(ctx)

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
	assert.Equal(t, 200, rr.Code)
	fmt.Println(rr.Body.String())
}

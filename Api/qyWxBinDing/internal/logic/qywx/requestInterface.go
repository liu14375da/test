package qywx

import (
	"ZeroProject/Api/qyWxBinDing/internal/model"
	"ZeroProject/common/global"
	"ZeroProject/common/tool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

/*
	请求企业微信的接口
*/

type AccessToken struct {
	AccessToken string    `json:"access_token"`         // 获取到的凭证，最长为512字节
	ExpiresIn   int64     `json:"expires_in,omitempty"` // 凭证的有效时间（秒），通常为2小时（7200秒）
	ExpireAt    time.Time `json:"expire_at,omitempty"`  // 过期时间，超过时重新获取
}

// AuthCodeURL 生成网页授权地址.
//  appId:       公众号的唯一标识
//  redirectURI: 授权后重定向的回调链接地址
//  scope:       应用授权作用域
//  state:       重定向后会带上 state 参数, 开发者可以填写 a-zA-Z0-9 的参数值, 最多128字节
func AuthCodeURL(appId, redirectURI, scope, state string) string {
	return "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + url.QueryEscape(appId) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
}

// 获取访问用户身份,根据code获取成员信息
func AccessTokenData(cache1, code string) string {
	getTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=" + cache1 + "&code=" + code
	client := &http.Client{}
	req, _ := client.Get(getTokenUrl)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	var jsonStr model.Userinfo
	json.Unmarshal(body, &jsonStr)
	userId := jsonStr.UserId
	fmt.Println("userId:", userId)
	return userId
}

// 获取access_token
func GetAccessToken() string {
	cache1 := tool.GetCache("title")
	if cache1 == "" {
		getTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + global.QyWxRpc.Corpid + "&corpsecret=" + global.QyWxRpc.Corpsecret
		client := &http.Client{}
		req, _ := client.Get(getTokenUrl)
		defer req.Body.Close()
		body, _ := ioutil.ReadAll(req.Body)
		var jsonStr AccessToken
		json.Unmarshal(body, &jsonStr)
		token := jsonStr.AccessToken
		tool.SetCache(token, "title")
		return token
	} else {
		return cache1
	}
}

// 成员添加到企业微信
func AddUser(userid, accessToken string) *model.User {
	getTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/user/authsucc?access_token=" + accessToken + "&userid=" + userid
	client := &http.Client{}
	req, _ := client.Get(getTokenUrl)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	var jsonStr model.User
	json.Unmarshal(body, &jsonStr)
	return &jsonStr
}

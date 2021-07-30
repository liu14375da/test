package logic

import (
	"ZeroProject/Rpc/qyWxBinDing/internal/model/qyWx"
	"ZeroProject/Rpc/qyWxBinDing/internal/svc"
	pb "ZeroProject/Rpc/qyWxBinDing/pb"
	"ZeroProject/common/errorx"
	"ZeroProject/common/global"
	"ZeroProject/common/jwt"
	"ZeroProject/common/tool"
	"time"

	"context"
	"encoding/json"
	"errors"
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"net/http"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func (l PingLogic) ObtainQyWx(request *pb.Request) (*pb.Response, error) {
	userInfo, err := l.svcCtx.BandingWx.SelectUser(request.ErpName)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if userInfo.UserPassword != request.Pwd {
		return nil, errors.New("用户密码不正确")
	}
	newAccessToken, err := GetAccessToken()
	if err != nil {
		return nil, errors.New("获取企业微信参数失败")
	}
	user := GetUser(request.WxName, newAccessToken)
	if user.UserID == "" {
		return nil, errors.New("绑定时请退出企业微信,在进行登录绑定")
	}
	qyRow, _ := l.svcCtx.BandingWx.SelectUserId(user.UserID)
	if qyRow == nil || qyRow.User_Id == "" {
		number, err := insertUser(l, user, userInfo.StaffId, request.Cron)
		if err != nil {
			return nil, errors.New("插入数据失败")
		}
		return &pb.Response{
			Msg: number,
		}, nil
	}
	return nil, errors.New("不可重复绑定")
	//return &pb.Response{
	//	Msg: "不可重复绑定",
	//}, nil
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func GetAccessToken() (string, error) {
	if cache1 := tool.GetCache("title"); cache1 == "" {
		getTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + global.QyWxRpc.Corpid + "&corpsecret=" + global.QyWxRpc.Corpsecret
		client := &http.Client{}
		req, _ := client.Get(getTokenUrl)
		defer req.Body.Close()
		body, _ := ioutil.ReadAll(req.Body)
		var jsonStr qyWx.AccessToken
		json.Unmarshal(body, &jsonStr)
		token := jsonStr.AccessToken
		return token, nil
	} else {
		return cache1, nil
	}
}

func GetUser(name, token string) *qyWx.User {
	getTokenUrl := "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=" + token + "&userid=" + name
	client := &http.Client{}
	req, _ := client.Get(getTokenUrl)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	var jsonStr qyWx.User
	json.Unmarshal(body, &jsonStr)
	return &jsonStr
}

// 这块加个删除的原因是因为当我们第一次绑定完，把对应的数据添加到数据库中，如果管理员在企业微信平台把我们删掉
//我们后台是不知道那个删除动作的，这里添加个删除功能是建立在我们绑定之后，删除，在绑定的情景下的，
//因为那个微信ueserid会重新生成的，所以我们把数据库的的删掉再次添加新的数据
func insertUser(l PingLogic, user *qyWx.User, id string, cron string) (string, error) {
	if staffId, err := l.svcCtx.BandingWx.SelectStaffID(user.UserID, id); err == nil {
		switch staffId {
		case 1:
			return "绑定成功!", nil
		case 2:
			qyRow, err := l.svcCtx.BandingWx.InsertUser(user.UserID, user.Name, user.Mobile, user.Email,
				user.Alias, user.Address, user.Position, user.Gender, id, user.Telephone,
				user.Enable, user.ToInvite, user.Status, cron)
			if err != nil {
				return "", err
			}
			if qyRow != 0 {
				return "绑定成功!", nil
			} else {
				return "绑定失败!", nil
			}
		default:
			return "", errors.New("删除绑定数据失败")
		}
	} else {
		return "", err
	}
	return "", errors.New("删除绑定数据失败")
}

/*
  根据企业微信userId返回token
*/
func (l PingLogic) QyUserId(request *pb.QyUserIdRequest) (*pb.QyUserIdResponse, error) {
	userInfo, err := l.svcCtx.BandingWx.GetUserId(request.UserId)
	if err != nil {
		return nil, errorx.NewDefaultError("企业微信userId不存在")
	}
	//生成token
	now := time.Now().Unix()
	accessExpire := global.JwtAuth.AccessExpire
	jwtToken, err := jwt.GetJwtToken(
		now,
		global.JwtAuth.AccessExpire,
		global.JwtAuth.AccessSecret,
		userInfo.User_Name,
		userInfo.Company_ID,
		userInfo.Staff_ID,
		userInfo.User_ID,
	)
	if err != nil {
		return nil, err
	}
	return &pb.QyUserIdResponse{
		Token:  jwtToken,
		Expire: now + accessExpire,
	}, nil
}

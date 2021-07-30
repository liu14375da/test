package bandingWx

import (
	vars "ZeroProject/Rpc/qyWxBinDing/internal/model"
	"ZeroProject/Rpc/qyWxBinDing/internal/model/qyWx"
	"ZeroProject/model/entity/tbsfUser"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
	"strconv"
	"strings"
	"time"
)

var (
	userFieldNames        = builderx.RawFieldNames(&tbsfUser.TbsfUser{})
	userFieldNames2       = builderx.RawFieldNames(&qyWx.User{})
	userRows              = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet = strings.Join(stringx.Remove(userFieldNames2, "UserID", "CreateTime"), ",")
	cacheUserNumberPrefix = "cache#User#number#"
	cacheUserIdPrefix     = "cache#User#id#"
)

type (
	BinDing interface {
		SelectUser(ErpName string) (*tbsfUser.TbsfUser, error)
		SelectUserId(id string) (*qyWx.SelectUser, error)
		SelectStaffID(userid, id string) (int64, error)
		InsertUser(userid, name, mobile, email, alias, address, position, gender, staffId, telephone string, enable int8, invite bool, status int8, cron string) (int64, error)
		GetUserId(id string) (*tbsfUser.UserId, error)
	}
	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}
)

func NewBandingWx(s sqlx.SqlConn, c cache.CacheConf) BinDing {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(s, c),
		table:      "TBSF_User",
	}
}

func (m defaultUserModel) SelectUser(ErpName string) (*tbsfUser.TbsfUser, error) {
	var resp tbsfUser.TbsfUser
	err := m.QueryRow(&resp, ErpName, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where user_name = ? ", userRows, m.table)
		str := strings.Replace(query, "`", "", -1)
		return conn.QueryRow(v, str, ErpName)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, vars.ErrNotFound
	default:
		return nil, err
	}
}

func (m defaultUserModel) SelectUserId(id string) (*qyWx.SelectUser, error) {
	var resp qyWx.SelectUser
	err := m.QueryRow(&resp, id, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select User_Id from TBSF_Corporate_WxUser where User_Id = ? ")
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, vars.ErrNotFound
	default:
		return nil, err
	}
}

func (m defaultUserModel) SelectStaffID(userid, id string) (int64, error) {
	var res qyWx.CorporateWxUser
	_ = m.QueryRow(&res, id, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select User_Id ,Staff_ID from  TBSF_Corporate_WxUser  where Staff_ID=?")
		return conn.QueryRow(v, query, id)
	})

	if res.StaffId != "" {
		result, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
			query := fmt.Sprintf("update TBSF_Corporate_WxUser set User_Id=? where Staff_ID=?")
			return conn.Exec(query, userid, id)
		}, userid)
		if err != nil {
			return 0, nil
		}
		affectedRow, err := result.RowsAffected()
		if err != nil {
			return 0, nil
		}
		return affectedRow, nil
	}
	return 2, nil
}

var timeLayoutStr = "2006-01-02 15:04:05"

func (m defaultUserModel) InsertUser(userid, name, mobile, email, alias, address, position, gender, staffId, telephone string, enable int8, invite bool, status int8, cron string) (int64, error) {
	t := time.Now()
	ts := t.Format(timeLayoutStr)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("INSERT INTO TBSF_Corporate_WxUser(User_Id,Staff_ID,Name,Alias,Mobile,Position,Gender,Email,Telephone,Enable,Address,To_invite,Status,CreateTime,Cron)VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		fmt.Println(query)
		return conn.Exec(query, userid, staffId, name, alias, mobile, position, gender, email, telephone, strconv.Itoa(int(enable)), address, invite, status, ts, cron)
	}, userid)
	id, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m defaultUserModel) GetUserId(id string) (*tbsfUser.UserId, error) {
	var resp qyWx.CorporateWxUser
	err := m.QueryRow(&resp, id, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select User_Id,Staff_ID from TBSF_Corporate_WxUser where User_Id= ? ")
		return conn.QueryRow(v, query, id)
	})
	var user tbsfUser.UserId
	_ = m.QueryRow(&user, resp.StaffId, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select User_ID,User_Name,Staff_ID,User_Password ,Company_ID from TBSF_User where Staff_ID= ? ")
		return conn.QueryRow(v, query, resp.StaffId)
	})
	switch err {
	case nil:
		return &user, nil
	case sqlc.ErrNotFound:
		return nil, vars.ErrNotFound
	default:
		return nil, err
	}
}

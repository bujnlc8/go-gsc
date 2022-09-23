package controller

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/afocus/captcha"
	"github.com/bujnlc8/go-gsc/models"
	"github.com/bujnlc8/go-gsc/util"
	"github.com/gin-gonic/gin"
)

var cap = captcha.New()

type UserFeedBackData struct {
	FeedBackType int64  `json:"feedback_type"`
	Remark       string `json:"remark"`
	Captcha      string `json:"captcha"`
	Token        string `json:"token"`
}

func Code2Session(ctx *gin.Context) {
	code := ctx.Param("code")
	appType := ctx.DefaultQuery("app_type", "wechat")
	appId := ""
	appSecret := ""
	host := ""
	if appType == "wechat" {
		appId = util.GetConfStr("wxAppId")
		appSecret = util.GetConfStr("wxAppSecret")
		host = "api.weixin.qq.com"
	} else if appType == "qq" {
		appId = util.GetConfStr("qqAppId")
		appSecret = util.GetConfStr("qqAppSecret")
		host = "api.q.qq.com"
	}
	if len(appId) == 0 {
		ctx.JSON(400, models.ErrorResp{Code: -1, Msg: "参数错误"})
		return
	}
	url := fmt.Sprintf("https://%s/sns/jscode2session?appid=%s&grant_type=authorization_code&js_code=%s&secret=%s", host, appId, code, appSecret)
	res, err := Get(url)
	if err != nil {
		ctx.JSON(500, models.ReturnOpenId{Code: -1, Data: models.LoginResponse{OpenID: "", SessionKey: "", UnionID: ""}})
		return
	}
	ctx.JSON(200, models.ReturnOpenId{Code: 0, Data: res})
}

func Get(url string) (lres models.LoginResponse, err error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return models.LoginResponse{}, err
	}
	defer resp.Body.Close()
	var data models.LLoginResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return models.LoginResponse{}, err
	}

	if data.Errcode != 0 {
		err = errors.New(data.Errmsg)
		return models.LoginResponse{}, err
	}
	lres = data.LoginResponse
	return
}

func HandleUserFeedBack(ctx *gin.Context) {
	open_id := ctx.Param("open_id")
	gsc_id := ctx.Param("gsc_id")
	checkResult := true
	if len(open_id) == 0 || len(gsc_id) == 0 {
		checkResult = false
	}
	if v, err := strconv.Atoi(gsc_id); err != nil || v <= 0 {
		checkResult = false
	}
	if !checkResult {
		ctx.JSON(400, models.ErrorResp{Code: -1, Msg: "参数错误"})
		return
	}
	userFeedBackData := UserFeedBackData{}
	if err := ctx.BindJSON(&userFeedBackData); err != nil {
		ctx.JSON(400, models.ErrorResp{Code: -1, Msg: "数据格式错误"})
		return
	}
	res, err := util.DB.Exec(
		"UPDATE captcha set is_valid = 0 where open_id = ? and md5 = ? and str = ? and is_valid = 1",
		open_id,
		userFeedBackData.Token,
		userFeedBackData.Captcha,
	)
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "数据库异常"})
		return
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		ctx.JSON(400, models.ErrorResp{Code: -1, Msg: "验证码错误"})
		return
	}
	if err := models.DoUserFeedBack(open_id, gsc_id, userFeedBackData.FeedBackType, userFeedBackData.Remark); err != nil {
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "数据库异常"})
		return
	}
	ctx.JSON(200, models.ErrorResp{Code: 0, Msg: "反馈成功"})
}

func HandleCaptcha(ctx *gin.Context) {
	openId := ctx.Param("open_id")
	if len(openId) == 0 {
		ctx.JSON(400, models.ErrorResp{Code: -1, Msg: "参数错误"})
		return
	}
	if row, err := util.DB.Query(
		"SELECT COUNT(1) AS c FROM captcha WHERE open_id = ? AND create_time > ? ",
		openId,
		time.Now().Format("2006-01-02"),
	); err != nil {
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "数据库异常"})
		return
	} else {
		var totalNum int64
		for row.Next() {
			row.Scan(&totalNum)
		}
		if totalNum > 100 {
			ctx.JSON(403, models.ErrorResp{Code: -1, Msg: "验证码超限"})
			return
		}
	}
	cap.SetFont("comic.ttf")
	cap.SetFrontColor(color.RGBA{97, 113, 114, 255})
	cap.SetBkgColor(color.RGBA{243, 243, 242, 255})
	cap.SetSize(160, 80)
	img, str := cap.Create(6, captcha.NUM)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "系统错误"})
		return
	}
	b64Data := base64.StdEncoding.EncodeToString(buf.Bytes())
	md5Data := util.GetMd5(str)
	if _, err := util.DB.Exec(
		"INSERT INTO captcha(open_id, str, md5, create_time) VALUES(?, ?, ?, NOW())",
		openId,
		str,
		md5Data,
	); err != nil {
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "系统错误"})
		return
	}
	ctx.JSON(200, models.CaptchaResp{Code: 0, Token: md5Data, Captcha: b64Data})
}

func HandleAd(ctx *gin.Context) {
	openId := ctx.Param("open_id")
	if len(openId) == 0 {
		ctx.JSON(400, models.ErrorResp{Code: -1, Msg: "参数错误"})
		return
	}
	if rows, err := util.DB.Query(
		"SELECT COUNT(1) AS c FROM ad_whitelist WHERE OPEN_ID = ? AND IS_VALID = 1",
		openId,
	); err != nil {
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "数据库错误"})
		return
	} else {
		var count int64
		for rows.Next() {
			rows.Scan(&count)
		}
		if count == 0 {
			ctx.JSON(200, models.ReturnLike{Code: 0, Data: "invalid"})
		} else {
			ctx.JSON(200, models.ReturnLike{Code: 0, Data: "valid"})
		}
	}
}

func Code2SessionAliPay(ctx *gin.Context) {
	code := ctx.Param("code")
	appId := util.GetConfStr("alipayAppId")
	time := time.Now().Format("2006-01-02 15:04:05")
	signS := fmt.Sprintf(
		"app_id=%s&charset=UTF-8&code=%s&format=json&grant_type=authorization_code&method=alipay.system.oauth.token&sign_type=RSA2&timestamp=%s&version=1.0",
		appId,
		code,
		time,
	)
	alipayRsa := util.GetConfStr("alipayRsa")
	sign := util.Rsa2Sign(signS, alipayRsa)
	url := fmt.Sprintf(
		"https://openapi.alipay.com/gateway.do?app_id=%s&charset=UTF-8&code=%s&format=json&grant_type=authorization_code&method=alipay.system.oauth.token&sign_type=RSA2&timestamp=%s&version=1.0&sign=%s",
		appId,
		code,
		url.QueryEscape(time),
		url.QueryEscape(sign),
	)
	res, err := GetAlipay(url)
	if err != nil {
		ctx.JSON(500, models.ReturnOpenId{Code: -1, Data: models.LoginResponse{OpenID: "", SessionKey: "", UnionID: ""}})
		return
	}
	ctx.JSON(200, models.ReturnOpenId{Code: 0, Data: models.LoginResponse{OpenID: res.AlipaySystemOauthTokenResponse.UserId, SessionKey: "", UnionID: ""}})
}

func GetAlipay(url string) (models.AlipayResponse, error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	req.Header.Set("Host", "openapi.alipay.com")
	resp, err := client.Do(req)
	if err != nil {
		return models.AlipayResponse{}, err
	}
	defer resp.Body.Close()
	var data models.AlipayResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return models.AlipayResponse{}, err
	}
	return data, nil
}

func HandleCheckCaptcha(ctx *gin.Context) {
	open_id := ctx.DefaultQuery("open_id", "")
	token := ctx.DefaultQuery("token", "")
	captcha := ctx.DefaultQuery("captcha", "")
	res, err := util.DB.Exec(
		"UPDATE captcha set is_valid = 0 where open_id = ? and md5 = ? and str = ? and is_valid = 1",
		open_id,
		token,
		captcha,
	)
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "数据库异常"})
		return
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		ctx.JSON(200, models.ErrorResp{Code: -1, Msg: "验证码错误"})
		return
	}
	ctx.JSON(200, models.ErrorResp{Code: 0, Msg: "验证码正确"})
}

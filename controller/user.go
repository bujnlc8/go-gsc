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
	wxappId := util.GetConfStr("wxAppId")
	wxappSecret := util.GetConfStr("wxappSecret")
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&grant_type=authorization_code&js_code=%s&secret=%s", wxappId, code, wxappSecret)
	res, err := Get(url)
	if err != nil {
		fmt.Println(err)
	}
	ctx.JSON(200, models.ReturnOpenId{Code: 0, Data: res})
}

func Get(url string) (lres models.LoginResponse, err error) {
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var data models.LLoginResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}

	if data.Errcode != 0 {
		err = errors.New(data.Errmsg)
		return
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
		fmt.Println(err)
		ctx.JSON(500, models.ErrorResp{Code: -1, Msg: "系统错误"})
		return
	}
	ctx.JSON(200, models.CaptchaResp{Code: 0, Token: md5Data, Captcha: b64Data})
}

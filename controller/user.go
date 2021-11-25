package controller

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bujnlc8/go-gsc/models"
	"github.com/bujnlc8/go-gsc/util"
	"github.com/gin-gonic/gin"
)

type UserFeedBackData struct {
	FeedBackType int64  `json:"feedback_type"`
	Remark       string `json:"remark"`
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
		ctx.JSON(200, models.ReturnLike{Code: -1, Data: "反馈失败"})
		return
	}
	userFeedBackData := UserFeedBackData{}
	if err := ctx.BindJSON(&userFeedBackData); err != nil {
		ctx.JSON(200, models.ReturnLike{Code: -1, Data: "数据格式错误"})
		return
	}
	if err := models.DoUserFeedBack(open_id, gsc_id, userFeedBackData.FeedBackType, userFeedBackData.Remark); err != nil {
		ctx.JSON(200, models.ReturnLike{Code: -1, Data: "数据库异常"})
		return
	}
	ctx.JSON(200, models.ReturnLike{Code: 0, Data: "反馈成功"})
}

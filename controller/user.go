package controller

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bujnlc8/go-gsc/models"
	"github.com/bujnlc8/go-gsc/util"
	"github.com/gin-gonic/gin"
)

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

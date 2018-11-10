package controller

import (
	"github.com/gin-gonic/gin"
	"gogsc/util"
	"fmt"
	"github.com/medivhzhan/weapp"
	"gogsc/models"
)

func Code2Session(ctx *gin.Context)  {
	code := ctx.Param("code")
	wxappId := util.GetConfStr("wxAppId")
	wxappSecret:=util.GetConfStr("wxappSecret")
	res, err := weapp.Login(wxappId, wxappSecret, code)
	if err!=nil{
		fmt.Println(err)
	}
	ctx.JSON(200, models.ReturnOpenId{Code: 0, Data:res})
}


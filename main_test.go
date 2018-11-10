package main

import (
	"testing"
	"gogsc/util"
	"gogsc/models"
	"fmt"
	"gogsc/controller"
)

func Test_DB(t *testing.T)  {
	if util.DB == nil{
		t.Error("数据库连接失败")
	}else{
		t.Log("数据库连接正常")
	}
}

func Test_GetGSCById(t *testing.T)  {
	controller.GetGSCById(1, "123")
}

func Test_GetGSC30(t *testing.T)  {
	fmt.Println(models.GetGSC30())
}

func Test_GSCQuery(t *testing.T)  {
	fmt.Println(models.GSCQuery("宴山亭"))
}

func Test_GSCQueryLike(t *testing.T)  {
	fmt.Println(models.GSCQueryLike("宴山亭", "123"))
}
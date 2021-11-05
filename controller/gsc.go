package controller

import (
	"strconv"

	"github.com/bujnlc8/go-gsc/models"
	"github.com/gin-gonic/gin"
)

func GetGSCById(id int64, open_id string) models.GSC {
	gsc := models.GetGSCById(id, open_id)
	return gsc
}

func HandleIndex(ctx *gin.Context) {
	id_ := ctx.Param("id")
	if id_ == "all" {
		HandleIndexAll(ctx)
	} else {
		id, _ := strconv.ParseInt(id_, 10, 64)
		open_id := ctx.Param("open_id")
		gsc := GetGSCById(id, open_id)
		Returndata := models.ReturnDataSingle{Code: 0, Data: models.ReturnDataInerSingle{Msg: "success", Data: gsc}}
		ctx.JSON(200, Returndata)
	}
}

func HandleIndexAll(ctx *gin.Context) {
	gscs := models.GetGSC30()
	if len(gscs) == 0 {
		gscs = make([]models.GSC, 0)
	}
	ReturnData := models.ReturnDataList{Code: 0, Data: models.ReturnDataIner{Msg: "success", Data: gscs}}
	ctx.JSON(200, ReturnData)
}

func HandleQuery(ctx *gin.Context) {
	q := ctx.Param("q")
	page := ctx.Param("page")
	open_id := ctx.Param("open_id")
	var gscs []models.GSC
	if page == "main" {
		gscs = models.GSCQuery(q)
	} else {
		gscs = models.GSCQueryLike(q, open_id)
	}
	if len(gscs) == 0 {
		gscs = make([]models.GSC, 0)
	}
	ReturnData := models.ReturnDataList{Code: 0, Data: models.ReturnDataIner{Msg: "success", Data: gscs}}
	ctx.JSON(200, ReturnData)
}

func QueryMyLike(ctx *gin.Context) {
	open_id := ctx.Param("open_id")
	var gscs []models.GSC
	gscs = models.GSCQueryLike("", open_id)
	if len(gscs) == 0 {
		gscs = make([]models.GSC, 0)
	}
	ReturnData := models.ReturnDataList{Code: 0, Data: models.ReturnDataIner{Msg: "success", Data: gscs}}
	ctx.JSON(200, ReturnData)
}

func SetUserLike(ctx *gin.Context) {
	open_id := ctx.Param("open_id")
	gsc_id := ctx.Param("gsc_id")
	check_result := true
	if len(open_id) == 0 || len(gsc_id) == 0 {
		check_result = false
	}
	if v, err := strconv.Atoi(gsc_id); err != nil || v <= 0 {
		check_result = false
	}
	if !check_result {
		ctx.JSON(200, models.ReturnLike{Code: -1, Data: "收藏失败"})
	} else {
		if models.SetLike(open_id, gsc_id, 1) {
			ctx.JSON(200, models.ReturnLike{Code: 0, Data: "收藏成功"})
		} else {
			ctx.JSON(200, models.ReturnLike{Code: -1, Data: "收藏失败"})
		}

	}
}

func SetUserDisLike(ctx *gin.Context) {
	open_id := ctx.Param("open_id")
	gsc_id := ctx.Param("gsc_id")
	check_result := true
	if len(open_id) == 0 || len(gsc_id) == 0 {
		check_result = false
	}
	if v, err := strconv.Atoi(gsc_id); err != nil || v <= 0 {
		check_result = false
	}
	if !check_result {
		ctx.JSON(200, models.ReturnLike{Code: -1, Data: "收藏失败"})
	} else {
		if models.SetLike(open_id, gsc_id, 0) {
			ctx.JSON(200, models.ReturnLike{Code: 0, Data: "取消成功"})
		} else {
			ctx.JSON(200, models.ReturnLike{Code: -1, Data: "取消失败"})

		}
	}
}

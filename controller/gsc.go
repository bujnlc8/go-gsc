package controller

import (
	"strconv"

	"github.com/bujnlc8/go-gsc/models"
	"github.com/bujnlc8/go-gsc/util"
	"github.com/gin-gonic/gin"
)

func HandleIndex(ctx *gin.Context) {
	id_ := ctx.Param("id")
	if id_ == "all" {
		HandleIndexAll(ctx)
	} else {
		id, _ := strconv.ParseInt(id_, 10, 64)
		open_id := ctx.Param("open_id")
		gsc, err := models.GetGSCById(id, open_id)
		if err != nil {
			ctx.JSON(500, models.ErrorResp{Code: 500, Msg: "系统错误"})
			return
		}
		Returndata := models.ReturnDataSingle{Code: 0, Data: models.ReturnDataInerSingle{Msg: "success", Data: gsc}}
		ctx.JSON(200, Returndata)
	}
}

func HandleIndexAll(ctx *gin.Context) {
	gscs, err := models.GetGSC30()
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: 500, Msg: "系统错误"})
		return
	}
	if len(gscs) == 0 {
		gscs = make([]models.GSC, 0)
	}
	ReturnData := models.ReturnDataList{Code: 0, Data: models.ReturnDataIner{Msg: "success", Data: gscs}}
	ctx.JSON(200, ReturnData)
}

func HandleShortIndex(ctx *gin.Context) {
	gscs, err := models.GetGSCSimple20()
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: 500, Msg: "系统错误"})
		return
	}
	if len(gscs) == 0 {
		gscs = make([]models.GSCSimple, 0)
	}
	ReturnData := models.ReturnSimpleDataList{Code: 0, Data: models.ReturnSimpleDataIner{Msg: "success", Data: gscs, Total: 20, SplitWords: ""}}
	ctx.JSON(200, ReturnData)
}

func HandleQuery(ctx *gin.Context) {
	q := ctx.Param("q")
	page := ctx.Param("page")
	open_id := ctx.Param("open_id")
	var gscs []models.GSC
	var err error
	if page == "main" {
		gscs, err = models.GSCQuery(q)
	} else {
		gscs, err = models.GSCQueryLike(q, open_id)
	}
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: 500, Msg: "系统错误"})
		return
	}
	if len(gscs) == 0 {
		gscs = make([]models.GSC, 0)
	}
	ReturnData := models.ReturnDataList{Code: 0, Data: models.ReturnDataIner{Msg: "success", Data: gscs}}
	ctx.JSON(200, ReturnData)
}

func HandleQueryByPage(ctx *gin.Context) {
	q := ctx.Param("q")
	page := ctx.Param("page")
	open_id := ctx.Param("open_id")
	page_size := ctx.DefaultQuery("page_size", "50")
	page_size_int, err := strconv.ParseInt(page_size, 10, 0)
	if err != nil {
		ctx.JSON(400, models.ErrorResp{Code: 400, Msg: "参数错误"})
		return
	}
	if page_size_int <= 0 {
		page_size_int = 50
	}
	// 页码
	page_num := ctx.DefaultQuery("page_num", "1")
	page_num_int, err := strconv.ParseInt(page_num, 10, 0)
	if err != nil {
		ctx.JSON(400, models.ErrorResp{Code: 400, Msg: "参数错误"})
		return
	}
	if page_num_int <= 0 {
		page_num_int = 1
	}
	search_pattern := ctx.DefaultQuery("search_pattern", "all")
	gscs := make([]models.GSCSimple, 0)
	var total int64
	var splitWords string
	if page == "main" {
		gscs, total, splitWords, err = models.GSCQueryByPage(q, page_size_int, page_num_int, search_pattern)
	} else {
		gscs, total, splitWords, err = models.GSCQueryLikeByPage(q, open_id, page_size_int, page_num_int, search_pattern)
	}
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: 500, Msg: "系统错误"})
		return
	}
	ReturnData := models.ReturnSimpleDataList{Code: 0, Data: models.ReturnSimpleDataIner{Msg: "success", Data: gscs, Total: total, SplitWords: splitWords}}
	ctx.JSON(200, ReturnData)
}

func QueryMyLike(ctx *gin.Context) {
	open_id := ctx.Param("open_id")
	var gscs []models.GSC
	var err error
	gscs, err = models.GSCQueryLike("", open_id)
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: 500, Msg: "系统错误"})
		return
	}
	if len(gscs) == 0 {
		gscs = make([]models.GSC, 0)
	}
	ReturnData := models.ReturnDataList{Code: 0, Data: models.ReturnDataIner{Msg: "success", Data: gscs}}
	ctx.JSON(200, ReturnData)
}

func QueryMyLikeByPage(ctx *gin.Context) {
	open_id := ctx.Param("open_id")
	page_size := ctx.DefaultQuery("page_size", "50")
	page_size_int, err := strconv.ParseInt(page_size, 10, 0)
	if err != nil {
		ctx.JSON(400, models.ErrorResp{Code: 400, Msg: "参数错误"})
		return
	}
	if page_size_int <= 0 {
		page_size_int = 50
	}
	// 页码
	page_num := ctx.DefaultQuery("page_num", "1")
	page_num_int, err := strconv.ParseInt(page_num, 10, 0)
	if err != nil {
		ctx.JSON(400, models.ErrorResp{Code: 400, Msg: "参数错误"})
		return
	}
	if page_num_int <= 0 {
		page_num_int = 1
	}
	search_pattern := ctx.DefaultQuery("search_pattern", "all")
	var gscs []models.GSCSimple
	var total int64
	var splitWords string
	gscs, total, splitWords, err = models.GSCQueryLikeByPage("", open_id, page_size_int, page_num_int, search_pattern)
	if err != nil {
		ctx.JSON(500, models.ErrorResp{Code: 500, Msg: "系统错误"})
		return
	}
	ReturnData := models.ReturnSimpleDataList{Code: 0, Data: models.ReturnSimpleDataIner{Msg: "success", Data: gscs, Total: total, SplitWords: splitWords}}
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

func GetSignUrlForAudio(ctx *gin.Context) {
	audioRequest := models.AudioRequest{}
	if err := ctx.BindJSON(&audioRequest); err != nil {
		ctx.JSON(400, models.ErrorResp{Code: -1, Msg: "参数错误"})
		return
	}
	ctx.JSON(200, models.ReturnLike{Code: 0, Data: util.GetMd5ForAudioUrl(audioRequest.FileName)})
}

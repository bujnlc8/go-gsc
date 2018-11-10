package models

import (
	"gogsc/util"
	"fmt"
	"crypto/rand"
	"math/big"
	"database/sql"
	"github.com/medivhzhan/weapp"
	"strings"
)

type ReturnDataList struct {
	Code int8 `json:"code"`
	Data ReturnDataIner `json:"data"`
}

type ReturnDataIner struct {
	Msg string `json:"msg"`
	Data []GSC `json:"data"`
}

type ReturnDataSingle struct {
	Code int8 `json:"code"`
	Data ReturnDataInerSingle `json:"data"`
}

type ReturnDataInerSingle struct {
	Msg string `json:"msg"`
	Data GSC `json:"data"`
}

type GSC struct {
	Id             int64  `json:"id"`
	Work_title     string `json:"work_title"`
	Work_author    string `json:"work_author"`
	Work_dynasty   string `json:"work_dynasty"`
	Content        string `json:"content"`
	Translation    string `json:"translation"`
	Intro          string `json:"intro"`
	Annotation_    string `json:"annotation"`
	Foreword       string `json:"foreword"`
	Appreciation   string `json:"appreciation"`
	Master_comment string `json:"master_comment"`
	Layout         string `json:"layout"`
	Audio_id       int64  `json:"audio_id"`
	Like           int8   `json:"like"`
}

type ReturnOpenId struct {
	Code int8 `json:"code"`
	Data weapp.LoginResponse `json:"data"`
}

type ReturnLike struct {
	Code int8 `json:"code"`
	Data string `json:"data"`
} 

func processRows(rows *sql.Rows) []GSC  {
	var GSCS []GSC
	for rows.Next() {
		var gsc = new(GSC)
		rows.Scan(&gsc.Id, &gsc.Work_title, &gsc.Work_author,
			&gsc.Work_dynasty, &gsc.Content, &gsc.Translation,
			&gsc.Intro, &gsc.Annotation_, &gsc.Foreword,
			&gsc.Appreciation, &gsc.Master_comment, &gsc.Layout,
			&gsc.Audio_id)
		GSCS = append(GSCS, *gsc)
	}
	return GSCS
}

func GetGSCById(id int64, open_id string) GSC {
	rows, err := util.DB.Query(
		"select id, work_title, work_author, work_dynasty, content, "+
			"translation, intro, annotation_, foreword, appreciation, "+
			"master_comment, layout, audio_id from gsc  where id = ? ", id)
	if err != nil {
		fmt.Println(err)
	}
	var gsc = new(GSC)
	for rows.Next() {
		rows.Scan(&gsc.Id, &gsc.Work_title, &gsc.Work_author,
			&gsc.Work_dynasty, &gsc.Content, &gsc.Translation,
			&gsc.Intro, &gsc.Annotation_, &gsc.Foreword,
			&gsc.Appreciation, &gsc.Master_comment, &gsc.Layout, &gsc.Audio_id)
	}
	//查询是否喜欢
	if gsc.Id != 0 {
		rows, err := util.DB.Query(
			"select id from user_like_gsc where open_id=? and gsc_id=? ", open_id, id)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			gsc.Like = 1
		}
	}
	return *gsc
}

// 获取随机30条数据
func GetGSC30() []GSC{
	var RandId int64
	for i:=0;i<30;i++{
		result, _ := rand.Int(rand.Reader, big.NewInt(8069))
		RandId = result.Int64()
	}
	rows, err := util.DB.Query(
		"select id, work_title, work_author, work_dynasty, content, "+
			"translation, intro, annotation_, foreword, appreciation, "+
			"master_comment, layout, audio_id from gsc where id > ?  order by audio_id desc limit 30", RandId)
	if err != nil {
		fmt.Println(err)
	}
	return processRows(rows)
}

func GSCQuery(q string) []GSC {
	var rows *sql.Rows
	var err error
	if q!="音频" {
		rows, err = util.DB.Query("select id, work_title, work_author, work_dynasty, "+
			"content, `translation`, intro, annotation_, foreword, appreciation, "+
			"master_comment, layout, audio_id from gsc "+
			"where work_title like ? or work_author like ? order by audio_id desc", "%"+q+"%", "%"+q+"%")
		if err != nil {
			fmt.Println(err)
		}
	}else {
		var RandId int64
		for i:=0;i<30;i++{
			result, _ := rand.Int(rand.Reader, big.NewInt(3000))
			RandId = result.Int64()
		}
		rows, err = util.DB.Query("select id, work_title, work_author, work_dynasty, "+
			"content, `translation`, intro, annotation_, foreword, appreciation, "+
			"master_comment, layout, audio_id from gsc "+
			"where audio_id > 0 and id > ? limit 100", RandId)
		if err != nil {
			fmt.Println(err)
		}
	}
	return processRows(rows)
}

func GSCQueryLike(q string, open_id string) []GSC {
	rows, err := util.DB.Query(
		"select gsc_id from user_like_gsc where open_id=? ", open_id)
	if err != nil {
		fmt.Println(err)
	}
	var gscids []string
	for rows.Next() {
		var gsc_id string
		rows.Scan(&gsc_id)
		gscids = append(gscids, gsc_id)
	}
	if len(gscids) == 0{
		gscids = append(gscids, "-1")
	}
	gscids_str:=strings.Join(gscids, ",")
	if q != "" {
		rows, err = util.DB.Query(
			"select id, work_title, work_author, work_dynasty, content, "+
				"translation, intro, annotation_, foreword, appreciation, master_comment, layout,"+
				"audio_id from gsc where id in ("+ gscids_str + ") and (work_title like ? "+
				" or work_author like ? ) order by audio_id desc",
			"%"+q+"%", "%"+q+"%")
		if err != nil {
			fmt.Println(err)
		}
	}else {
		rows, err = util.DB.Query(
			"select id, work_title, work_author, work_dynasty, content, "+
				"translation, intro, annotation_, foreword, appreciation, master_comment, layout,"+
				"audio_id from gsc where id in ("+ gscids_str + ") order by audio_id desc")
		if err != nil {
			fmt.Println(err)
		}
	}
	return processRows(rows)
}

func SetLike(open_id string, gsc_id string, operate int8) bool{
	if operate == 1{
		result, err := util.DB.Exec("insert into user_like_gsc(open_id, gsc_id)values(?, ?)", open_id, gsc_id)
		if err!=nil{
			fmt.Println(err)
			return false
		}
		rows_affected, _ := result.RowsAffected();if rows_affected==1{
			return true
		}
	}else{
		result, err := util.DB.Exec("delete from user_like_gsc where open_id=? and gsc_id=?", open_id, gsc_id)
		if err!=nil{
			fmt.Println(err)
			return false
		}
		rows_affected, _ := result.RowsAffected();if rows_affected==1{
			return true
		}
	}
	return false
}



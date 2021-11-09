package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/bujnlc8/go-gsc/util"
)

type ErrorResp struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type MayNullStr sql.NullString

func (s *MayNullStr) Scan(value interface{}) error {
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	if reflect.TypeOf(value) == nil {
		*s = MayNullStr{i.String, false}
	} else {
		*s = MayNullStr{i.String, true}
	}
	return nil
}

func (s MayNullStr) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(""), nil
}

type ReturnDataList struct {
	Code int8           `json:"code"`
	Data ReturnDataIner `json:"data"`
}

type ReturnSimpleDataList struct {
	Code int8                 `json:"code"`
	Data ReturnSimpleDataIner `json:"data"`
}

type ReturnDataIner struct {
	Msg  string `json:"msg"`
	Data []GSC  `json:"data"`
}

type ReturnSimpleDataIner struct {
	Msg   string      `json:"msg"`
	Data  []GSCSimple `json:"data"`
	Total int64       `json:"total"`
}

type ReturnDataSingle struct {
	Code int8                 `json:"code"`
	Data ReturnDataInerSingle `json:"data"`
}

type ReturnDataInerSingle struct {
	Msg  string `json:"msg"`
	Data GSC    `json:"data"`
}

type GSC struct {
	Id             int64      `json:"id"`
	Work_title     string     `json:"work_title"`
	Work_author    string     `json:"work_author"`
	Work_dynasty   string     `json:"work_dynasty"`
	Content        MayNullStr `json:"content"`
	Translation    MayNullStr `json:"translation"`
	Intro          MayNullStr `json:"intro"`
	Annotation_    MayNullStr `json:"annotation"`
	Foreword       MayNullStr `json:"foreword"`
	Appreciation   MayNullStr `json:"appreciation"`
	Master_comment MayNullStr `json:"master_comment"`
	Layout         MayNullStr `json:"layout"`
	Audio_id       int64      `json:"audio_id"`
	Like           int8       `json:"like"`
	Score          float64    `json:"score"`
}

type GSCSimple struct {
	Id           int64      `json:"id"`
	Work_title   string     `json:"work_title"`
	Work_author  string     `json:"work_author"`
	Work_dynasty string     `json:"work_dynasty"`
	Content      MayNullStr `json:"content"`
	Audio_id     int64      `json:"audio_id"`
	Score        float64    `json:"score"`
}

type ReturnOpenId struct {
	Code int8          `json:"code"`
	Data LoginResponse `json:"data"`
}

type ReturnLike struct {
	Code int8   `json:"code"`
	Data string `json:"data"`
}

type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
type LoginResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

type LLoginResponse struct {
	Response
	LoginResponse
}

func processRows(rows *sql.Rows) []GSC {
	var GSCS []GSC
	for rows.Next() {
		var gsc = new(GSC)
		rows.Scan(&gsc.Id, &gsc.Work_title, &gsc.Work_author,
			&gsc.Work_dynasty, &gsc.Content, &gsc.Translation,
			&gsc.Intro, &gsc.Annotation_, &gsc.Foreword,
			&gsc.Appreciation, &gsc.Master_comment, &gsc.Layout,
			&gsc.Audio_id, &gsc.Score)
		GSCS = append(GSCS, *gsc)
	}
	return GSCS
}

func processSimpleRows(rows *sql.Rows) []GSCSimple {
	var GSCS []GSCSimple
	for rows.Next() {
		var gsc = new(GSCSimple)
		rows.Scan(&gsc.Id, &gsc.Work_title, &gsc.Work_author,
			&gsc.Work_dynasty, &gsc.Content,
			&gsc.Audio_id, &gsc.Score)
		GSCS = append(GSCS, *gsc)
	}
	return GSCS
}

func GetGSCById(id int64, open_id string) GSC {
	rows, err := util.DB.Query(
		"SELECT `id`, work_title, work_author, work_dynasty, content, "+
			"translation, intro, annotation_, foreword, appreciation, "+
			"master_comment, layout, audio_id, 0 FROM gsc  WHERE `id` = ? ", id)
	if err != nil {
		fmt.Println(err)
	}
	var gsc = new(GSC)
	for rows.Next() {
		rows.Scan(&gsc.Id, &gsc.Work_title, &gsc.Work_author,
			&gsc.Work_dynasty, &gsc.Content, &gsc.Translation,
			&gsc.Intro, &gsc.Annotation_, &gsc.Foreword,
			&gsc.Appreciation, &gsc.Master_comment, &gsc.Layout, &gsc.Audio_id, &gsc.Score)
	}
	//查询是否喜欢
	if gsc.Id != 0 {
		rows, err := util.DB.Query(
			"SELECT `id` FROM user_like_gsc WHERE open_id=? AND  gsc_id=? ", open_id, id)
		if err != nil {
			fmt.Println(err)
		}
		for rows.Next() {
			gsc.Like = 1
		}
	}
	return *gsc
}

// GetGSC30 获取随机30条数据
func GetGSC30() []GSC {
	rows, err := util.DB.Query(
		"SELECT `id`, work_title, work_author, work_dynasty, content, " +
			"translation, intro, annotation_, foreword, appreciation, " +
			"master_comment, layout, audio_id, 0 FROM gsc WHERE audio_id > 0 ORDER BY RAND() LIMIT 30")
	if err != nil {
		fmt.Println(err)
	}
	return processRows(rows)
}

func GetGSCSimple30() []GSCSimple {
	rows, err := util.DB.Query(
		"SELECT `id`, work_title, work_author, work_dynasty, SUBSTRING(content, 1, 50), " +
			"audio_id, 0 FROM gsc WHERE audio_id > 0 ORDER BY RAND() LIMIT 30")
	if err != nil {
		fmt.Println(err)
	}
	return processSimpleRows(rows)
}

func GSCQuery(q string) []GSC {
	var rows *sql.Rows
	var err error
	if q != "音频" {
		againstS := util.AgainstString(q)
		rows, err = util.DB.Query(
			"SELECT `id`, work_title, work_author, work_dynasty, " +
				"content, translation, intro, annotation_, foreword, appreciation, " +
				"master_comment, layout, audio_id , MATCH(work_author, work_title, work_dynasty, content)" +
				" AGAINST ('" + againstS + "' IN BOOLEAN MODE) AS score FROM gsc " +
				" WHERE MATCH(work_author, work_title, work_dynasty, content) " +
				"AGAINST ('" + againstS + "' IN  BOOLEAN MODE) ORDER BY audio_id DESC,score DESC LIMIT 500")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		rows, err = util.DB.Query("SELECT `id`, work_title, work_author, work_dynasty, " +
			"content, `translation`, intro, annotation_, foreword, appreciation, " +
			"master_comment, layout, audio_id, 0 FROM gsc " +
			"WHERE audio_id > 0 ORDER BY RAND() LIMIT 100")
		if err != nil {
			fmt.Println(err)
		}
	}
	return processRows(rows)
}

func GSCQueryByPage(q string, page_size int64, page_num int64, search_pattern string) ([]GSCSimple, int64, error) {
	var rows *sql.Rows
	var err error
	offset := (page_num - 1) * page_size
	var total int64
	if q != "音频" {
		againstS := util.AgainstString(q)
		matchS := util.MatchStringBySearchPattern(search_pattern)
		sql := fmt.Sprintf("SELECT `id`, work_title, work_author, work_dynasty, "+
			"SUBSTRING(content, 1, 50) AS c, audio_id , %s  AGAINST ('%s' IN BOOLEAN MODE) AS score FROM gsc WHERE %s "+
			"AGAINST ('%s' IN  BOOLEAN MODE) ORDER BY audio_id DESC,score DESC LIMIT %d OFFSET %d", matchS, againstS, matchS, againstS, page_size, offset)
		rows, err = util.DB.Query(sql)
		if err != nil {
			return nil, 0, err
		}
		sql = fmt.Sprintf("SELECT count(1) AS c FROM gsc WHERE %s AGAINST ('%s' IN  BOOLEAN MODE)", matchS, againstS)
		total_rows, err := util.DB.Query(sql)
		if err != nil {
			return nil, 0, err
		}
		for total_rows.Next() {
			total_rows.Scan(&total)
		}
	} else {
		rows, err = util.DB.Query("SELECT `id`, work_title, work_author, work_dynasty, " +
			"SUBSTRING(content, 1, 50) AS c, audio_id, 0 FROM gsc " +
			"WHERE audio_id > 0 ORDER BY RAND() LIMIT 100")
		if err != nil {
			return nil, 0, err
		}
		total = int64(100)
	}
	return processSimpleRows(rows), total, nil
}

func GSCQueryLike(q string, open_id string) []GSC {
	rows, err := util.DB.Query(
		"SELECT gsc_id FROM user_like_gsc WHERE open_id=? ", open_id)
	if err != nil {
		fmt.Println(err)
	}
	var gscids []string
	for rows.Next() {
		var gsc_id string
		rows.Scan(&gsc_id)
		gscids = append(gscids, gsc_id)
	}
	if len(gscids) == 0 {
		gscids = append(gscids, "-1")
	}
	gscids_str := strings.Join(gscids, ",")
	if q != "" {
		againstS := util.AgainstString(q)
		rows, err = util.DB.Query(
			"SELECT `id`, work_title, work_author, work_dynasty, content, " +
				"translation, intro, annotation_, foreword, appreciation, master_comment, layout," +
				"audio_id,  MATCH(work_author, work_title, work_dynasty, content) AGAINST ('" + againstS + "' IN BOOLEAN MODE) AS score " +
				"FROM gsc WHERE MATCH(work_author, work_title, work_dynasty, content) " +
				"AGAINST ('" + againstS + "' IN BOOLEAN MODE) AND  `id` IN (" + gscids_str + ") ORDER BY audio_id DESC, score DESC")
		if err != nil {
			fmt.Println(err)
		}
	} else {
		rows, err = util.DB.Query(
			"SELECT `id`, work_title, work_author, work_dynasty, content, " +
				"translation, intro, annotation_, foreword, appreciation, master_comment, layout," +
				"audio_id, 0 FROM gsc WHERE `id` IN (" + gscids_str + ") ORDER BY audio_id DESC")
		if err != nil {
			fmt.Println(err)
		}
	}
	return processRows(rows)
}

func GSCQueryLikeByPage(q string, open_id string, page_size int64, page_num int64, search_pattern string) ([]GSCSimple, int64, error) {
	rows, err := util.DB.Query(
		"SELECT gsc_id FROM user_like_gsc WHERE open_id=? ", open_id)
	if err != nil {
		return nil, 0, err
	}
	var gscids []string
	for rows.Next() {
		var gsc_id string
		rows.Scan(&gsc_id)
		gscids = append(gscids, gsc_id)
	}
	total := int64(len(gscids))
	if len(gscids) == 0 {
		gscids = append(gscids, "-1")
	}
	offset := (page_num - 1) * page_size
	gscids_str := strings.Join(gscids, ",")
	if q != "" {
		againstS := util.AgainstString(q)
		matchS := util.MatchStringBySearchPattern(search_pattern)
		sql := fmt.Sprintf("SELECT `id`, work_title, work_author, work_dynasty, SUBSTRING(content, 1, 50) AS c, audio_id, "+
			" %s AGAINST ('%s' IN BOOLEAN MODE) AS score FROM gsc WHERE %s "+
			"AGAINST ('%s' IN BOOLEAN MODE) AND  `id` IN (%s) ORDER BY audio_id DESC, score DESC LIMIT %d OFFSET %d", matchS, againstS, matchS, againstS, gscids_str, page_size, offset)
		rows, err = util.DB.Query(sql)
		if err != nil {
			return nil, 0, err
		}
		sql = fmt.Sprintf("SELECT count(1) as c FROM gsc WHERE %s AGAINST ('%s' IN BOOLEAN MODE) AND  `id` IN (%s)", matchS, againstS, gscids_str)
		total_rows, err := util.DB.Query(sql)
		if err != nil {
			return nil, 0, err
		}
		for total_rows.Next() {
			total_rows.Scan(&total)
		}
	} else {
		sql := fmt.Sprintf("SELECT `id`, work_title, work_author, work_dynasty, SUBSTRING(content, 1, 50) AS c, audio_id, 0 FROM gsc WHERE `id` IN (%s) ORDER BY audio_id DESC LIMIT %d OFFSET %d", gscids_str, page_size, offset)
		rows, err = util.DB.Query(sql)
		if err != nil {
			return nil, 0, err
		}
	}
	return processSimpleRows(rows), total, nil
}

func SetLike(open_id string, gsc_id string, operate int8) bool {
	if operate == 1 {
		result, err := util.DB.Exec("INSERT  INTO  user_like_gsc(open_id, gsc_id) VALUES (?, ?)", open_id, gsc_id)
		if err != nil {
			fmt.Println(err)
			return false
		}
		rows_affected, _ := result.RowsAffected()
		if rows_affected == 1 {
			return true
		}
	} else {
		result, err := util.DB.Exec("DELETE FROM user_like_gsc WHERE open_id=? AND gsc_id=?", open_id, gsc_id)
		if err != nil {
			fmt.Println(err)
			return false
		}
		rows_affected, _ := result.RowsAffected()
		if rows_affected == 1 {
			return true
		}
	}
	return false
}

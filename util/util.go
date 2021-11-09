package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/yanyiwu/gojieba"
)

type CONF map[string]interface{}

var Confs = make(CONF)

var JieBa = gojieba.NewJieba()

// GetPrevDir ...
func GetPrevDir(path string) string {
	latsindex := strings.LastIndex(path, "/")
	return path[:latsindex]
}

// InitConf ...
func InitConf() {
	if len(Confs) != 0 {
		fmt.Println(Confs)
		return
	}
	CurrentDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fl, err := os.Open(CurrentDir + "/conf/config")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer fl.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if n == 0 {
			break
		}
	}
	for _, v := range strings.Split(string(buf), "\n") {
		v = strings.Replace(v, " ", "", -1)
		index := strings.Index(v, ":")
		if index == -1 {
			continue
		}
		k := strings.Replace(v[:index], " ", "", -1)
		v := strings.Replace(v[index+1:], " ", "", -1)
		Confs[k] = v
	}
}

// GetConf ...
func GetConf(name string) interface{} {
	return Confs[name]
}

func GetConfStr(name string) string {
	if os.Getenv("GSC_DEBUG") == "true" {
		return fmt.Sprintf("%v", Confs[name])
	}
	return os.Getenv(name)
}

var RE = regexp.MustCompile("[。？！，；、, ? . ! ; \\s ·]")

func SplitString(s string) []string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) <= 2 {
		return []string{s}
	}
	res := JieBa.Cut(s, true)
	newRes := make([]string, 0)
	for _, ss := range res {
		ss = RE.ReplaceAllString(ss, "")
		if utf8.RuneCountInString(ss) <= 0 {
			continue
		}
		newRes = append(newRes, ss)
	}
	return newRes
}

func AgainstString(s string) string {
	splitRes := SplitString(s)
	// 如果分词数量小于4，全部匹配
	splitLen := len(splitRes)
	if splitLen <= 3 {
		return "+" + strings.Join(splitRes, " +")
	}
	res := ""
	l := int(float64(splitLen) * 0.8)
	for i, qq := range splitRes {
		if i < l {
			res += " +" + qq
		} else {
			res += " " + qq
		}
	}
	return res
}

func MatchStringBySearchPattern(search_pattern string) string {
	if search_pattern == "title" {
		return "MATCH(work_title)"
	} else if search_pattern == "content" {
		return "MATCH(content, foreword)"
	} else if search_pattern == "author" {
		return "MATCH(work_author)"
	}
	return "MATCH(work_author, work_title, work_dynasty, content)"
}

func init() {
	if os.Getenv("GSC_DEBUG") == "true" {
		InitConf()
	}
}

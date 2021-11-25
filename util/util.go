package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/yanyiwu/gojieba"
)

var JieBa = gojieba.NewJieba()

func GetConfStr(name string) string {
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
	} else if search_pattern == "dynasty" {
		return "MATCH(work_dynasty)"
	}
	return "MATCH(work_author, work_title, work_dynasty, content)"
}

func GetMd5(s string) string {
	data := []byte(s + os.Getenv("md5Secret"))
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

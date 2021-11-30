package util

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
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

func GetMd5ForAudioUrl(fileName string) string {
	time := time.Now().Format("200601021504")
	data := []byte(os.Getenv("audioSecret") + time + fileName)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	h := hex.EncodeToString(md5Ctx.Sum(nil))
	return fmt.Sprintf("%s/%s/%s%s", os.Getenv("audioDomain"), time, h, fileName)
}

const (
	// 私钥 PEMBEGIN 开头
	PEMBEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	// 私钥 PEMEND 结尾
	PEMEND = "\n-----END RSA PRIVATE KEY-----"
	// 公钥 PEMBEGIN 开头
	PUBPEMBEGIN = "-----BEGIN PUBLIC KEY-----\n"
	// 公钥 PEMEND 结尾
	PUBPEMEND = "\n-----END PUBLIC KEY-----"
)

func FormatPrivateKey(privateKey string) string {
	if !strings.HasPrefix(privateKey, PEMBEGIN) {
		privateKey = PEMBEGIN + privateKey
	}
	if !strings.HasSuffix(privateKey, PEMEND) {
		privateKey = privateKey + PEMEND
	}
	return privateKey
}

func Rsa2Sign(s string, privateKey string) string {
	privateKey = FormatPrivateKey(privateKey)
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return ""
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}
	h := sha256.New()
	h.Write([]byte(s))
	digest := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, digest)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(signature)
}

package util

import (
	"fmt"
	"os"
	"strings"
)

type CONF map[string]interface{}

var Confs = make(CONF)

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

func init() {
	if os.Getenv("GSC_DEBUG") == "true" {
		InitConf()
	}
}

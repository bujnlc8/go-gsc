package util

import (
	"testing"
)

func Test_SplitString(t *testing.T) {
	s := "苏轼 定州   远人 "
	res := SplitString(s)
	if len(res) != 3 {
		t.Error(" SplitString error, expect length is 3, but get ", len(res), res)
	}
	s = " 苏轼定州怀古 "
	res = SplitString(s)
	if len(res) != 3 {
		t.Error(" SplitString error, expect length is 3, but get ", len(res), res)
	}
	res = SplitString("纳兰性德")
	t.Log(res)
	res = SplitString("谏逐客书")
	t.Log(res)
	res = SplitString("张抡")
	t.Log(res)
	res = SplitString("不负如来不负卿")
	t.Log(res)
}

package util

import (
	"testing"
)

func Test_SplitString(t *testing.T) {
	s := "苏轼 定州   远人 "
	res := SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
	if len(res) != 3 {
		t.Error(" SplitString error, expect length is 3, but get ", len(res), res)
	}
	s = " 苏轼定州怀古 "
	res = SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
	if len(res) != 3 {
		t.Error(" SplitString error, expect length is 3, but get ", len(res), res)
	}
	s = " 纳兰性德 "
	res = SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
	s = " 谏逐客书 "
	res = SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
	s = "张抡"
	res = SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
	s = "不负如来不负卿"
	res = SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
	s = "裁剪冰绡，轻叠数重，淡着燕脂匀注。新样靓妆，艳溢香融"
	res = SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
	s = "浣溪沙 · 游蕲水清泉寺，寺临兰溪，溪水西流"
	res = SplitString(s)
	t.Log(res)
	t.Log(AgainstString(s))
}

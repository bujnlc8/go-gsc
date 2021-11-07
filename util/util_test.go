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

}

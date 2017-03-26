package pac

import (
	"testing"
	"github.com/gwuhaolin/gfwlist4go/gfwlist"
)

func TestMake(t *testing.T) {
	blankList, err := gfwlist.BlankList()
	if err != nil {
		t.Error(err)
	}
	pac := Pac{
		BlankList: blankList,
		WhiteList: gfwlist.WHITE_LIST,
		Proxy:     "SOCKS5 127.0.0.1:1080",
	}
	t.Log(pac.String())
}

package pac

import (
	"github.com/gwuhaolin/gfwlist4go/gfwlist"
	"testing"
)

func TestMake(t *testing.T) {
	blankList, err := gfwlist.BlankList()
	if err != nil {
		t.Error(err)
	}
	pac := pac{
		BlankList: blankList,
		WhiteList: gfwlist.WHITE_LIST,
		Proxy:     "SOCKS5 127.0.0.1:1080",
	}
	t.Log(pac.String())
}

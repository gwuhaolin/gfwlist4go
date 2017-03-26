package main

import (
	"os"
	"github.com/gwuhaolin/gfwlist4go/gfwlist"
	"github.com/gwuhaolin/gfwlist4go/pac"
	"io/ioutil"
	"path"
	"log"
)

const (
	OUTPUT = "proxy.pac"
)

func main() {
	proxy := "SOCKS5 127.0.0.1"
	if len(os.Args) == 2 {
		proxy = os.Args[1]
	}
	blankList, err := gfwlist.BlankList()
	if err != nil {
		log.Fatal("获取被墙名单失败", err)
	}
	doc := pac.Pac{
		BlankList: blankList,
		WhiteList: gfwlist.WHITE_LIST,
		Proxy:     proxy,
	}
	str := doc.String()
	err = ioutil.WriteFile(OUTPUT, []byte(str), 0644)
	if err != nil {
		log.Fatal("写文件失败", err)
	}
	dir, _ := os.Getwd()
	log.Print("pac文件输出在 ", path.Join(dir, OUTPUT))
}

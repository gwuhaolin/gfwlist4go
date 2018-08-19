# 为golang编写的自动代理脚本生成

## 在命令行中使用

### 安装

## 在golang中调用
```go
package main

import (
	"github.com/gwuhaolin/gfwlist4go/pac"
	"log"
	"os"
	"path"
)

const (
	OUTPUT = "proxy.js"
)

func main() {
	proxy := "SOCKS5 127.0.0.1"
	if len(os.Args) == 2 {
		proxy = os.Args[1]
	}
	err := pac.SavePac(proxy, OUTPUT)
	if err != nil {
		log.Fatal("写文件失败", err)
	}
	dir, _ := os.Getwd()
	log.Print("pac文件输出在 ", path.Join(dir, OUTPUT))
}
```





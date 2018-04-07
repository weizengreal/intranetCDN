package main

import (
	"fmt"
	"./download"
	"strings"
	"os"
)

func init()  {
	// 初始化全局配置，该函数由 golang 内部调用
	fmt.Println("I am init function!")
}

func main() {
	url := "http://kjds-cdn.aibeike.com/webkjdsfiles/3a6c910b45e644739a80522bfd92d4ea.zip"
	uu := []byte(url)
	ss := strings.LastIndex(url,"/")
	fmt.Println(ss)
	fullname := string(uu[ss+1:])
	fmt.Println(fullname)

	fmt.Println(os.Getwd())
	fmt.Println(download.SendHead(url))
	//download.SendGet("http://kjds-cdn.aibeike.com/webkjdsfiles/3a6c910b45e644739a80522bfd92d4ea.zip",
	//	"/Users/weizeng/vagrant/go/coursewareDelivery/cache",0,100);
}

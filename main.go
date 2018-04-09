package main

import (
	"fmt"
	//"./httputil"
	"./download"
)

func init()  {
	// 初始化全局配置，该函数由 golang 内部调用
	fmt.Println("this is main package init function!")
}

func main() {
	url := "httputil://kjds-cdn.aibeike.com/webkjdsfiles/3a6c910b45e644739a80522bfd92d4ea.zip"
	fmt.Println(download.Download(url))
	//httputil.Server()

	//block := &base.Block{
	//	Start : 0,
	//	End : 200,
	//	BlockId : "123321123321",
	//	BlockSize : 200,
	//	BlockMd5 :"dwdwadw",
	//}
	//fmt.Println(block.AttemptCount)

}

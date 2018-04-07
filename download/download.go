package download

import (
	"../base"
	"flag"
	"os"
	"log"
	"errors"
	"fmt"
)

var root string

var context *base.Context = new(base.Context)

func init() {
	fmt.Println("this is download package init function!")
	r,_ := os.Getwd()
	flag.StringVar(&root,"root", r + "/cache","file where you want to save!")
}

func download(url string) error {
	// 第一步，根据 url 本身和 HEAD 请求初始化文件信息


	// 第二步，根据上下文决定采用哪种模式下载


	// 第三步，根据分片下载的数据合并成完整文件


	return nil
}

// 初始化下载所需要的参数
func assignInit(url string) error {
	prefixName , fullName := base.GetUriName(url)
	length, support, md5, err := SendHead(url)
	if err != nil {
		log.Fatalln(err,"have err!")
		return err
	}
	if !base.CheckFileStat(root) {
		err := os.Mkdir(root,0666)
		if err != nil {
			panic(err)
		}
	}
	// Accept-Ranges 不存在，不支持断点续传模式，要求前端自行处理
	if !support {
		return errors.New("不支持断点续传")
	}
	//
	file := &base.File{
		Path : root+fullName,
		Url : url,
		Name : prefixName,
		Length : length,
		FileMd5 : md5,
	}
	context.File = file


	return nil
}
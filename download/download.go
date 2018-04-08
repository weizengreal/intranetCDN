package download

import (
	"../base"
	"flag"
	"os"
	"log"
	"fmt"
	"sync"
)

const (
	GOODNUM int64 = 1024
	BLOCKSIZE int64 = GOODNUM * GOODNUM * 8
)

// 记录当前下载任务的上下文信息，调试阶段先这么写
var context *base.Context = new(base.Context)

var group sync.WaitGroup

var root string

var blockNum int64

var attemptCount int64

func init() {
	r,_ := os.Getwd()
	flag.StringVar(&root,"root", r + "/cache/","file where you want to save!")
	flag.Int64Var(&blockNum,"block",1,"how many block you want to download once!")
	flag.Int64Var(&attemptCount,"attempt",3,"if download failed,attempt time!")
}

func Download(url string) error {
	// 第一步，根据 url 本身和 HEAD 请求初始化文件信息
	err := assignInit(url)
	if err != nil {
		log.Println(err)
		return err
	}
	// 第二步，根据上下文决定开始下载
	for key,block := range context.FileMap {
		if base.CheckBlockStat(key,block) {
			continue
		}
		group.Add(1)
		go atomDownload(context.Res.Url,key,block)
	}
	group.Wait()

	// 第三步，根据分片下载的数据合并成完整文件
	if context.Res.Support {
		for tmpFilePath,_ := range context.FileMap {
			bytes := base.ReadFile(tmpFilePath)
			err = base.AppendToFile(context.Res.Path,bytes)
			if err != nil {
				log.Println("append file failed!",err)
				return err
			}
		}
	}
	return nil
}

// 初始化下载所需要的参数
func assignInit(url string) error {
	prefixName , fullName := base.GetUriName(url)
	length, support, md5, err := SendHead(url)
	if err != nil {
		log.Println(err,"have err during sendHead to url!")
		return err
	}
	if !base.CheckFileStat(root) {
		err := os.Mkdir(root,0666)
		if err != nil {
			log.Println(err,"create new dir error!")
			return err
		}
	}
	// 当前资源对象
	context.Res = &base.Resource{
		Path : root+fullName,
		Url : url,
		Name : prefixName,
		Length : length,
		FileMd5 : md5,
		Support : support,
	}
	initBlock()
	// 不支持断点续传模式，直接返回交由上层处理
	//if !support {
	//	return nil
	//}


	fmt.Println(*context.Res)
	fmt.Println(context.FileMap)

	return nil
}

// 初始化上下文资源的 block ，区间划分上左闭右开
func initBlock()  {
	length := context.Res.Length
	block := new(base.Block)
	context.FileMap = make(map[string] *base.Block)
	if !context.Res.Support {
		tmpName := base.MD5(context.Res.Name)
		block.Start = 0
		block.End = length
		block.BlockSize = block.End - block.Start
		block.BlockId = tmpName
		context.FileMap[context.Res.Path] = block
		return
	}
	blockUnit := BLOCKSIZE * blockNum
	for i := 0; i <= int(length/(blockUnit)); i++ {
		if i != 0 {
			block.Next = new(base.Block)
			block = block.Next
		}
		tmpName := base.MD5(context.Res.Name + string(i))
		block.Start = int64(i) * blockUnit
		block.End = int64(i + 1) * blockUnit
		block.BlockSize = block.End - block.Start
		block.BlockId = tmpName
		context.FileMap[root + tmpName] = block
	}
}

// 最小的原子化下载工具
func atomDownload(url string, path string, block *base.Block) error {
	length,err := SendGet(url,path,block.Start,block.End - 1)
	if err != nil || length != block.BlockSize{
		if block.AttemptCount < attemptCount {
			block.AttemptCount++
			err = atomDownload(url,path,block)
			log.Println("下载重试中")
		}
	} else {
		// 释放 wait 信息
		group.Done()
	}
	return err
}
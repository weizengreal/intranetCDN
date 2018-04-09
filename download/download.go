package download

import (
	"../base"
	"../comhttp"
	"log"
	"fmt"
)

// 记录当前下载任务的上下文信息，调试阶段先这么写
var context *base.Context = new(base.Context)

func Download(url string) error {
	// 第一步，根据 url 本身和 HEAD 请求初始化文件信息
	err := assignInit(url)
	if err != nil {
		log.Println(err)
		return err
	}
	// 第二步，根据上下文决定开始下载
	for tmpFilePath,block := range context.FileMap {
		if base.CheckBlockStat(tmpFilePath,block) {
			continue
		}
		context.Group.Add(1)
		go atomDownload(context.Res.Url,tmpFilePath,block)
	}
	context.Group.Wait()
	// 第三步，根据分片下载的数据合并成完整文件
	err = base.CreateFileOnly(context.Res.Path)
	if err != nil {
		fmt.Println("create file failed",err)
		return err
	}
	if context.Res.Support {
		for i := 0; i < len(context.TmpPath); i++ {
			bytes := base.ReadFile(context.TmpPath[i])
			err = base.AppendToFile(context.Res.Path,bytes)
			if err != nil {
				log.Println("append file failed!",err)
				return err
			}
		}
	}
	return nil
}

// 最小的原子化下载工具
func atomDownload(url string, path string, block *base.Block) error {
	length,err := comhttp.SendGet(url,path,block.Start,block.End - 1)
	if err != nil || length != block.BlockSize{
		if block.AttemptCount < attemptCount {
			block.AttemptCount++
			err = atomDownload(url,path,block)
			log.Println("download retry!",err,length,*block)
		} else {
			context.Group.Done()
		}
	} else {
		// 释放 wait 信息
		context.Group.Done()
	}
	return err
}
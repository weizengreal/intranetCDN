package download

import (
	"../base"
	"log"
	"fmt"
	"crypto/md5"
	//"context"
)

func Download(url string,context *base.Context) error {
	// 第一步，根据 url 本身和 HEAD 请求初始化文件信息
	err := AssignInit(url,context)
	if err != nil {
		log.Println(err)
		return err
	}
	// 第二步，根据上下文决定开始下载
	DownloadAllBlock(context)
	// 第三步，根据分片下载的数据合并成完整文件
	MergeBlock(context)
	return nil
}

// 将上下文中所有的 block 全部下载下来
func DownloadAllBlock(context *base.Context)  {
	for tmpFilePath,block := range context.FileMap {
		if base.CheckBlockStat(tmpFilePath,block) {
			continue
		}
		context.Group.Add(1)
		go atomDownload(tmpFilePath,block,context)
	}
	context.Group.Wait()
	context.Res.IsComplete = true
}

// 将下载完成的数据合称为一个完整的文件
func MergeBlock(context *base.Context) error {
	err := base.CreateFileOnly(context.Res.Path)
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
			// 计算文件的 MD5 值
			md5Bytes := md5.Sum(bytes)
			context.FileMap[context.TmpPath[i]].BlockMd5 = string(md5Bytes[:])
		}
	}
	return err
}

// 最小的原子化下载工具
func atomDownload(path string, block *base.Block,context *base.Context) error {
	length,err := SendGet(context.Res.Url,path,block.Start,block.End - 1)
	if err != nil || length != block.BlockSize{
		if block.AttemptCount < attemptCount {
			block.AttemptCount++
			err = atomDownload(path,block,context)
			log.Println("download retry!",err,length,*block)
		} else {
			// 重试次数超出限制，释放锁并设置该 block 的下载状态Wie失败
			block.Status = false
			context.Group.Done()
		}
	} else {
		// 释放 wait 信息
		block.Status = true
		context.Group.Done()
	}
	return err
}

// 检测一个上下文信息是否出现变化，false 表示没有变化，true 表示发生变化了
func IsChange(context *base.Context) bool {
	if _,_,_,header ,err := SendHead(context.Res.Url); err != nil {
		return true
	} else if md5Str,ok := header["Content-Md5"]; ok && md5Str == context.Res.FileMd5 {
		return false
	}
	return true
}
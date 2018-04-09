package download

import (
	"../base"
	"../http"
	"log"
	"os"
	"fmt"
	"flag"
)

const (
	GOODNUM int64 = 1024
	BLOCKSIZE int64 = GOODNUM * GOODNUM * 8
)

var root string

var blockNum int64

var attemptCount int64

func init() {
	r,_ := os.Getwd()
	flag.StringVar(&root,"root", r + "/cache/","file where you want to save!")
	flag.Int64Var(&blockNum,"block",1,"how many block you want to download once!")
	flag.Int64Var(&attemptCount,"attempt",3,"if download failed,attempt time!")
}

// 初始化下载所需要的参数
func assignInit(url string) error {
	prefixName , fullName := base.GetUriName(url)
	length, support, md5, err := http.SendHead(url)
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

	fmt.Println(*context.Res)
	fmt.Println(context.FileMap)
	fmt.Println(context.TmpPath)

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
		//context.TmpPath = append(context.TmpPath, context.Res.Path)
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
		// End 的阈值不大于 length
		if block.End > length {
			block.End = length
		}
		block.BlockSize = block.End - block.Start
		block.BlockId = tmpName
		context.FileMap[root + tmpName] = block
		context.TmpPath = append(context.TmpPath, root + tmpName)
	}
}
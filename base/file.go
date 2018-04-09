package base

import "sync"

// 当前任务应该使用的文件信息
type Resource struct {
	Path string
	Url string
	Name string
	Length int64
	FileMd5 string
	Support bool
}

// 一个即将下载的 block 块
type Block struct {
	Start int64
	End int64
	BlockId string
	BlockSize int64
	BlockMd5 string // 该 md5 值用于检测缓存文件
	AttemptCount int64
	Next *Block
}

// 文件下载上下文
type Context struct {
	Res *Resource
	FileMap map[string]*Block // key 为每一个 block 存储的名称，block 为当前需要下载的内容
	TmpPath []string // 保存在 FileMap 中的数据是无序的，这里使用 slice 保存有序的数据分割路径
	Group sync.WaitGroup  // 记录上下文的锁
}

// 文件保存指针，用于告诉 SendGet 函数应该将文件保存在哪里
type FileStroage struct {
	mode int
	path string
	Res *Resource
}

// 为 FileStroage 实现 Writer 接口
func (fs *FileStroage)Writer(p []byte) (int ,error)  {
	if fs.mode == 1 {
		// 复制文件流到文件中
	}else if fs.mode == 2 {
		// 复制文件流到 httputil 的 response 中
	} else {
		// 同时处理1、2两个mode
	}
	// errors.New("")
	return 0,nil
}
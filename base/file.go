package base


// 当前任务应该使用的文件信息
type File struct {
	Path string
	Url string
	Name string
	Length int64
	FileMd5 string
}

// 一个即将下载的 block 块
type Block struct {
	BlockId string
	Start int64
	End int64
	Count int64
	Next *Block
	BlockMd5 string
}

// 文件下载上下文
type Context struct {
	File *File
	FileMap map[string]*Block // key 为每一个 block 存储的位置，block 为当前需要下载的内容
}

// 文件保存指针，用于高速 SendGet 函数应该将文件保存在哪里
type FileStroage struct {
	mode int
	path string
	file *File
}

// 为 FileStroage 实现 Writer 接口
func (fs *FileStroage)Writer(p []byte) (int ,error)  {
	if fs.mode == 1 {
		// 复制文件流到文件中
	}else if fs.mode == 2 {
		// 复制文件流到 http 的 response 中
	} else {
		// 同时处理1、2两个mode
	}
	// errors.New("")
	return 0,nil
}
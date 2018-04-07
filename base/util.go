package base

import (
	"os"
	"log"
)

// 创建一个文件
func CreateFile(path string) (f *os.File, err error) {
	if CheckFileStat(path) {
		DeleteFile(path)
	}
	f, err = os.Create(path)
	if err != nil {
		log.Println("create file" + path + "faild!" )
	}
	return f,err
}

// 删除一个文件
func DeleteFile(path string) error {
	if !CheckFileStat(path) {
		return nil
	}
	err := os.Remove(path)
	if err != nil {
		log.Println("delete file " + path + "faild!")
	}
	return err
}

// 检测文件是否存在
func CheckFileStat(path string) bool {
	if _,err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

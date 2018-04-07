package base

import (
	"os"
	"log"
	"strings"
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

// 根据一个链接返回当前下载内容的文件名称
func GetUriName(url string) (prefixName, fullName string) {
	urlArr := []byte(url)
	index := strings.LastIndex(url,"/")
	if index == -1 {
		fullName = string(urlArr[0:])
	} else {
		fullName = string(urlArr[(index+1):])
	}
	fullNameArr := []byte(fullName)
	pointIndex := strings.LastIndex(fullName,".")
	if pointIndex == -1 {
		prefixName = fullName
	} else {
		prefixName = string(fullNameArr[:pointIndex])
	}
	return prefixName,fullName
}